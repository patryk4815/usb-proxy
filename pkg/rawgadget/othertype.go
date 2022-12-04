package rawgadget

import (
	"bytes"
	"github.com/lunixbochs/struc"
)

type UsbEventCtrl struct {
	RawEventType   Usb_raw_event_type
	Type           Usb_Type
	Request        Usb_Req
	DescriptorType Usb_DT
	Direction      Usb_Dir

	RawCtrlReq Usb_ctrlrequest
}

func (s *UsbEventCtrl) UpdateFrom(newobj *Usb_raw_event) error {
	s.RawEventType = newobj.Type

	tmp := &Usb_ctrlrequest{}
	err := struc.Unpack(bytes.NewReader(newobj.Data), tmp)
	if err != nil {
		return err
	}
	s.RawCtrlReq = *tmp

	s.Type = Usb_Type(tmp.BRequestType & uint8(USB_TYPE_MASK))
	s.Request = Usb_Req(tmp.BRequest)
	s.DescriptorType = Usb_DT(tmp.WValue >> 8)

	isIn := (uint8(tmp.BRequestType) & uint8(USB_DIR_IN)) > 0
	if isIn {
		s.Direction = USB_DIR_IN
	} else {
		s.Direction = USB_DIR_OUT
	}
	return nil
}
