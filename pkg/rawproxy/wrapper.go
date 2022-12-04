package rawproxy

import (
	"context"
	"fmt"
	_ "github.com/lunixbochs/struc"
	"github.com/patryk4815/usb-proxy/pkg/ctxproxy"
	"github.com/patryk4815/usb-proxy/pkg/rawgadget"
	log "github.com/sirupsen/logrus"
)

type WrapperUsbEP0 struct {
	Self *XX_www
}

func (s *WrapperUsbEP0) ReadContext(ctx context.Context, output []byte) (int, error) {
	var readed int
	event := <-s.Self.chReader

	log.WithField("bytes", event.Event.RawCtrlReq.WLength).Debug("ep0: IN(raw) copying bytes")
	defer func() {
		log.WithField("bytes", readed).WithField("input", event.Event.RawCtrlReq.WLength).Debug("ep0: IN(raw) transferred bytes")
	}()

	if len(output) < int(event.Event.RawCtrlReq.WLength) {
		err := fmt.Errorf("output short len has=%d needed=%d", len(output), event.Event.RawCtrlReq.WLength)
		panic(err)
	}

	if info := ctxproxy.CtxEP0Data(ctx); info != nil {
		info.Event = event.Event
		info.Close = func(err error) {
			event.ChDone <- err
		}
	}

	ep := &rawgadget.Usb_raw_ep_io{}
	ep.Data = make([]byte, event.Event.RawCtrlReq.WLength)

	n, err := s.Self.otgDevice.EP0Read(ep)
	readed = n
	if n < 0 && err == nil {
		return 0, fmt.Errorf("WrapperUsbEP0 ReadContext control minus value = %d", n)
	}
	if n >= 0 {
		copy(output, ep.Data[:n])
	}

	return n, err // TODO: wrap err
}

func (s *WrapperUsbEP0) WriteContext(ctx context.Context, input []byte) (int, error) {
	var readed int

	log.WithField("bytes", len(input)).Debug("ep0: OUT(raw) writing bytes")
	defer func() {
		log.WithField("bytes", readed).WithField("input", len(input)).Debug("ep0: OUT(raw) transferred bytes")
	}()

	ep := &rawgadget.Usb_raw_ep_io{}
	ep.Data = make([]byte, len(input))
	copy(ep.Data, input)

	n, err := s.Self.otgDevice.EP0Write(ep)
	readed = n
	if n < 0 && err == nil {
		return 0, fmt.Errorf("WrapperUsbEP0 WriteContext control minus value = %d", n)
	}
	return n, err // TODO: wrap err
}

type WrapperUsbEPX struct {
	Self        *XX_www
	EpNum       int
	ReadedBytes int
	WritedBytes int
}

func (s *WrapperUsbEPX) Read(input []byte) (int, error) {
	var readed int
	pkt := &rawgadget.Usb_raw_ep_io{}
	pkt.Ep = uint16(s.EpNum)
	pkt.Data = make([]byte, len(input))

	log.WithField("bytes", len(input)).WithField("ep", pkt.Ep).Debug("epX: IN(raw) copying bytes")
	defer func() {
		logger := log.WithField("bytes", readed).WithField("input", len(input)).WithField("ep", pkt.Ep)
		logger = logger.WithField("whole_r", float64(s.ReadedBytes)/1024/1024)
		logger = logger.WithField("buf", input)
		logger.Debug("epX: IN(raw) transferred bytes")
	}()

	n, err := s.Self.otgDevice.EPRead(pkt)
	readed = n
	if n < 0 && err == nil {
		return 0, fmt.Errorf("WrapperUsbEPX Read control minus value = %d", n)
	}
	if n >= 0 {
		copy(input, pkt.Data[:n])
		s.ReadedBytes += n
	}
	return n, err
}

func (s *WrapperUsbEPX) Write(input []byte) (int, error) {
	var readed int
	pkt := &rawgadget.Usb_raw_ep_io{}
	pkt.Ep = uint16(s.EpNum)
	pkt.Data = make([]byte, len(input))
	copy(pkt.Data, input)

	logger := log.WithField("bytes", len(input)).WithField("ep", pkt.Ep)
	logger = logger.WithField("whole_w", float64(s.WritedBytes)/1024/1024)
	logger = logger.WithField("buf", input)
	logger.Debug("epX: OUT(raw) writing bytes")
	defer func() {
		log.WithField("bytes", readed).WithField("input", len(input)).WithField("ep", pkt.Ep).Debug("epX: OUT(raw) transferred bytes")
	}()

	n, err := s.Self.otgDevice.EPWrite(pkt)
	readed = n
	if n < 0 && err == nil {
		return 0, fmt.Errorf("WrapperUsbEPX Write control minus value = %d", n)
	}
	s.WritedBytes += n

	return n, err
}
