package main

const UDC_NAME_LENGTH_MAX uint = 128

type usb_raw_init struct {
	driver_name [UDC_NAME_LENGTH_MAX]byte
	device_name [UDC_NAME_LENGTH_MAX]byte
	speed       byte
}

type usb_raw_event_type uint32

const (
	USB_RAW_EVENT_INVALID usb_raw_event_type = 0
	USB_RAW_EVENT_CONNECT usb_raw_event_type = 1
	USB_RAW_EVENT_CONTROL usb_raw_event_type = 2
)

type usb_raw_event struct {
	xtype  usb_raw_event_type
	length uint32
}

type usb_raw_ep_io struct {
	ep     uint16
	flags  uint16
	length uint32
}

const USB_RAW_EPS_NUM_MAX = 30
const USB_RAW_EP_NAME_MAX = 16
const USB_RAW_EP_ADDR_ANY = 0xff

type usb_raw_ep_caps_casted uint32

func (s usb_raw_ep_caps_casted) GetType() usb_raw_ep_caps {
	return usb_raw_ep_caps{
		type_control: (s & 0b000001) > 0,
		type_iso:     (s & 0b000010) > 0,
		type_bulk:    (s & 0b000100) > 0,
		type_int:     (s & 0b001000) > 0,
		dir_in:       (s & 0b010000) > 0,
		dir_out:      (s & 0b100000) > 0,
	}
}

type usb_raw_ep_caps struct {
	type_control bool
	type_iso     bool
	type_bulk    bool
	type_int     bool
	dir_in       bool
	dir_out      bool
}

type usb_raw_ep_limits struct {
	maxpacket_limit uint16
	max_streams     uint16
	reserved        uint32
}

type usb_raw_ep_info struct {
	name   [USB_RAW_EP_NAME_MAX]byte
	addr   uint32
	caps   usb_raw_ep_caps_casted
	limits usb_raw_ep_limits
}

type usb_raw_eps_info struct {
	eps [USB_RAW_EPS_NUM_MAX]usb_raw_ep_info
}

var (
	USB_RAW_IOCTL_INIT          = _IOW('U', 0, sizeOf(usb_raw_init{}))
	USB_RAW_IOCTL_RUN           = _IO('U', 1)
	USB_RAW_IOCTL_EVENT_FETCH   = _IOR('U', 2, sizeOf(usb_raw_event{}))
	USB_RAW_IOCTL_EP0_WRITE     = _IOW('U', 3, sizeOf(usb_raw_ep_io{}))
	USB_RAW_IOCTL_EP0_READ      = _IOWR('U', 4, sizeOf(usb_raw_ep_io{}))
	USB_RAW_IOCTL_EP_ENABLE     = _IOW('U', 5, 0x9) // sizeOf(usb_endpoint_descriptor{})
	USB_RAW_IOCTL_EP_DISABLE    = _IOW('U', 6, sizeOf(uint32(0)))
	USB_RAW_IOCTL_EP_WRITE      = _IOW('U', 7, sizeOf(usb_raw_ep_io{}))
	USB_RAW_IOCTL_EP_READ       = _IOWR('U', 8, sizeOf(usb_raw_ep_io{}))
	USB_RAW_IOCTL_CONFIGURE     = _IO('U', 9)
	USB_RAW_IOCTL_VBUS_DRAW     = _IOW('U', 10, sizeOf(uint32(0)))
	USB_RAW_IOCTL_EPS_INFO      = _IOR('U', 11, sizeOf(usb_raw_eps_info{}))
	USB_RAW_IOCTL_EP0_STALL     = _IO('U', 12)
	USB_RAW_IOCTL_EP_SET_HALT   = _IOW('U', 13, sizeOf(uint32(0)))
	USB_RAW_IOCTL_EP_CLEAR_HALT = _IOW('U', 14, sizeOf(uint32(0)))
	USB_RAW_IOCTL_EP_SET_WEDGE  = _IOW('U', 15, sizeOf(uint32(0)))
)

type raw_gadget_interface_descriptor struct {
	xinterface usb_interface_descriptor
	endpoints  *usb_endpoint_descriptor
}

type raw_gadget_interface struct {
	altsetting     *raw_gadget_interface_descriptor
	num_altsetting int
}

type raw_gadget_config_descriptor struct {
	config     usb_config_descriptor
	interfaces *raw_gadget_interface
}

const EP_MAX_PACKET_CONTROL = 512
const EP_MAX_PACKET_TRANSFER = 512

type usb_raw_control_event struct {
	inner usb_raw_event
	ctrl  usb_ctrlrequest
}

type usb_raw_control_io struct {
	inner usb_raw_ep_io
	data  [EP_MAX_PACKET_CONTROL]byte
}

type usb_raw_transfer_io struct {
	inner usb_raw_ep_io
	data  [EP_MAX_PACKET_TRANSFER]byte
}
