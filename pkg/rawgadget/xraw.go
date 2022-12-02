package rawgadget

import (
	"bytes"
	"errors"
	"github.com/google/gousb"
	"github.com/lunixbochs/struc"
	"log"
	"syscall"
	"unsafe"
)

type XX_Raw struct {
	fd int
}

func NewRawGadget() *XX_Raw {
	return &XX_Raw{}
}

func (s *XX_Raw) ioctlPtr(req int, ptr unsafe.Pointer) (r1, r2 uintptr, err error) {
	r1, r2, errno := ioctlPtr(s.fd, req, ptr)
	log.Printf("ioctlPtr: r1=%x, r2=%x, errno=%v\n", r1, r2, errno)
	if errno > 0 {
		return r1, r2, errno
	}
	return r1, r2, nil
}

func (s *XX_Raw) ioctlInt(req int, val uintptr) (r1, r2 uintptr, err error) {
	r1, r2, errno := ioctlInt(s.fd, req, val)
	log.Printf("ioctlInt: r1=%x, r2=%x, errno=%v\n", r1, r2, errno)
	if errno > 0 {
		return r1, r2, errno
	}
	return r1, r2, nil
}

func (s *XX_Raw) Open() error {
	// TODO: custom name
	fd, err := syscall.Open("/dev/raw-gadget", syscall.O_RDWR, 0666)
	if err != nil {
		return err
	}
	s.fd = fd
	return nil
}

func (s *XX_Raw) Close() error {
	return syscall.Close(s.fd)
}

func (s *XX_Raw) Init(speed gousb.Speed, driver, device string) error {
	if len(driver) > UDC_NAME_LENGTH_MAX || len(device) > UDC_NAME_LENGTH_MAX {
		return errors.New("length device or driver")
	}

	arg := Usb_raw_init{
		driver_name: [UDC_NAME_LENGTH_MAX]byte{},
		device_name: [UDC_NAME_LENGTH_MAX]byte{},
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

func (s *XX_Raw) Run() error {
	_, _, err := s.ioctlInt(USB_RAW_IOCTL_RUN, 0)
	return err
}

func (s *XX_Raw) EventFetch(ev *Usb_raw_event) (int, error) {
	buf := bytes.NewBuffer(nil)
	in := &_packer_Usb_raw_event{
		Type:   uint32(ev.Type),
		Length: uint32(len(ev.Data)),
		Data:   ev.Data,
	}
	err := struc.Pack(buf, in)
	if err != nil {
		return 0, err
	}

	log.Printf("USB_RAW_IOCTL_EVENT_FETCH in: %#v\n", in)
	log.Printf("USB_RAW_IOCTL_EVENT_FETCH pre: %q\n", buf.Bytes())

	_, _, err = s.ioctlPtr(USB_RAW_IOCTL_EVENT_FETCH, unsafe.Pointer(&(buf.Bytes()[0])))
	if err != nil {
		return 0, err
	}
	log.Printf("USB_RAW_IOCTL_EVENT_FETCH post: %q\n", buf.Bytes())

	out := &_packer_Usb_raw_event{}
	err = struc.Unpack(buf, out)
	if err != nil {
		return 0, err
	}

	ev.Type = Usb_raw_event_type(out.Type)
	ev.Data = append(ev.Data[:0], out.Data...) // hack
	log.Printf("Usb_raw_event post: %#v\n", ev)

	n := int(out.Length)
	return n, nil
}

func (s *XX_Raw) _epReadWriter(ioctl int, isReader bool, ep *Usb_raw_ep_io) (int, error) {
	if len(ep.Data) > PAGE_SIZE {
		return 0, errors.New("err buf over limit PAGE_SIZE")
	}
	if ep.Ep > USB_RAW_EPS_NUM_MAX {
		return 0, errors.New("err Ep over limit USB_RAW_EPS_NUM_MAX")
	}

	buf := bytes.NewBuffer(nil)
	in := &_packer_Usb_raw_ep_io{
		Ep:     uint16(ep.Ep),
		Flags:  uint16(ep.Flags),
		Length: uint32(len(ep.Data)),
		Data:   ep.Data,
	}
	err := struc.Pack(buf, in)
	if err != nil {
		return 0, err
	}

	num, _, err := s.ioctlPtr(ioctl, unsafe.Pointer(&(buf.Bytes()[0])))
	if err != nil {
		return int(num), err
	}

	if isReader {
		out := &_packer_Usb_raw_ep_io{}
		err = struc.Unpack(buf, out)
		if err != nil {
			return 0, err
		}
		ep.Data = append(ep.Data[:0], out.Data...) // hack overwrite
	}

	return int(num), nil
}

func (s *XX_Raw) EP0Read(ep *Usb_raw_ep_io) (int, error) {
	return s._epReadWriter(USB_RAW_IOCTL_EP0_READ, true, ep)
}

func (s *XX_Raw) EP0Write(ep *Usb_raw_ep_io) (int, error) {
	return s._epReadWriter(USB_RAW_IOCTL_EP0_WRITE, false, ep)
}

func (s *XX_Raw) EPRead(ep *Usb_raw_ep_io) (int, error) {
	return s._epReadWriter(USB_RAW_IOCTL_EP_READ, true, ep)
}

func (s *XX_Raw) EPWrite(ep *Usb_raw_ep_io) (int, error) {
	return s._epReadWriter(USB_RAW_IOCTL_EP_WRITE, false, ep)
}

func (s *XX_Raw) EPEnable(desc *Usb_endpoint_descriptor) (int, error) {
	buf := [...]byte{
		desc.BLength,                             // uint8
		desc.BDescriptorType,                     // uint8
		desc.BEndpointAddress,                    // uint8
		desc.BmAttributes,                        // uint8
		uint8((desc.WMaxPacketSize) & 0xff),      // uint16
		uint8((desc.WMaxPacketSize >> 8) & 0xff), // uint16
		desc.BInterval,                           // uint8
		desc.BRefresh,                            // uint8
		desc.BSynchAddress,                       // uint8
	}
	epNum, _, err := s.ioctlPtr(USB_RAW_IOCTL_EP_ENABLE, unsafe.Pointer(&buf[0]))
	if err != nil {
		return int(epNum), err
	}
	return int(epNum), nil
}

func (s *XX_Raw) EPDisable(num uint32) error {
	_, _, err := s.ioctlInt(USB_RAW_IOCTL_EP_DISABLE, uintptr(num))
	if err != nil {
		return err
	}
	return nil
}

func (s *XX_Raw) Configured() error {
	_, _, err := s.ioctlInt(USB_RAW_IOCTL_CONFIGURE, 0)
	if err != nil {
		return err
	}
	return nil
}

func (s *XX_Raw) VBusDraw(power uint32) error {
	_, _, err := s.ioctlInt(USB_RAW_IOCTL_VBUS_DRAW, uintptr(power))
	if err != nil {
		return err
	}
	return nil
}

func (s *XX_Raw) EpsInfo() ([]Usb_raw_ep_info, error) {
	info := Usb_raw_eps_info{}
	num, _, err := s.ioctlPtr(USB_RAW_IOCTL_EPS_INFO, unsafe.Pointer(&info))
	if err != nil {
		return nil, err
	}
	return info.Eps[:int(num)], nil
}

func (s *XX_Raw) Ep0Stall() error {
	_, _, err := s.ioctlInt(USB_RAW_IOCTL_EP0_STALL, 0)
	if err != nil {
		return err
	}
	return nil
}

func (s *XX_Raw) EpSetHalt(ep int) error {
	_, _, err := s.ioctlInt(USB_RAW_IOCTL_EP_SET_HALT, uintptr(ep))
	if err != nil {
		return err
	}
	return nil
}

func (s *XX_Raw) EpClearHalt(ep int) error {
	_, _, err := s.ioctlInt(USB_RAW_IOCTL_EP_CLEAR_HALT, uintptr(ep))
	if err != nil {
		return err
	}
	return nil
}

func (s *XX_Raw) EpSetWedge(ep int) error {
	_, _, err := s.ioctlInt(USB_RAW_IOCTL_EP_SET_WEDGE, uintptr(ep))
	if err != nil {
		return err
	}
	return nil
}

func (s *XX_Raw) DebugShowEps() error {
	infos, err := s.EpsInfo()
	if err != nil {
		return err
	}

	for i, info := range infos {
		caps := info.Caps.GetType()
		log.Printf("Ep #%d:"+
			"  Name: %s"+
			"  Addr: %x"+
			"  type: iso=%v bulk=%v int=%v"+
			"  dir: in=%v out=%v"+
			"  Maxpacket_limit: %x"+
			"  Max_streams: %x\n",
			i,
			info.Name,
			info.Addr,
			caps.Type_iso, caps.Type_bulk, caps.Type_int,
			caps.Dir_in, caps.Dir_out,
			info.Limits.Maxpacket_limit,
			info.Limits.Max_streams,
		)
	}
	return nil
}
