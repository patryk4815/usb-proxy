package rawgadget

import "fmt"

// TODO: zalezy od wersji kernela: https://github.com/torvalds/linux/blob/master/include/uapi/linux/usb/ch9.h#L407

const VENDOR_REQ_OUT Usb_Req = 0x5b
const VENDOR_REQ_IN Usb_Req = 0x5c

/*
 * USB directions
 *
 * This bit flag is used in endpoint descriptors' BEndpointAddress field.
 * It's also one of three fields in control requests BRequestType.
 */

var USB_DIR_mapper = map[Usb_Dir]string{
	USB_DIR_OUT: "USB_DIR_OUT",
	USB_DIR_IN:  "USB_DIR_IN",
}

type Usb_Dir uint8

func (s Usb_Dir) String() string {
	r, ok := USB_DIR_mapper[s]
	if !ok {
		return fmt.Sprintf("unknown (%d)", s)
	}
	return r
}

const USB_DIR_OUT Usb_Dir = 0   /* to device */
const USB_DIR_IN Usb_Dir = 0x80 /* to host */

/*
 * USB types, the second of three BRequestType fields
 */

var USB_TYPE_mapper = map[Usb_Type]string{
	USB_TYPE_STANDARD: "USB_TYPE_STANDARD",
	USB_TYPE_CLASS:    "USB_TYPE_CLASS",
	USB_TYPE_VENDOR:   "USB_TYPE_VENDOR",
	USB_TYPE_RESERVED: "USB_TYPE_RESERVED",
}

type Usb_Type uint8

func (s Usb_Type) String() string {
	r, ok := USB_TYPE_mapper[s]
	if !ok {
		return fmt.Sprintf("unknown (%d)", s)
	}
	return r
}

const USB_TYPE_MASK Usb_Type = (0x03 << 5)
const USB_TYPE_STANDARD Usb_Type = (0x00 << 5)
const USB_TYPE_CLASS Usb_Type = (0x01 << 5)
const USB_TYPE_VENDOR Usb_Type = (0x02 << 5)
const USB_TYPE_RESERVED Usb_Type = (0x03 << 5)

/*
 * USB recipients, the third of three BRequestType fields
 */
const USB_RECIP_MASK = 0x1f
const USB_RECIP_DEVICE = 0x00
const USB_RECIP_INTERFACE = 0x01
const USB_RECIP_ENDPOINT = 0x02
const USB_RECIP_OTHER = 0x03

/* From Wireless USB 1.0 */
const USB_RECIP_PORT = 0x04
const USB_RECIP_RPIPE = 0x05

/*
 * Standard requests, for the BRequest field of a SETUP packet.
 *
 * These are qualified by the BRequestType field, so that for example
 * TYPE_CLASS or TYPE_VENDOR specific feature Flags could be retrieved
 * by a GET_STATUS request.
 */

var USB_REQ_mapper = map[Usb_Req]string{
	USB_REQ_GET_STATUS:        "USB_REQ_GET_STATUS",
	USB_REQ_CLEAR_FEATURE:     "USB_REQ_CLEAR_FEATURE",
	USB_REQ_SET_FEATURE:       "USB_REQ_SET_FEATURE",
	USB_REQ_SET_ADDRESS:       "USB_REQ_SET_ADDRESS",
	USB_REQ_GET_DESCRIPTOR:    "USB_REQ_GET_DESCRIPTOR",
	USB_REQ_SET_DESCRIPTOR:    "USB_REQ_SET_DESCRIPTOR",
	USB_REQ_GET_CONFIGURATION: "USB_REQ_GET_CONFIGURATION",
	USB_REQ_SET_CONFIGURATION: "USB_REQ_SET_CONFIGURATION",
	USB_REQ_GET_INTERFACE:     "USB_REQ_GET_INTERFACE",
	USB_REQ_SET_INTERFACE:     "USB_REQ_SET_INTERFACE",
	USB_REQ_SYNCH_FRAME:       "USB_REQ_SYNCH_FRAME",
	USB_REQ_SET_ENCRYPTION:    "USB_REQ_SET_ENCRYPTION",
	USB_REQ_GET_ENCRYPTION:    "USB_REQ_GET_ENCRYPTION",
	//USB_REQ_RPIPE_ABORT:         "USB_REQ_RPIPE_ABORT",
	USB_REQ_SET_HANDSHAKE: "USB_REQ_SET_HANDSHAKE",
	//USB_REQ_RPIPE_RESET:         "USB_REQ_RPIPE_RESET",
	USB_REQ_GET_HANDSHAKE:       "USB_REQ_GET_HANDSHAKE",
	USB_REQ_SET_CONNECTION:      "USB_REQ_SET_CONNECTION",
	USB_REQ_SET_SECURITY_DATA:   "USB_REQ_SET_SECURITY_DATA",
	USB_REQ_GET_SECURITY_DATA:   "USB_REQ_GET_SECURITY_DATA",
	USB_REQ_SET_WUSB_DATA:       "USB_REQ_SET_WUSB_DATA",
	USB_REQ_LOOPBACK_DATA_WRITE: "USB_REQ_LOOPBACK_DATA_WRITE",
	USB_REQ_LOOPBACK_DATA_READ:  "USB_REQ_LOOPBACK_DATA_READ",
	USB_REQ_SET_INTERFACE_DS:    "USB_REQ_SET_INTERFACE_DS",
	//USB_REQ_GET_PARTNER_PDO:     "USB_REQ_GET_PARTNER_PDO",
	//USB_REQ_GET_BATTERY_STATUS:  "USB_REQ_GET_BATTERY_STATUS",
	//USB_REQ_SET_PDO:             "USB_REQ_SET_PDO",
	//USB_REQ_GET_VDM:             "USB_REQ_GET_VDM",
	//USB_REQ_SEND_VDM:            "USB_REQ_SEND_VDM",
	USB_REQ_SET_SEL:         "USB_REQ_SET_SEL",
	USB_REQ_SET_ISOCH_DELAY: "USB_REQ_SET_ISOCH_DELAY",
}

type Usb_Req uint8

func (s Usb_Req) String() string {
	r, ok := USB_REQ_mapper[s]
	if !ok {
		return fmt.Sprintf("unknown (%d)", s)
	}
	return r
}

const USB_REQ_GET_STATUS Usb_Req = 0x00
const USB_REQ_CLEAR_FEATURE Usb_Req = 0x01
const USB_REQ_SET_FEATURE Usb_Req = 0x03
const USB_REQ_SET_ADDRESS Usb_Req = 0x05
const USB_REQ_GET_DESCRIPTOR Usb_Req = 0x06
const USB_REQ_SET_DESCRIPTOR Usb_Req = 0x07
const USB_REQ_GET_CONFIGURATION Usb_Req = 0x08
const USB_REQ_SET_CONFIGURATION Usb_Req = 0x09
const USB_REQ_GET_INTERFACE Usb_Req = 0x0A
const USB_REQ_SET_INTERFACE Usb_Req = 0x0B
const USB_REQ_SYNCH_FRAME Usb_Req = 0x0C
const USB_REQ_SET_ENCRYPTION Usb_Req = 0x0D
const USB_REQ_GET_ENCRYPTION Usb_Req = 0x0E

// const USB_REQ_RPIPE_ABORT = 0x0E
const USB_REQ_SET_HANDSHAKE Usb_Req = 0x0F

// const USB_REQ_RPIPE_RESET = 0x0F
const USB_REQ_GET_HANDSHAKE Usb_Req = 0x10
const USB_REQ_SET_CONNECTION Usb_Req = 0x11
const USB_REQ_SET_SECURITY_DATA Usb_Req = 0x12
const USB_REQ_GET_SECURITY_DATA Usb_Req = 0x13
const USB_REQ_SET_WUSB_DATA Usb_Req = 0x14
const USB_REQ_LOOPBACK_DATA_WRITE Usb_Req = 0x15
const USB_REQ_LOOPBACK_DATA_READ Usb_Req = 0x16
const USB_REQ_SET_INTERFACE_DS Usb_Req = 0x17

// const USB_REQ_GET_PARTNER_PDO = 20
// const USB_REQ_GET_BATTERY_STATUS = 21
// const USB_REQ_SET_PDO = 22
// const USB_REQ_GET_VDM = 23
// const USB_REQ_SEND_VDM = 24
const USB_REQ_SET_SEL Usb_Req = 0x30
const USB_REQ_SET_ISOCH_DELAY Usb_Req = 0x31

var USB_DT_mapper = map[Usb_DT]string{
	USB_DT_DEVICE:                 "USB_DT_DEVICE",
	USB_DT_CONFIG:                 "USB_DT_CONFIG",
	USB_DT_STRING:                 "USB_DT_STRING",
	USB_DT_INTERFACE:              "USB_DT_INTERFACE",
	USB_DT_ENDPOINT:               "USB_DT_ENDPOINT",
	USB_DT_DEVICE_QUALIFIER:       "USB_DT_DEVICE_QUALIFIER",
	USB_DT_OTHER_SPEED_CONFIG:     "USB_DT_OTHER_SPEED_CONFIG",
	USB_DT_INTERFACE_POWER:        "USB_DT_INTERFACE_POWER",
	USB_DT_OTG:                    "USB_DT_OTG",
	USB_DT_DEBUG:                  "USB_DT_DEBUG",
	USB_DT_INTERFACE_ASSOCIATION:  "USB_DT_INTERFACE_ASSOCIATION",
	USB_DT_SECURITY:               "USB_DT_SECURITY",
	USB_DT_KEY:                    "USB_DT_KEY",
	USB_DT_ENCRYPTION_TYPE:        "USB_DT_ENCRYPTION_TYPE",
	USB_DT_BOS:                    "USB_DT_BOS",
	USB_DT_DEVICE_CAPABILITY:      "USB_DT_DEVICE_CAPABILITY",
	USB_DT_WIRELESS_ENDPOINT_COMP: "USB_DT_WIRELESS_ENDPOINT_COMP",
	USB_DT_WIRE_ADAPTER:           "USB_DT_WIRE_ADAPTER",
	USB_DT_RPIPE:                  "USB_DT_RPIPE",
	USB_DT_CS_RADIO_CONTROL:       "USB_DT_CS_RADIO_CONTROL",
	USB_DT_PIPE_USAGE:             "USB_DT_PIPE_USAGE",
	USB_DT_SS_ENDPOINT_COMP:       "USB_DT_SS_ENDPOINT_COMP",
	USB_DT_SSP_ISOC_ENDPOINT_COMP: "USB_DT_SSP_ISOC_ENDPOINT_COMP",
}

type Usb_DT uint8

func (s Usb_DT) String() string {
	r, ok := USB_DT_mapper[s]
	if !ok {
		return fmt.Sprintf("unknown (%d)", s)
	}
	return r
}

const USB_DT_DEVICE Usb_DT = 0x01
const USB_DT_CONFIG Usb_DT = 0x02
const USB_DT_STRING Usb_DT = 0x03
const USB_DT_INTERFACE Usb_DT = 0x04
const USB_DT_ENDPOINT Usb_DT = 0x05
const USB_DT_DEVICE_QUALIFIER Usb_DT = 0x06
const USB_DT_OTHER_SPEED_CONFIG Usb_DT = 0x07
const USB_DT_INTERFACE_POWER Usb_DT = 0x08
const USB_DT_OTG Usb_DT = 0x09
const USB_DT_DEBUG Usb_DT = 0x0a
const USB_DT_INTERFACE_ASSOCIATION Usb_DT = 0x0b
const USB_DT_SECURITY Usb_DT = 0x0c
const USB_DT_KEY Usb_DT = 0x0d
const USB_DT_ENCRYPTION_TYPE Usb_DT = 0x0e
const USB_DT_BOS Usb_DT = 0x0f
const USB_DT_DEVICE_CAPABILITY Usb_DT = 0x10
const USB_DT_WIRELESS_ENDPOINT_COMP Usb_DT = 0x11
const USB_DT_WIRE_ADAPTER Usb_DT = 0x21
const USB_DT_RPIPE Usb_DT = 0x22
const USB_DT_CS_RADIO_CONTROL Usb_DT = 0x23
const USB_DT_PIPE_USAGE Usb_DT = 0x24
const USB_DT_SS_ENDPOINT_COMP Usb_DT = 0x30
const USB_DT_SSP_ISOC_ENDPOINT_COMP Usb_DT = 0x31

type Usb_endpoint_descriptor struct {
	BLength          uint8
	BDescriptorType  uint8
	BEndpointAddress uint8
	BmAttributes     uint8
	WMaxPacketSize   uint16
	BInterval        uint8
	BRefresh         uint8
	BSynchAddress    uint8
}

type Usb_ctrlrequest struct {
	BRequestType uint8  `struc:"uint8,little"`
	BRequest     uint8  `struc:"uint8,little"`
	WValue       uint16 `struc:"uint16,little"`
	WIndex       uint16 `struc:"uint16,little"`
	WLength      uint16 `struc:"uint16,little"`
}

type Usb_config_descriptor struct {
	BLength             uint8
	BDescriptorType     uint8
	WTotalLength        uint16
	BNumInterfaces      uint8
	BConfigurationValue uint8
	IConfiguration      uint8
	BmAttributes        uint8
	BMaxPower           uint8
}

type Usb_interface_descriptor struct {
	BLength            uint8
	BDescriptorType    uint8
	BInterfaceNumber   uint8
	BAlternateSetting  uint8
	BNumEndpoints      uint8
	BInterfaceClass    uint8
	BInterfaceSubClass uint8
	BInterfaceProtocol uint8
	IInterface         uint8
}
