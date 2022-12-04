package rawproxy

import (
	"github.com/google/gousb"
	_ "github.com/lunixbochs/struc"
	"github.com/patryk4815/usb-proxy/pkg/rawgadget"
	log "github.com/sirupsen/logrus"
	"io"
	"time"
)

type XX_RawReaderWriter struct {
	ChDone chan error
	Event  rawgadget.UsbEventCtrl
}

type XX_www struct {
	otgDevice *rawgadget.XX_Raw

	pleaseStopEp0            bool
	setConfigurationDoneOnce bool

	chReader chan XX_RawReaderWriter
	chWriter chan XX_RawReaderWriter

	hostDevice IHostDevice

	endpointsInfo map[gousb.EndpointAddress]gousb.EndpointDesc
}

func New() *XX_www {
	s := rawgadget.NewRawGadget()
	if err := s.Open(); err != nil {
		panic(err)
	}

	www := &XX_www{
		otgDevice: s,

		chReader: make(chan XX_RawReaderWriter),
		chWriter: make(chan XX_RawReaderWriter),

		endpointsInfo: make(map[gousb.EndpointAddress]gousb.EndpointDesc),
	}
	return www
}

func (s *XX_www) Open(driver, device string) {
	if err := s.otgDevice.Init(gousb.SpeedHigh, driver, device); err != nil {
		log.WithError(err).Fatalf("raw-gaget init err")
	}
	if err := s.otgDevice.Run(); err != nil {
		log.WithError(err).Fatalf("raw-gaget Run err")
	}
}

func (s *XX_www) GetChWriter() <-chan XX_RawReaderWriter {
	return s.chWriter
}

func (s *XX_www) Close() error {
	// TODO: implement better
	return s.otgDevice.Close()
}

type IHostDevice interface {
	GetConfig(wValue uint16) (*gousb.Config, error)
}

func (s *XX_www) SetHostProxy(h IHostDevice) {
	s.hostDevice = h
}

func (s *XX_www) FuncOnConfigure(event *rawgadget.UsbEventCtrl) error {
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
			addr := addr
			endpoint := endpoint

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
			log.WithField("ifface.Number", ifface.Number).WithField("endpoint", endpoint).Debug("host endpoints")

			epNum, err := s.otgDevice.EPEnable(&infoEndpoint)
			if err != nil {
				panic(err)
			}
			s.endpointsInfo[gousb.EndpointAddress(epNum)] = endpoint

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

	return nil
}

func (s *XX_www) EventFetchCtrlReq(out *rawgadget.UsbEventCtrl) (int, error) {
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

func (s *XX_www) EP0Request(event *rawgadget.UsbEventCtrl) bool {
	if event.RawEventType != rawgadget.USB_RAW_EVENT_CONTROL {
		return true
	}
	switch event.Type {
	case rawgadget.USB_TYPE_CLASS:
		//if event.RawCtrlReq.WLength <= 0 {
		//	return false
		//}
	case rawgadget.USB_TYPE_STANDARD:
		switch event.Request {
		case rawgadget.USB_REQ_SET_INTERFACE:
			log.Fatalf("TODO implement")
		case rawgadget.USB_REQ_SET_CONFIGURATION:
			if s.setConfigurationDoneOnce {
				return true
			}
			s.setConfigurationDoneOnce = true

			s.FuncOnConfigure(event)

			if err := s.otgDevice.Configured(); err != nil {
				log.WithError(err).Fatalf("not configured")
			}

			// TODO: skad to pobrac?
			//if err := s.otgDevice.VBusDraw(0x32 * 5); err != nil {
			//	panic(err)
			//}

			s.otgDevice.DebugShowEps()
			return true
		}
	}
	return true
}

func (s *XX_www) EP0Loop() error {
	log.Println("Start for EP0, thread id")

	for !s.pleaseStopEp0 {
		log.Println("EP0Loop??")

		event := rawgadget.UsbEventCtrl{}
		_, err := s.EventFetchCtrlReq(&event) // TODO: n
		log.Printf("EventFetch: %+v\n", event)

		if err == io.EOF && event.RawEventType == rawgadget.USB_RAW_EVENT_CONNECT {
			// OK XD
			continue
		} else if err != nil {
			panic(err)
			return err
		}

		if !s.EP0Request(&event) {
			log.Warning("Ep0Stall")
			s.otgDevice.Ep0Stall() // TODO: err
			continue
		}

		if event.RawEventType == rawgadget.USB_RAW_EVENT_CONTROL {
			if event.Direction == rawgadget.USB_DIR_IN {
				// write
				ref := XX_RawReaderWriter{
					ChDone: make(chan error, 1),
					Event:  event,
				}
				s.chWriter <- ref
				err := <-ref.ChDone
				if err != nil {
					log.Warning("Ep0Stall_TMP write")
					//s.otgDevice.Ep0Stall() // TODO: err
				}
			} else {
				// read
				ref := XX_RawReaderWriter{
					ChDone: make(chan error, 1),
					Event:  event,
				}
				s.chReader <- ref
				err := <-ref.ChDone
				if err != nil {
					log.Warning("Ep0Stall_TMP read")
					//s.otgDevice.Ep0Stall() // TODO: err
				}
			}
		}
	}
	return nil
}

func (s *XX_www) _EpLoopRead(in *gousb.InEndpoint, epNum int) error {
	log.Println("_EpLoopRead")
	writer := &WrapperUsbEPX{Self: s, EpNum: epNum}
	_ = writer

	for {
		buf := make([]byte, in.Desc.MaxPacketSize)
		n, err := in.Read(buf)
		if err != nil {
			log.WithError(err).WithField("PollInterval", in.Desc.PollInterval).Warning("sleeping err read _EpLoopRead")
			time.Sleep(in.Desc.PollInterval)
			continue
		}
		buf = buf[:n]

		_, err = writer.Write(buf)
		if err != nil {
			log.WithError(err).WithField("PollInterval", in.Desc.PollInterval).Warning("sleeping err write _EpLoopRead")
			time.Sleep(in.Desc.PollInterval)
			continue
		}
	}

	return nil
}

func (s *XX_www) _EpLoopWrite(out *gousb.OutEndpoint, epNum int) error {
	log.Println("_EpLoopWrite")
	reader := &WrapperUsbEPX{Self: s, EpNum: epNum}
	_ = reader

	for {
		buf := make([]byte, out.Desc.MaxPacketSize)
		n, err := reader.Read(buf)
		if err != nil {
			log.WithError(err).WithField("PollInterval", out.Desc.PollInterval).Warning("sleeping err read _EpLoopWrite")
			time.Sleep(out.Desc.PollInterval)
			continue
		}
		buf = buf[:n]

		_, err = out.Write(buf)
		if err != nil {
			log.WithError(err).WithField("PollInterval", out.Desc.PollInterval).Warning("sleeping err write _EpLoopWrite")
			time.Sleep(out.Desc.PollInterval)
			continue
		}
	}

	return nil
}
