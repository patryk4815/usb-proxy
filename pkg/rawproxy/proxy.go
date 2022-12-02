package rawproxy

import (
	"bytes"
	"fmt"
	"github.com/google/gousb"
	"github.com/lunixbochs/struc"
	_ "github.com/lunixbochs/struc"
	"github.com/patryk4815/usb-proxy/pkg/rawgadget"
	"io"
	"log"
)

type XX_RawReaderWriter struct {
	ChDone chan bool
	Event  UsbEventCtrl
}

type XX_www struct {
	otgDevice *rawgadget.XX_Raw

	pleaseStopEp0            bool
	setConfigurationDoneOnce bool

	chReader chan XX_RawReaderWriter
	chWriter chan XX_RawReaderWriter

	lastEvent UsbEventCtrl

	hostDevice IHostDevice
}

func New() *XX_www {
	s := rawgadget.NewRawGadget()
	if err := s.Open(); err != nil {
		panic(err)
	}

	// TODO: config argument

	www := &XX_www{
		otgDevice: s,

		chReader: make(chan XX_RawReaderWriter),
		chWriter: make(chan XX_RawReaderWriter),
	}
	return www
}

func (s *XX_www) Open() {
	if err := s.otgDevice.Init(gousb.SpeedHigh, "fe980000.usb", "fe980000.usb"); err != nil {
		panic(err)
	}
	if err := s.otgDevice.Run(); err != nil {
		panic(err)
	}
}

func (s *XX_www) GetChWriter() <-chan XX_RawReaderWriter {
	return s.chWriter
}

func (s *XX_www) GetLastEvent() UsbEventCtrl {
	return s.lastEvent
}

func (s *XX_www) Close() error {
	// TODO: implement better
	return s.otgDevice.Close()
}

type UsbEventCtrl struct {
	RawEventType   rawgadget.Usb_raw_event_type
	Type           rawgadget.Usb_Type
	Request        rawgadget.Usb_Req
	DescriptorType rawgadget.Usb_DT
	Direction      rawgadget.Usb_Dir

	RawCtrlReq rawgadget.Usb_ctrlrequest
}

func (s *UsbEventCtrl) UpdateFrom(newobj *rawgadget.Usb_raw_event) error {
	s.RawEventType = newobj.Type

	tmp := &rawgadget.Usb_ctrlrequest{}
	err := struc.Unpack(bytes.NewReader(newobj.Data), tmp)
	if err != nil {
		return err
	}
	s.RawCtrlReq = *tmp

	s.Type = rawgadget.Usb_Type(tmp.BRequestType & uint8(rawgadget.USB_TYPE_MASK))
	s.Request = rawgadget.Usb_Req(tmp.BRequest)
	s.DescriptorType = rawgadget.Usb_DT(tmp.WValue >> 8)

	isIn := (uint8(tmp.BRequestType) & uint8(rawgadget.USB_DIR_IN)) > 0
	if isIn {
		s.Direction = rawgadget.USB_DIR_IN
	} else {
		s.Direction = rawgadget.USB_DIR_OUT
	}
	return nil
}

type IHostDevice interface {
	GetConfig(wValue uint16) (*gousb.Config, error)
}

func (s *XX_www) SetHostProxy(h IHostDevice) {
	s.hostDevice = h
}

func (s *XX_www) FuncOnConfigure(event *UsbEventCtrl) error {
	cfg, err := s.hostDevice.GetConfig(event.RawCtrlReq.WValue)
	if err != nil {
		panic(err) // TODO: fix
	}

	for _, ifface := range cfg.Desc.Interfaces {
		ifs, err := cfg.Interface(ifface.Number, 0) // TODO: alt?
		log.Printf("err=%#v\n", err)
		if err != nil {
			panic(err)
		}

		for addr, endpoint := range ifs.Setting.Endpoints {
			infoEndpoint := rawgadget.Usb_endpoint_descriptor{
				BLength:          uint8(endpoint.Original.BLength()),
				BDescriptorType:  uint8(endpoint.Original.BDescriptorType()),
				BEndpointAddress: uint8(endpoint.Original.BEndpointAddress()),
				BmAttributes:     uint8(endpoint.Original.BmAttributes()),
				WMaxPacketSize:   uint16(endpoint.Original.WMaxPacketSize()),
				BInterval:        uint8(endpoint.Original.BInterval()),
				BRefresh:         uint8(endpoint.Original.BRefresh()),
				BSynchAddress:    uint8(endpoint.Original.BSynchAddress()),
			}
			epNum, err := s.otgDevice.EPEnable(&infoEndpoint)
			if err != nil {
				panic(err)
			}

			log.Printf("endpoint addr=%#v EpNum=%d\n", addr, epNum)
			if endpoint.Direction == gousb.EndpointDirectionIn {
				ine, err := ifs.InEndpoint(endpoint.Number)
				if err != nil {
					log.Printf("endpoint ERR addr=%#v EpNum=%d err=%v\n", addr, epNum, err)
					panic(err)
				} else {
					go s._EpLoopRead(ine, epNum)
				}
			}
			if endpoint.Direction == gousb.EndpointDirectionOut {
				oute, err := ifs.OutEndpoint(endpoint.Number)
				if err != nil {
					log.Printf("endpoint ERR addr=%#v EpNum=%d err=%v\n", addr, epNum, err)
					panic(err)
				} else {
					go s._EpLoopWrite(oute, epNum)
				}
			}
		}
	}

	// TODO: to chyba powinno byc na koncu?
	if err := s.otgDevice.Configured(); err != nil {
		panic(err)
	}

	return nil
}

func (s *XX_www) EventFetchCtrlReq(out *UsbEventCtrl) (int, error) {
	out2 := &rawgadget.Usb_ctrlrequest{}
	out1 := &rawgadget.Usb_raw_event{}
	out1.Data = make([]byte, rawgadget.Sizeof(*out2))

	n, err := s.otgDevice.EventFetch(out1)
	if err != nil {
		return 0, err
	}

	if err := out.UpdateFrom(out1); err != nil {
		return n, err
	}

	return n, nil
}

func (s *XX_www) EP0Request(event *UsbEventCtrl) bool {
	if event.RawEventType != rawgadget.USB_RAW_EVENT_CONTROL {
		return true
	}

	bRequestType := event.Type
	bRequest := event.Request
	wValue := event.DescriptorType
	log.Printf("EP0Request: BRequestType=%s BRequest=%s WValue=%s\n", bRequestType, bRequest, wValue)

	switch bRequestType {
	case rawgadget.USB_TYPE_STANDARD:
		switch bRequest {
		case rawgadget.USB_REQ_GET_DESCRIPTOR:
			switch wValue {
			case rawgadget.USB_DT_DEVICE:
				// io->Inner.Length = sizeof(usb_device);
				return true
			case rawgadget.USB_DT_DEVICE_QUALIFIER:
				//io->Inner.Length = sizeof(usb_qualifier);
				return true
			case rawgadget.USB_DT_CONFIG:
				// o->Inner.Length =
				//					build_config(&io->Data[0],
				//						sizeof(io->Data), false);
				return true
			case rawgadget.USB_DT_OTHER_SPEED_CONFIG:
				// io->Inner.Length =
				//					build_config(&io->Data[0],
				//						sizeof(io->Data), true);
				return true
			case rawgadget.USB_DT_INTERFACE:
				return true
			case rawgadget.USB_DT_ENDPOINT:
				//io->Inner.Length = sizeof(usb_endpoint_bulk_in);
				return true
			case rawgadget.USB_DT_STRING:
				// io->Data[0] = 4;
				// io->Data[1] = USB_DT_STRING;
				// if ((Event->Ctrl.WValue & 0xff) == 0) {
				// 	io->Data[2] = 0x09;
				// 	io->Data[3] = 0x04;
				// } else {
				// 	io->Data[2] = 'x';
				// 	io->Data[3] = 0x00;
				// }
				// io->Inner.Length = 4;
				return true
			case rawgadget.USB_DT_BOS:
				// if (BCD_USB < 0x0201)
				//	 return false;
				// memcpy(&io->Data[0], &usb_bos, sizeof(usb_bos));
				// io->Inner.Length = sizeof(usb_bos);
				return true
			case rawgadget.USB_DT_RPIPE:
				return true
			default:
				panic(fmt.Sprintf("unsupported WValue = %d - %s", wValue, wValue))
			}
		case rawgadget.USB_REQ_SET_CONFIGURATION:
			if s.setConfigurationDoneOnce {
				return true
			}
			s.setConfigurationDoneOnce = true

			s.FuncOnConfigure(event)
			s.otgDevice.DebugShowEps()

			//io->Inner.Length = 0;
			return true
		case rawgadget.USB_REQ_SET_INTERFACE:
			// // TODO: enable/disable Endpoints, etc.
			// alt_index = Event->Ctrl.WValue;
			// io->Inner.Length = 0;
			return true
		case rawgadget.USB_REQ_GET_INTERFACE:
			// io->Data[0] = usb_interface_alts[alt_index]-> BAlternateSetting;
			// io->Inner.Length = 1;
			return true
		default:
			panic(fmt.Sprintf("unsupported USB_TYPE_STANDARD BRequest = %d - %s", bRequest, bRequest))
		}
	case rawgadget.USB_TYPE_CLASS:
		switch bRequest {
		case rawgadget.USB_REQ_GET_INTERFACE:
			return true
		case rawgadget.USB_REQ_SET_INTERFACE:
			return true
		case rawgadget.USB_REQ_SET_CONFIGURATION:
			return true
		case 0xfe: // TODO: WTF?
			return true
		case 0xff: // TODO: WTF?
			return true
		default:
			panic(fmt.Sprintf("unsupported USB_TYPE_CLASS BRequest = %d - %s", bRequest, bRequest))
		}
	case rawgadget.USB_TYPE_VENDOR:
		switch event.Request {
		case rawgadget.VENDOR_REQ_OUT:
			// io->Inner.Length = Event->Ctrl.WLength;
			return true
		case rawgadget.VENDOR_REQ_IN:
			// memcpy(&io->Data[0], &vendor_buffer[0], Event->Ctrl.WLength);
			// io->Inner.Length = Event->Ctrl.WLength;
			return true
		default:
			panic(fmt.Sprintf("unsupported USB_TYPE_VENDOR BRequest = %d - %s", bRequest, bRequest))
		}
	default:
		panic("unsupported BRequestType")
	}

	return true
}

func (s *XX_www) EP0Loop() error {
	log.Println("Start for EP0, thread id")

	for !s.pleaseStopEp0 {
		log.Println("EP0Loop??")

		event := UsbEventCtrl{}
		_, err := s.EventFetchCtrlReq(&event) // TODO: n
		if err == io.EOF && event.RawEventType == rawgadget.USB_RAW_EVENT_CONNECT {
			// OK XD
		} else if err != nil {
			panic(err)
			return err
		}
		log.Printf("EventFetch: %#v\n", event)

		if !s.EP0Request(&event) {
			log.Printf("loop Ep0Stall\n")
			s.otgDevice.Ep0Stall() // TODO: err
			continue
		}

		if event.RawEventType == rawgadget.USB_RAW_EVENT_CONTROL {
			if event.Direction == rawgadget.USB_DIR_IN {
				// write
				ref := XX_RawReaderWriter{
					ChDone: make(chan bool),
					Event:  event,
				}
				s.chWriter <- ref
				<-ref.ChDone
			} else {
				// read
				ref := XX_RawReaderWriter{
					ChDone: make(chan bool),
					Event:  event,
				}
				s.chReader <- ref
				<-ref.ChDone
			}
		}
	}
	return nil
}

func (s *XX_www) _EpLoopRead(in *gousb.InEndpoint, epNum int) error {
	log.Println("_EpLoopRead")
	writer := &WrapperUsbRawTransferIO{Self: s, EpNum: epNum}
	_ = writer

	for {
		buf := make([]byte, in.Desc.MaxPacketSize)
		n, err := in.Read(buf)
		buf = buf[:n]
		log.Printf("_EpLoopRead-end: err=%v n=%d buf=%q\n", err, n, buf)
		writer.Write(buf)
	}

	return nil
}

func (s *XX_www) _EpLoopWrite(out *gousb.OutEndpoint, epNum int) error {
	log.Println("_EpLoopWrite")
	reader := &WrapperUsbRawTransferIO{Self: s, EpNum: epNum}
	_ = reader

	for {
		buf := make([]byte, out.Desc.MaxPacketSize)
		n, err := reader.Read(buf)
		buf = buf[:n]
		log.Printf("_EpLoopWrite-end: err=%v n=%d buf=%q\n", err, n, buf)
		out.Write(buf)
	}

	return nil
}
