package main

import (
	"github.com/patryk4815/usb-proxy/pkg/hostproxy"
	"github.com/patryk4815/usb-proxy/pkg/rawproxy"
	"log"
	"time"
)

func main() {
	log.Println("start")

	host := hostproxy.New()
	defer host.Close()
	host.Open()

	raw := rawproxy.New()
	defer raw.Close()
	raw.Open()

	raw.SetHostProxy(host)
	host.SetRawProxy(raw)

	go raw.EP0Loop()

	tunnelEp0 := &rawproxy.WrapperUsbRawControlIO{Self: raw}
	go customCopy(tunnelEp0, host) // MAX COPY = PAGE_SIZE kernel
	customCopy(host, tunnelEp0)

	log.Println("finished click ctrl+c")
	time.Sleep(time.Second * 600)
}
