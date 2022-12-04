package hostproxy

import (
	"context"
	"fmt"
	"github.com/google/gousb"
	"github.com/patryk4815/usb-proxy/pkg/ctxproxy"
	"github.com/patryk4815/usb-proxy/pkg/rawproxy"
	log "github.com/sirupsen/logrus"
	"io"
	"time"
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
}

func (s *XX_Host) SetRawProxy(h IRawDevice) {
	s.rawDevice = h
}

func (s *XX_Host) ReadContext(ctx context.Context, output []byte) (int, error) {
	var readed int
	event := <-s.rawDevice.GetChWriter()
	if info := ctxproxy.CtxEP0Data(ctx); info != nil {
		info.Event = event.Event
		info.Close = func(err error) {
			event.ChDone <- err
		}
	}

	log.WithField("bytes", event.Event.RawCtrlReq.WLength).Debug("ep0: IN(raw) copying bytes")
	defer func() {
		log.WithField("bytes", readed).WithField("input", event.Event.RawCtrlReq.WLength).Debug("ep0: IN(raw) transferred bytes")
	}()

	if len(output) < int(event.Event.RawCtrlReq.WLength) {
		err := fmt.Errorf("output short len has=%d needed=%d", len(output), event.Event.RawCtrlReq.WLength)
		panic(err)
	}

	out := make([]byte, event.Event.RawCtrlReq.WLength)
	n, err := s.Dev.Control(event.Event.RawCtrlReq.BRequestType, event.Event.RawCtrlReq.BRequest, event.Event.RawCtrlReq.WValue, event.Event.RawCtrlReq.WIndex, out)
	readed = n
	if n < 0 && err == nil {
		return 0, fmt.Errorf("ReadContext control minus value = %d", n)
	}
	if n >= 0 {
		copy(output, out[:n])
	}
	return n, err // TODO: wrap error
}

func (s *XX_Host) WriteContext(ctx context.Context, input []byte) (int, error) {
	var readed int
	info := ctxproxy.CtxEP0Data(ctx)
	if info == nil {
		panic(fmt.Errorf("context should have info"))
	}

	log.WithField("bytes", len(input)).Debug("ep0: OUT(host) writing bytes")
	defer func() {
		log.WithField("bytes", readed).WithField("input", len(input)).Debug("ep0: OUT(host) transferred bytes")
	}()

	out := make([]byte, len(input))
	copy(out, input)

	n, err := s.Dev.Control(info.Event.RawCtrlReq.BRequestType, info.Event.RawCtrlReq.BRequest, info.Event.RawCtrlReq.WValue, info.Event.RawCtrlReq.WIndex, out)
	readed = n
	if n < 0 && err == nil {
		return 0, fmt.Errorf("WriteContext control minus value = %d", n)
	}
	return n, err // TODO: wrap error
}

func (s *XX_Host) GetConfig(wValue uint16) (*gousb.Config, error) {
	return s.Dev.Config(int(wValue))
}

func (s *XX_Host) Close() {
	s.Dev.Close()
	s.ctx.Close()
}

func (s *XX_Host) Open(vid, pid gousb.ID) {
	ctx := gousb.NewContext()
	s.ctx = ctx

	// TODO: auto select when usb plug-in?

	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		log.Printf("Scan OpenDevices: %s\n", desc.String())
		return desc.Vendor == vid && desc.Product == pid
	})

	if err != nil {
		log.Fatalf("OpenDevices(): %v", err)
	}
	if len(devs) == 0 {
		log.Fatalf("no devices found matching VID %s and PID %s", vid, pid)
	}

	dev := devs[0]
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
	s.Dev.ControlTimeout = time.Second * 999
}
