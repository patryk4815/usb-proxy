package rawproxy

import (
	_ "github.com/lunixbochs/struc"
	"github.com/patryk4815/usb-proxy/pkg/rawgadget"
	"log"
)

type WrapperUsbRawControlIO struct {
	Self *XX_www
}

func (s *WrapperUsbRawControlIO) Read(output []byte) (int, error) {
	defer func() {
		log.Printf("XX_www Read-End\n")
	}()

	event := <-s.Self.chReader
	log.Printf("XX_www Read lenEvent=%d lenOutput=%d\n", event.Event.RawCtrlReq.WLength, len(output))
	s.Self.lastEvent = event.Event
	defer close(event.ChDone)

	ep := &rawgadget.Usb_raw_ep_io{}
	ep.Data = make([]byte, event.Event.RawCtrlReq.WLength)

	n, err := s.Self.otgDevice.EP0Read(ep)
	log.Printf("XX_www Read-ret len=%d err=%#v\n", n, err)

	copy(output, ep.Data[:n])
	log.Printf("[dump] readOTG: Ep=%d Data=%q\n", ep.Ep, output)

	return n, err
}

func (s *WrapperUsbRawControlIO) Write(input []byte) (int, error) {
	defer func() {
		log.Printf("XX_www Write-End\n")
	}()
	log.Printf("XX_www Write len=%d\n", len(input))

	ep := &rawgadget.Usb_raw_ep_io{}
	ep.Data = make([]byte, len(input))
	copy(ep.Data[:], input)

	log.Printf("[dump] writeOTG: Ep=%d Data=%q\n",
		ep.Ep,
		input)

	n, err := s.Self.otgDevice.EP0Write(ep)
	log.Printf("XX_www Write-ret len=%d err=%#v\n", n, err)

	return n, err
}

type WrapperUsbRawTransferIO struct {
	Self  *XX_www
	EpNum int
}

func (s *WrapperUsbRawTransferIO) Read(p []byte) (int, error) {
	pkt := &rawgadget.Usb_raw_ep_io{}
	pkt.Ep = uint16(s.EpNum)
	pkt.Data = make([]byte, len(p))

	log.Printf("WrapperUsbRawTransferIO Read len=%d\n", len(p))

	n, err := s.Self.otgDevice.EPRead(pkt)

	log.Printf("WrapperUsbRawTransferIO Read-result len=%d, n=%d, err=%v\n", len(p), n, err)

	if err != nil {
		panic(err)
	}
	copy(p, pkt.Data[:n])

	return n, nil
}

func (s *WrapperUsbRawTransferIO) Write(p []byte) (int, error) {
	pkt := &rawgadget.Usb_raw_ep_io{}
	pkt.Ep = uint16(s.EpNum)
	pkt.Data = make([]byte, len(p))
	copy(pkt.Data[:], p)

	log.Printf("WrapperUsbRawTransferIO Write len=%d\n", len(p))

	n, err := s.Self.otgDevice.EPWrite(pkt)

	log.Printf("WrapperUsbRawTransferIO Write-result len=%d n=%d err=%#v\n", len(p), n, err)
	if err != nil {
		panic(err)
	}
	return n, nil
}
