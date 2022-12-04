package main

import (
	"context"
	"github.com/patryk4815/usb-proxy/pkg/ctxproxy"
	"github.com/patryk4815/usb-proxy/pkg/rawgadget"
	log "github.com/sirupsen/logrus"
)

type ReaderCtx interface {
	ReadContext(ctx context.Context, p []byte) (n int, err error)
}

type WriterCtx interface {
	WriteContext(ctx context.Context, p []byte) (n int, err error)
}

func customCopy(dst WriterCtx, src ReaderCtx) {
	for {
		customCopySingle(dst, src)
	}
}

func customCopySingle(dst WriterCtx, src ReaderCtx) {
	ctxInfo := &ctxproxy.InfoCtxEP0Data{}
	ctx := ctxproxy.NewCtxEP0Data(context.Background(), ctxInfo)

	buf := make([]byte, rawgadget.PAGE_SIZE)
	nr, err := src.ReadContext(ctx, buf)
	if err != nil {
		log.WithError(err).Errorf("[panic] Read err=%#v\n", err)
		ctxInfo.Close(err)
		return
	}
	buf = buf[:nr]

	nw, err := dst.WriteContext(ctx, buf)
	if err != nil {
		log.WithError(err).Errorf("[panic] Write err=%#v\n", err)
		ctxInfo.Close(err)
		return
	}

	if nr != nw {
		log.WithField("nr", nr).WithField("nw", nw).Errorf("customCopySingle WTF size not match")
	}

	ctxInfo.Close(nil)
}
