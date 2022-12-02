package hostproxy

import (
	"github.com/google/gousb"
	"github.com/patryk4815/usb-proxy/pkg/rawproxy"
	"io"
	"log"
)

type XX_Host struct {
	io.Reader
	io.Writer

	ctx *gousb.Context
	Dev *gousb.Device

	rawDevice IRawDevice
}

func New() *XX_Host {
	return &XX_Host{}
}

type IRawDevice interface {
	GetChWriter() <-chan rawproxy.XX_RawReaderWriter
	GetLastEvent() rawproxy.UsbEventCtrl
}

func (s *XX_Host) SetRawProxy(h IRawDevice) {
	s.rawDevice = h
}

func (s *XX_Host) Read(output []byte) (int, error) {
	defer func() {
		log.Printf("XX_Host Read-End\n")
	}()
	event := <-s.rawDevice.GetChWriter()
	eventData := event.Event
	log.Printf("XX_Host Read: len=%d, %+v\n", len(output), eventData)
	defer close(event.ChDone)

	out := make([]byte, eventData.RawCtrlReq.WLength) // TODO: output len???? WTF, moze jednak nie
	n, err := s.Dev.Control(eventData.RawCtrlReq.BRequestType, eventData.RawCtrlReq.BRequest, eventData.RawCtrlReq.WValue, eventData.RawCtrlReq.WIndex, out)

	log.Printf("XX_Host Read-result n=%d, err=%#v\n", n, err)
	if n <= 0 {
		panic(err)
		return n, err
	}
	copy(output, out)

	log.Printf("[dump] readHost: info=%+v data=%q\n", eventData, out[:n])

	return n, nil
}

func (s *XX_Host) Write(input []byte) (int, error) {
	defer func() {
		log.Printf("XX_Host Write-End\n")
	}()

	eventData := s.rawDevice.GetLastEvent()
	log.Printf("XX_Host Write: len=%d cap=%d %+v\n", len(input), cap(input), eventData)

	// TODO: ogarnac ten length tutaj
	//out := make([]byte, eventData.RawCtrlReq.WLength)
	//nwrited := copy(out, input)
	//_ = nwrited
	out := make([]byte, len(input))
	copy(out, input)

	n, err := s.Dev.Control(eventData.RawCtrlReq.BRequestType, eventData.RawCtrlReq.BRequest, eventData.RawCtrlReq.WValue, eventData.RawCtrlReq.WIndex, out)
	log.Printf("XX_Host Write-result n=%d, err=%#v\n", n, err)

	if n < 0 {
		panic(err)
		return n, err
	}

	return n, nil
}

func (s *XX_Host) GetConfig(wValue uint16) (*gousb.Config, error) {
	return s.Dev.Config(int(wValue))
}

func (s *XX_Host) Close() {
	s.Dev.Close()
	s.ctx.Close()
}

func (s *XX_Host) Open() {
	ctx := gousb.NewContext()
	s.ctx = ctx

	// TODO: select device
	// TODO: auto select when usb plug-in?

	// 1038:160e
	// 090c:1000
	// 0930:6544

	//2022/12/02 02:55:22 Scan OpenDevices: 1.1: 1d6b:0002 (available configs: [1])
	//2022/12/02 02:55:22 Scan OpenDevices: 3.1: 1d6b:0003 (available configs: [1])
	//2022/12/02 02:55:22 Scan OpenDevices: 2.4: 0930:6544 (available configs: [1])
	//2022/12/02 02:55:22 Scan OpenDevices: 2.2: 2109:3431 (available configs: [1])
	//2022/12/02 02:55:22 Scan OpenDevices: 2.1: 1d6b:0002 (available configs: [1])

	vid, pid := gousb.ID(0x0930), gousb.ID(0x6544) //bialy
	//vid, pid := gousb.ID(0x090c), gousb.ID(0x1000) // z rzadu
	//vid, pid := gousb.ID(0x04e8), gousb.ID(0x60b3)
	//vid, pid := gousb.ID(0x090c), gousb.ID(0x1000)
	//vid, pid := gousb.ID(0x1038), gousb.ID(0x160e)
	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		log.Printf("Scan OpenDevices: %s\n", desc.String())
		return desc.Vendor == vid && desc.Product == pid
		//return desc.Bus == 2 && desc.Address == 4
	})

	if err != nil {
		log.Fatalf("OpenDevices(): %v", err)
	}
	if len(devs) == 0 {
		log.Fatalf("no devices found matching VID %s and PID %s", vid, pid)
	}

	dev := devs[0]
	//{
	//	out, _ := Dev.Product()
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

	s.Dev = dev
	//s.Dev.ControlTimeout = time.Second
}
