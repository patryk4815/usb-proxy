package rawgadget

import _ "github.com/lunixbochs/struc"

const UDC_NAME_LENGTH_MAX = 128

type Usb_raw_init struct {
	driver_name [UDC_NAME_LENGTH_MAX]byte
	device_name [UDC_NAME_LENGTH_MAX]byte
	speed       byte
}

type Usb_raw_event_type uint32 // TODO: stringify

const (
	USB_RAW_EVENT_INVALID Usb_raw_event_type = 0
	USB_RAW_EVENT_CONNECT Usb_raw_event_type = 1
	USB_RAW_EVENT_CONTROL Usb_raw_event_type = 2
)

type Usb_raw_event struct {
	Type Usb_raw_event_type
	Data []byte
}

type _ioctl_Usb_raw_event struct {
	Type   uint32
	Length uint32
	// Data []byte
}

type _packer_Usb_raw_event struct {
	Type   uint32 `struc:"uint32,little"`
	Length uint32 `struc:"uint32,sizeof=Data,little"`
	Data   []byte `struc:"[]uint8"`
}

type Usb_raw_ep_io struct {
	Ep    uint16 // MAX = USB_RAW_EPS_NUM_MAX
	Flags uint16
	Data  []byte
}

type _ioctl_Usb_raw_ep_io struct {
	Ep     uint16
	Flags  uint16
	Length uint32
	// Data []byte
}

type _packer_Usb_raw_ep_io struct {
	Ep     uint16 `struc:"uint16,little"`
	Flags  uint16 `struc:"uint16,little"`
	Length uint32 `struc:"uint32,sizeof=Data,little"`
	Data   []byte `struc:"[]uint8"`
}

const USB_RAW_EPS_NUM_MAX = 30
const USB_RAW_EP_NAME_MAX = 16
const USB_RAW_EP_ADDR_ANY = 0xff

type Usb_raw_ep_caps_casted uint32

func (s Usb_raw_ep_caps_casted) GetType() Usb_raw_ep_caps {
	return Usb_raw_ep_caps{
		Type_control: (s & 0b000001) > 0,
		Type_iso:     (s & 0b000010) > 0,
		Type_bulk:    (s & 0b000100) > 0,
		Type_int:     (s & 0b001000) > 0,
		Dir_in:       (s & 0b010000) > 0,
		Dir_out:      (s & 0b100000) > 0,
	}
}

type Usb_raw_ep_caps struct {
	Type_control bool
	Type_iso     bool
	Type_bulk    bool
	Type_int     bool
	Dir_in       bool
	Dir_out      bool
}

type Usb_raw_ep_limits struct {
	Maxpacket_limit uint16
	Max_streams     uint16
	Reserved        uint32
}

type Usb_raw_ep_info struct {
	Name   [USB_RAW_EP_NAME_MAX]byte
	Addr   uint32
	Caps   Usb_raw_ep_caps_casted
	Limits Usb_raw_ep_limits
}

type Usb_raw_eps_info struct {
	Eps [USB_RAW_EPS_NUM_MAX]Usb_raw_ep_info
}

var (
	USB_RAW_IOCTL_INIT          = _IOW('U', 0, Sizeof(Usb_raw_init{}))
	USB_RAW_IOCTL_RUN           = _IO('U', 1)
	USB_RAW_IOCTL_EVENT_FETCH   = _IOR('U', 2, Sizeof(_ioctl_Usb_raw_event{}))
	USB_RAW_IOCTL_EP0_WRITE     = _IOW('U', 3, Sizeof(_ioctl_Usb_raw_ep_io{}))
	USB_RAW_IOCTL_EP0_READ      = _IOWR('U', 4, Sizeof(_ioctl_Usb_raw_ep_io{}))
	USB_RAW_IOCTL_EP_ENABLE     = _IOW('U', 5, 0x9) // Sizeof(Usb_endpoint_descriptor{})
	USB_RAW_IOCTL_EP_DISABLE    = _IOW('U', 6, Sizeof(uint32(0)))
	USB_RAW_IOCTL_EP_WRITE      = _IOW('U', 7, Sizeof(_ioctl_Usb_raw_ep_io{}))
	USB_RAW_IOCTL_EP_READ       = _IOWR('U', 8, Sizeof(_ioctl_Usb_raw_ep_io{}))
	USB_RAW_IOCTL_CONFIGURE     = _IO('U', 9)
	USB_RAW_IOCTL_VBUS_DRAW     = _IOW('U', 10, Sizeof(uint32(0)))
	USB_RAW_IOCTL_EPS_INFO      = _IOR('U', 11, Sizeof(Usb_raw_eps_info{}))
	USB_RAW_IOCTL_EP0_STALL     = _IO('U', 12)
	USB_RAW_IOCTL_EP_SET_HALT   = _IOW('U', 13, Sizeof(uint32(0)))
	USB_RAW_IOCTL_EP_CLEAR_HALT = _IOW('U', 14, Sizeof(uint32(0)))
	USB_RAW_IOCTL_EP_SET_WEDGE  = _IOW('U', 15, Sizeof(uint32(0)))
)

type Raw_gadget_interface_descriptor struct {
	Interface Usb_interface_descriptor
	Endpoints *Usb_endpoint_descriptor
}

type Raw_gadget_interface struct {
	Altsetting     *Raw_gadget_interface_descriptor
	Num_altsetting int
}

type Raw_gadget_config_descriptor struct {
	Config     Usb_config_descriptor
	Interfaces *Raw_gadget_interface
}

const PAGE_SIZE = 4096 // TODO: get from system
