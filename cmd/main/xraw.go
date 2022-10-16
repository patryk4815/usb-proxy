package main

import (
	"errors"
	"fmt"
	"github.com/google/gousb"
	"log"
	"syscall"
	"unsafe"
)

type xRawReader struct {
	chDone chan bool
	event  usb_raw_control_event
}

type xRawWriter struct {
	chDone chan bool
	event  usb_raw_control_event
}

type xRaw struct {
	fd                       int
	pleaseStopEp0            bool
	setConfigurationDoneOnce bool

	chReader chan xRawReader
	chWriter chan xRawWriter

	lastEvent usb_raw_control_event

	lolHost *xHost
}

func NewRawGadget() *xRaw {
	return &xRaw{
		chReader: make(chan xRawReader),
		chWriter: make(chan xRawWriter),
	}
}

func (s *xRaw) ioctlPtr(req int, ptr unsafe.Pointer) (r1, r2 uintptr, err error) {
	r1, r2, errno := ioctlPtr(s.fd, req, ptr)
	log.Printf("ioctlPtr: r1=%x, r2=%x, errno=%v\n", r1, r2, errno)
	if errno > 0 {
		return r1, r2, errno
	}
	return r1, r2, nil
}

func (s *xRaw) ioctlInt(req int, val uintptr) (r1, r2 uintptr, err error) {
	r1, r2, errno := ioctlInt(s.fd, req, val)
	log.Printf("ioctlInt: r1=%x, r2=%x, errno=%v\n", r1, r2, errno)
	if errno > 0 {
		return r1, r2, errno
	}
	return r1, r2, nil
}

func (s *xRaw) Open() error {
	fd, err := syscall.Open("/dev/raw-gadget", syscall.O_RDWR, 0666)
	if err != nil {
		return err
	}
	s.fd = fd
	return nil
}

func (s *xRaw) Close() error {
	return syscall.Close(s.fd)
}

func (s *xRaw) Init(speed gousb.Speed, driver, device string) error {
	if len(driver) > 127 || len(device) > 127 {
		return errors.New("length invalid driver, device")
	}

	arg := usb_raw_init{
		driver_name: [128]byte{},
		device_name: [128]byte{},
		speed:       byte(speed),
	}
	copy(arg.driver_name[:], driver)
	copy(arg.device_name[:], device)
	arg.driver_name[127] = 0
	arg.device_name[127] = 0

	_, _, err := s.ioctlPtr(USB_RAW_IOCTL_INIT, unsafe.Pointer(&arg))
	if err != nil {
		return err
	}
	return nil
}

func (s *xRaw) Run() error {
	_, _, err := s.ioctlInt(USB_RAW_IOCTL_RUN, 0)
	return err
}

func (s *xRaw) EventFetch(event *usb_raw_control_event) error {
	_, _, err := s.ioctlPtr(USB_RAW_IOCTL_EVENT_FETCH, unsafe.Pointer(event))
	if err != nil {
		//if errno, ok := err.(syscall.Errno); ok && errno == syscall.EINTR {
		//	// TODO: event->length = 4294967295;  // 0xffffffff
		//	return nil
		//}
		return err
	}
	return nil
}

func (s *xRaw) EP0Read(io *usb_raw_control_io) error {
	_, _, err := s.ioctlPtr(USB_RAW_IOCTL_EP0_READ, unsafe.Pointer(io))
	if err != nil {
		//if errno, ok := err.(syscall.Errno); ok && errno == syscall.EBUSY {
		//	return nil
		//}
		return err
	}
	return nil
}

func (s *xRaw) EP0Write(io *usb_raw_control_io) error {
	_, _, err := s.ioctlPtr(USB_RAW_IOCTL_EP0_WRITE, unsafe.Pointer(io))
	if err != nil {
		return err
	}
	return nil
}

func (s *xRaw) EPEnable(desc *usb_endpoint_descriptor) (int, error) {
	// TODO: better handle length size?
	buf := [9]byte{
		desc.bLength,                             //          uint8
		desc.bDescriptorType,                     //  uint8
		desc.bEndpointAddress,                    // uint8
		desc.bmAttributes,                        //     uint8
		uint8((desc.wMaxPacketSize) & 0xff),      //   uint16
		uint8((desc.wMaxPacketSize >> 8) & 0xff), //   uint16
		desc.bInterval,                           //        uint8
		desc.bRefresh,                            //         uint8
		desc.bSynchAddress,                       //    uint8
	}
	epNum, _, err := s.ioctlPtr(USB_RAW_IOCTL_EP_ENABLE, unsafe.Pointer(&buf[0]))
	if err != nil {
		return int(epNum), err
	}
	return int(epNum), nil
}

func (s *xRaw) EPDisable(num uint32) error {
	_, _, err := s.ioctlInt(USB_RAW_IOCTL_EP_DISABLE, uintptr(num))
	if err != nil {
		return err
	}
	return nil
}

func (s *xRaw) EPRead(io *usb_raw_transfer_io) (int, error) {
	// TODO: better handle length size?
	n, _, err := s.ioctlPtr(USB_RAW_IOCTL_EP_READ, unsafe.Pointer(io))
	if err != nil {
		// Ignore failures caused by the test that halts endpoints.
		//if errno, ok := err.(syscall.Errno); ok &&
		//	(errno == syscall.EINPROGRESS || errno == syscall.EBUSY) {
		//	return nil
		//}
		return int(n), err
	}
	return int(n), nil
}

func (s *xRaw) EPWrite(io *usb_raw_transfer_io) (int, error) {
	// TODO: better handle length size?
	n, _, err := s.ioctlPtr(USB_RAW_IOCTL_EP_WRITE, unsafe.Pointer(io))
	if err != nil {
		// Ignore failures caused by the test that halts endpoints.
		//if errno, ok := err.(syscall.Errno); ok &&
		//	(errno == syscall.EINPROGRESS || errno == syscall.EBUSY) {
		//	return nil
		//}
		return int(n), err
	}
	return int(n), nil
}

func (s *xRaw) Configure() error {
	_, _, err := s.ioctlInt(USB_RAW_IOCTL_CONFIGURE, 0)
	if err != nil {
		return err
	}
	return nil
}

func (s *xRaw) VBusDraw(power uint32) error {
	_, _, err := s.ioctlInt(USB_RAW_IOCTL_VBUS_DRAW, uintptr(power))
	if err != nil {
		return err
	}
	return nil
}

func (s *xRaw) EpsInfo() ([]usb_raw_ep_info, error) {
	info := usb_raw_eps_info{}
	num, _, err := s.ioctlPtr(USB_RAW_IOCTL_EPS_INFO, unsafe.Pointer(&info))
	if err != nil {
		return nil, err
	}
	return info.eps[:int(num)], nil
}

func (s *xRaw) Ep0Stall() error {
	_, _, err := s.ioctlInt(USB_RAW_IOCTL_EP0_STALL, 0)
	if err != nil {
		return err
	}
	return nil
}

func (s *xRaw) EpSetHalt(ep int) error {
	_, _, err := s.ioctlInt(USB_RAW_IOCTL_EP_SET_HALT, uintptr(ep))
	if err != nil {
		return err
	}
	return nil
}

type WrapperUsbRawControlIO struct {
	*xRaw
}

func (s *WrapperUsbRawControlIO) Read(output []byte) (int, error) {
	defer func() {
		log.Printf("xRaw Read-End\n")
	}()
	event := <-s.chReader
	log.Printf("xRaw Read lenEvent=%d lenOutput=%d\n", event.event.ctrl.wLength, len(output))
	s.lastEvent = event.event
	defer close(event.chDone)

	io := usb_raw_control_io{}
	io.inner.ep = 0
	io.inner.flags = 0
	io.inner.length = uint32(event.event.ctrl.wLength)
	err := s.EP0Read(&io)
	log.Printf("xRaw Read-ret len=%d err=%#v\n", io.inner.length, err)

	n := copy(output, io.data[:io.inner.length])

	log.Printf("[dump] readOTG: ep=%d data=%q\n",
		io.inner.ep,
		output[:n])

	return n, err
}

func (s *WrapperUsbRawControlIO) Write(input []byte) (int, error) {
	defer func() {
		log.Printf("xRaw Write-End\n")
	}()
	log.Printf("xRaw Write len=%d\n", len(input))

	io := usb_raw_control_io{}
	io.inner.ep = 0
	io.inner.flags = 0
	io.inner.length = uint32(len(input))
	n := copy(io.data[:], input)

	log.Printf("[dump] writeOTG: ep=%d data=%q\n",
		io.inner.ep,
		input)

	err := s.EP0Write(&io)
	log.Printf("xRaw Write-ret len=%d err=%#v\n", io.inner.length, err)

	return n, err
}

func (s *xRaw) EP0Request(event *usb_raw_control_event) bool {
	if event.inner.xtype != USB_RAW_EVENT_CONTROL {
		return true
	}

	bRequestType := USB_TYPE_mapper[int(event.ctrl.bRequestType&USB_TYPE_MASK)]
	bRequest := USB_REQ_mapper[int(event.ctrl.bRequest)]
	wValue := USB_DT_mapper[int(event.ctrl.wValue>>8)]
	log.Printf("EP0Request: bRequestType=%s bRequest=%s wValue=%s length=%d wIndex=%x wValue=%x\n", bRequestType, bRequest, wValue, event.ctrl.wLength, event.ctrl.wIndex, event.ctrl.wValue)

	switch event.ctrl.bRequestType & USB_TYPE_MASK {
	case USB_TYPE_STANDARD:
		switch event.ctrl.bRequest {
		case USB_REQ_GET_DESCRIPTOR:
			switch event.ctrl.wValue >> 8 {
			case USB_DT_DEVICE:
				// io->inner.length = sizeof(usb_device);
				return true
			case USB_DT_DEVICE_QUALIFIER:
				//io->inner.length = sizeof(usb_qualifier);
				return true
			case USB_DT_CONFIG:
				// o->inner.length =
				//					build_config(&io->data[0],
				//						sizeof(io->data), false);
				return true
			case USB_DT_OTHER_SPEED_CONFIG:
				// io->inner.length =
				//					build_config(&io->data[0],
				//						sizeof(io->data), true);
				return true
			case USB_DT_INTERFACE:
				return true
			case USB_DT_ENDPOINT:
				//io->inner.length = sizeof(usb_endpoint_bulk_in);
				return true
			case USB_DT_STRING:
				// io->data[0] = 4;
				// io->data[1] = USB_DT_STRING;
				// if ((event->ctrl.wValue & 0xff) == 0) {
				// 	io->data[2] = 0x09;
				// 	io->data[3] = 0x04;
				// } else {
				// 	io->data[2] = 'x';
				// 	io->data[3] = 0x00;
				// }
				// io->inner.length = 4;
				return true
			case USB_DT_BOS:
				// if (BCD_USB < 0x0201)
				//	 return false;
				// memcpy(&io->data[0], &usb_bos, sizeof(usb_bos));
				// io->inner.length = sizeof(usb_bos);
				return true
			case USB_DT_RPIPE:
				return true
			default:
				panic(fmt.Sprintf("unsupported wValue = %d - %s", event.ctrl.wValue>>8, USB_DT_mapper[int(event.ctrl.wValue>>8)]))
			}
		case USB_REQ_SET_CONFIGURATION:
			if s.setConfigurationDoneOnce {
				return true
			}
			s.setConfigurationDoneOnce = true

			config, err := s.lolHost.dev.Config(int(event.ctrl.wValue))
			log.Printf("err=%#v\n", err)
			if err != nil {
				panic(err)
			}

			for _, ifface := range config.Desc.Interfaces {
				ifs, err := config.Interface(ifface.Number, 0)
				log.Printf("err=%#v\n", err)
				if err != nil {
					panic(err)
				}

				for addr, endpoint := range ifs.Setting.Endpoints {
					infoEndpoint := usb_endpoint_descriptor{
						bLength:          uint8(endpoint.Original.BLength()),
						bDescriptorType:  uint8(endpoint.Original.BDescriptorType()),
						bEndpointAddress: uint8(endpoint.Original.BEndpointAddress()),
						bmAttributes:     uint8(endpoint.Original.BmAttributes()),
						wMaxPacketSize:   uint16(endpoint.Original.WMaxPacketSize()),
						bInterval:        uint8(endpoint.Original.BInterval()),
						bRefresh:         uint8(endpoint.Original.BRefresh()),
						bSynchAddress:    uint8(endpoint.Original.BSynchAddress()),
					}
					epNum, err := s.EPEnable(&infoEndpoint)
					if err != nil {
						panic(err)
					}
					s.Configure()

					log.Printf("endpoint addr=%#v epNum=%d\n", addr, epNum)
					if endpoint.Direction == gousb.EndpointDirectionIn {
						ine, err := ifs.InEndpoint(endpoint.Number)
						if err != nil {
							log.Printf("endpoint ERR addr=%#v epNum=%d err=%v\n", addr, epNum, err)
							panic(err)
						} else {
							go s._EpLoopRead(ine, epNum)
						}
					}
					if endpoint.Direction == gousb.EndpointDirectionOut {
						oute, err := ifs.OutEndpoint(endpoint.Number)
						if err != nil {
							log.Printf("endpoint ERR addr=%#v epNum=%d err=%v\n", addr, epNum, err)
							panic(err)
						} else {
							go s._EpLoopWrite(oute, epNum)
						}
					}
				}
			}

			s._ShowEps()

			//io->inner.length = 0;
			return true
		case USB_REQ_SET_INTERFACE:
			// // TODO: enable/disable endpoints, etc.
			// alt_index = event->ctrl.wValue;
			// io->inner.length = 0;
			return true
		case USB_REQ_GET_INTERFACE:
			// io->data[0] = usb_interface_alts[alt_index]-> bAlternateSetting;
			// io->inner.length = 1;
			return true
		default:
			panic(fmt.Sprintf("unsupported USB_TYPE_STANDARD bRequest = %d - %s", event.ctrl.bRequest, USB_REQ_mapper[int(event.ctrl.bRequest)]))
		}
	case USB_TYPE_CLASS:
		switch event.ctrl.bRequest {
		case USB_REQ_GET_INTERFACE:
			return true
		case USB_REQ_SET_INTERFACE:
			return true
		case USB_REQ_SET_CONFIGURATION:
			return true
		case 0xfe:
			return true
		case 0xff:
			return true
		default:
			panic(fmt.Sprintf("unsupported USB_TYPE_CLASS bRequest = %d - %s", event.ctrl.bRequest, USB_REQ_mapper[int(event.ctrl.bRequest)]))
		}
	case USB_TYPE_VENDOR:
		switch event.ctrl.bRequest {
		case VENDOR_REQ_OUT:
			// io->inner.length = event->ctrl.wLength;
			return true
		case VENDOR_REQ_IN:
			// memcpy(&io->data[0], &vendor_buffer[0], event->ctrl.wLength);
			// io->inner.length = event->ctrl.wLength;
			return true
		default:
			panic(fmt.Sprintf("unsupported USB_TYPE_VENDOR bRequest = %d - %s", event.ctrl.bRequest, USB_REQ_mapper[int(event.ctrl.bRequest)]))
		}
	default:
		panic("unsupported bRequestType")
	}

	return true
}

func (s *xRaw) EP0Loop() error {
	log.Println("Start for EP0, thread id")

	for !s.pleaseStopEp0 {
		log.Println("EP0Loop")

		event := usb_raw_control_event{}
		event.inner.xtype = 0
		event.inner.length = uint32(sizeOf(event.ctrl))
		if err := s.EventFetch(&event); err != nil {
			return err
		}

		log.Printf("EventFetch: inner=%#v ctrl=%#v\n", event.inner, event.ctrl)

		if !s.EP0Request(&event) {
			log.Printf("loop Ep0Stall\n")
			s.Ep0Stall()
			continue
		}

		if event.inner.xtype == USB_RAW_EVENT_CONTROL {
			if event.ctrl.bRequestType&USB_DIR_IN > 0 {
				// write
				ref := xRawWriter{
					chDone: make(chan bool),
					event:  event,
				}
				s.chWriter <- ref
				<-ref.chDone
			} else {
				// read
				ref := xRawReader{
					chDone: make(chan bool),
					event:  event,
				}
				s.chReader <- ref
				<-ref.chDone
			}
		}
	}
	return nil
}

func (s *xRaw) _ShowEps() error {
	infos, err := s.EpsInfo()
	if err != nil {
		return err
	}

	for i, info := range infos {
		caps := info.caps.GetType()
		log.Printf("ep #%d:"+
			"  name: %s"+
			"  addr: %x"+
			"  type: iso=%v bulk=%v int=%v"+
			"  dir: in=%v out=%v"+
			"  maxpacket_limit: %x"+
			"  max_streams: %x\n",
			i,
			info.name,
			info.addr,
			caps.type_iso, caps.type_bulk, caps.type_int,
			caps.dir_in, caps.dir_out,
			info.limits.maxpacket_limit,
			info.limits.max_streams,
		)
	}
	return nil
}

type WrapperUsbRawTransferIO struct {
	*xRaw
	epNum int
}

func (s *WrapperUsbRawTransferIO) Read(p []byte) (int, error) {
	pkt := usb_raw_transfer_io{}
	//pkt := usb_raw_int_io{}
	pkt.inner.ep = uint16(s.epNum)
	pkt.inner.flags = 0
	pkt.inner.length = uint32(sizeOf(pkt.data))

	log.Printf("WrapperUsbRawTransferIO Read len=%d, pLen=%d\n", pkt.inner.length, len(p))

	n, err := s.EPRead(&pkt)

	log.Printf("WrapperUsbRawTransferIO Read-result len=%d, n=%d, err=%v\n", pkt.inner.length, n, err)

	if err != nil {
		panic(err)
	}
	copy(p, pkt.data[:n])

	return n, nil
}

func (s *WrapperUsbRawTransferIO) Write(p []byte) (int, error) {
	pkt := usb_raw_transfer_io{}
	//pkt := usb_raw_int_io{}
	pkt.inner.ep = uint16(s.epNum)
	pkt.inner.flags = 0
	pkt.inner.length = uint32(len(p))
	copy(pkt.data[:], p)

	log.Printf("WrapperUsbRawTransferIO Write len=%d\n", len(p))

	n, err := s.EPWrite(&pkt)

	log.Printf("WrapperUsbRawTransferIO Write-result len=%d n=%d err=%#v\n", len(p), n, err)
	if err != nil {
		panic(err)
	}
	return n, nil
}

func (s *xRaw) _EpLoopRead(in *gousb.InEndpoint, epNum int) error {
	log.Println("_EpLoopRead")
	writer := &WrapperUsbRawTransferIO{xRaw: s, epNum: epNum}
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

func (s *xRaw) _EpLoopWrite(out *gousb.OutEndpoint, epNum int) error {
	log.Println("_EpLoopWrite")
	reader := &WrapperUsbRawTransferIO{xRaw: s, epNum: epNum}
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

func (s *xRaw) _WorkerDo() {
	if err := s.Open(); err != nil {
		panic(err)
	}

	//if err := s.Init(gousb.SpeedHigh, "fe980000.usb", "fe980000.usb"); err != nil {
	if err := s.Init(gousb.SpeedHigh, "fe980000.usb", "fe980000.usb"); err != nil {
		panic(err)
	}

	if err := s.Run(); err != nil {
		panic(err)
	}

	if err := s.EP0Loop(); err != nil {
		panic(err)
	}
}
