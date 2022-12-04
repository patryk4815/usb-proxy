package ctxproxy

import "github.com/patryk4815/usb-proxy/pkg/rawgadget"

type XX_RawReaderWriter struct {
	ChDone chan error
	Event  rawgadget.UsbEventCtrl
}

type InfoCtxEP0Data struct {
	Event rawgadget.UsbEventCtrl
	Close func(err error)
}
