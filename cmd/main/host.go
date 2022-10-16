package main

import (
	"github.com/google/gousb"
	"io"
	"log"
)

type xHost struct {
	io.Reader
	io.Writer

	ctx *gousb.Context
	dev *gousb.Device

	lolRaw *xRaw
}

func NewHost() *xHost {
	return &xHost{}
}

func (s *xHost) Read(output []byte) (int, error) {
	defer func() {
		log.Printf("xHost Read-End\n")
	}()
	event := <-s.lolRaw.chWriter
	eventData := event.event
	log.Printf("xHost Read: len=%d, %+v\n", len(output), eventData.ctrl)
	defer close(event.chDone)

	// 	int result = libusb_control_transfer(dev_handle,
	//					setup_packet->bRequestType, setup_packet->bRequest,
	//					setup_packet->wValue, setup_packet->wIndex, *dataptr,
	//					setup_packet->wLength, timeout);
	out := make([]byte, eventData.ctrl.wLength)

	n, err := s.dev.Control(eventData.ctrl.bRequestType, eventData.ctrl.bRequest, eventData.ctrl.wValue, eventData.ctrl.wIndex, out)
	log.Printf("xHost Read-result n=%d, err=%#v\n", n, err)
	if n <= 0 {
		//s.lolRaw.Ep0Stall()
		//s.dev.Reset()
		//s.dev.SetAutoDetach(false)
		panic(err)
		return n, err
	}
	copy(output, out)

	log.Printf("[dump] readHost: bRequestType=%d bRequest=%d wValue=%x wIndex=%x data=%q\n",
		eventData.ctrl.bRequestType, eventData.ctrl.bRequest,
		eventData.ctrl.wValue, eventData.ctrl.wIndex,
		out[:n])

	return n, nil
}

func (s *xHost) Write(input []byte) (int, error) {
	defer func() {
		log.Printf("xHost Write-End\n")
	}()
	eventData := s.lolRaw.lastEvent
	log.Printf("xHost Write: %+v\n", eventData.ctrl)

	out := make([]byte, eventData.ctrl.wLength)
	nwrited := copy(out, input)

	log.Printf("[dump] writeHost: bRequestType=%d bRequest=%d wValue=%x wIndex=%x data=%q\n",
		eventData.ctrl.bRequestType, eventData.ctrl.bRequest,
		eventData.ctrl.wValue, eventData.ctrl.wIndex,
		out[:nwrited])

	n, err := s.dev.Control(eventData.ctrl.bRequestType, eventData.ctrl.bRequest, eventData.ctrl.wValue, eventData.ctrl.wIndex, out)
	log.Printf("xHost Write-result n=%d, err=%#v\n", n, err)
	if n < 0 {
		panic(err)
		return n, err
	}

	return n, nil
}

func (s *xHost) Close() {
	s.dev.Close()
	s.ctx.Close()
}

func (s *xHost) Open() {
	ctx := gousb.NewContext()
	s.ctx = ctx

	// 1038:160e
	vid, pid := gousb.ID(0x0930), gousb.ID(0x6544)
	//vid, pid := gousb.ID(0x04e8), gousb.ID(0x60b3)
	//vid, pid := gousb.ID(0x090c), gousb.ID(0x1000)
	//vid, pid := gousb.ID(0x1038), gousb.ID(0x160e)
	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		log.Printf("Scan OpenDevices: %#v\n", desc)
		return desc.Vendor == vid && desc.Product == pid
	})

	if err != nil {
		log.Fatalf("OpenDevices(): %v", err)
	}
	if len(devs) == 0 {
		log.Fatalf("no devices found matching VID %s and PID %s", vid, pid)
	}

	dev := devs[0]
	//{
	//	out, _ := dev.Product()
	//	log.Printf("Product: %s\n", out)
	//}

	{
		err := dev.SetAutoDetach(true)
		log.Printf("SetAutoDetach: %#v\n", err)

	}

	{
		err := dev.Reset()
		log.Printf("Reset: %#v\n", err)
	}

	{
		out, err := dev.GetStringDescriptor(0)
		log.Printf("GetStringDescriptor: %q, %#v\n", out, err)
	}

	s.dev = dev
	//s.dev.ControlTimeout = time.Second
}
