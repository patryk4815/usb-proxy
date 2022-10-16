package main

// TODO: zalezy od wersji kernela: https://github.com/torvalds/linux/blob/master/include/uapi/linux/usb/ch9.h#L407

const VENDOR_REQ_OUT = 0x5b
const VENDOR_REQ_IN = 0x5c

/*
 * USB directions
 *
 * This bit flag is used in endpoint descriptors' bEndpointAddress field.
 * It's also one of three fields in control requests bRequestType.
 */
const USB_DIR_OUT = 0   /* to device */
const USB_DIR_IN = 0x80 /* to host */

/*
 * USB types, the second of three bRequestType fields
 */

var USB_TYPE_mapper = map[int]string{
	USB_TYPE_STANDARD: "USB_TYPE_STANDARD",
	USB_TYPE_CLASS:    "USB_TYPE_CLASS",
	USB_TYPE_VENDOR:   "USB_TYPE_VENDOR",
	USB_TYPE_RESERVED: "USB_TYPE_RESERVED",
}

const USB_TYPE_MASK = (0x03 << 5)
const USB_TYPE_STANDARD = (0x00 << 5)
const USB_TYPE_CLASS = (0x01 << 5)
const USB_TYPE_VENDOR = (0x02 << 5)
const USB_TYPE_RESERVED = (0x03 << 5)

/*
 * USB recipients, the third of three bRequestType fields
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
 * Standard requests, for the bRequest field of a SETUP packet.
 *
 * These are qualified by the bRequestType field, so that for example
 * TYPE_CLASS or TYPE_VENDOR specific feature flags could be retrieved
 * by a GET_STATUS request.
 */

var USB_REQ_mapper = map[int]string{
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

const USB_REQ_GET_STATUS = 0x00
const USB_REQ_CLEAR_FEATURE = 0x01
const USB_REQ_SET_FEATURE = 0x03
const USB_REQ_SET_ADDRESS = 0x05
const USB_REQ_GET_DESCRIPTOR = 0x06
const USB_REQ_SET_DESCRIPTOR = 0x07
const USB_REQ_GET_CONFIGURATION = 0x08
const USB_REQ_SET_CONFIGURATION = 0x09
const USB_REQ_GET_INTERFACE = 0x0A
const USB_REQ_SET_INTERFACE = 0x0B
const USB_REQ_SYNCH_FRAME = 0x0C
const USB_REQ_SET_ENCRYPTION = 0x0D
const USB_REQ_GET_ENCRYPTION = 0x0E

//const USB_REQ_RPIPE_ABORT = 0x0E
const USB_REQ_SET_HANDSHAKE = 0x0F

//const USB_REQ_RPIPE_RESET = 0x0F
const USB_REQ_GET_HANDSHAKE = 0x10
const USB_REQ_SET_CONNECTION = 0x11
const USB_REQ_SET_SECURITY_DATA = 0x12
const USB_REQ_GET_SECURITY_DATA = 0x13
const USB_REQ_SET_WUSB_DATA = 0x14
const USB_REQ_LOOPBACK_DATA_WRITE = 0x15
const USB_REQ_LOOPBACK_DATA_READ = 0x16
const USB_REQ_SET_INTERFACE_DS = 0x17

//const USB_REQ_GET_PARTNER_PDO = 20
//const USB_REQ_GET_BATTERY_STATUS = 21
//const USB_REQ_SET_PDO = 22
//const USB_REQ_GET_VDM = 23
//const USB_REQ_SEND_VDM = 24
const USB_REQ_SET_SEL = 0x30
const USB_REQ_SET_ISOCH_DELAY = 0x31

var USB_DT_mapper = map[int]string{
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

const USB_DT_DEVICE = 0x01
const USB_DT_CONFIG = 0x02
const USB_DT_STRING = 0x03
const USB_DT_INTERFACE = 0x04
const USB_DT_ENDPOINT = 0x05
const USB_DT_DEVICE_QUALIFIER = 0x06
const USB_DT_OTHER_SPEED_CONFIG = 0x07
const USB_DT_INTERFACE_POWER = 0x08
const USB_DT_OTG = 0x09
const USB_DT_DEBUG = 0x0a
const USB_DT_INTERFACE_ASSOCIATION = 0x0b
const USB_DT_SECURITY = 0x0c
const USB_DT_KEY = 0x0d
const USB_DT_ENCRYPTION_TYPE = 0x0e
const USB_DT_BOS = 0x0f
const USB_DT_DEVICE_CAPABILITY = 0x10
const USB_DT_WIRELESS_ENDPOINT_COMP = 0x11
const USB_DT_WIRE_ADAPTER = 0x21
const USB_DT_RPIPE = 0x22
const USB_DT_CS_RADIO_CONTROL = 0x23
const USB_DT_PIPE_USAGE = 0x24
const USB_DT_SS_ENDPOINT_COMP = 0x30
const USB_DT_SSP_ISOC_ENDPOINT_COMP = 0x31

type usb_endpoint_descriptor struct {
	bLength          uint8
	bDescriptorType  uint8
	bEndpointAddress uint8
	bmAttributes     uint8
	wMaxPacketSize   uint16
	bInterval        uint8
	bRefresh         uint8
	bSynchAddress    uint8
}

type usb_ctrlrequest struct {
	bRequestType uint8
	bRequest     uint8
	wValue       uint16
	wIndex       uint16
	wLength      uint16
}

type usb_config_descriptor struct {
	bLength             uint8
	bDescriptorType     uint8
	wTotalLength        uint16
	bNumInterfaces      uint8
	bConfigurationValue uint8
	iConfiguration      uint8
	bmAttributes        uint8
	bMaxPower           uint8
}

type usb_interface_descriptor struct {
	bLength            uint8
	bDescriptorType    uint8
	bInterfaceNumber   uint8
	bAlternateSetting  uint8
	bNumEndpoints      uint8
	bInterfaceClass    uint8
	bInterfaceSubClass uint8
	bInterfaceProtocol uint8
	iInterface         uint8
}
