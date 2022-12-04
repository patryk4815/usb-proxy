package main

import (
	"errors"
	"flag"
	"github.com/google/gousb"
	"github.com/patryk4815/usb-proxy/pkg/hostproxy"
	"github.com/patryk4815/usb-proxy/pkg/rawproxy"
	log "github.com/sirupsen/logrus"
	"runtime"
	"strconv"
	"time"
)

func parseVendorIDOrProductID(v string) int {
	if len(v) != 4 {
		log.WithError(errors.New("len(v) != 4")).Fatalf("invalid vendor_id or product_id")
	}

	vv, err := strconv.ParseInt(v, 16, 64)
	if err != nil {
		log.WithError(err).Fatalf("invalid vendor_id or product_id")
	}
	return int(vv)
}

func main() {
	//debug.SetGCPercent(-1)

	driver := flag.String("driver", "", "eg. '-driver fe980000.usb'")
	device := flag.String("device", "", "eg. '-device fe980000.usb'")
	vendorIdStr := flag.String("vendor_id", "", "Bus 002 Device 012: ID 0abb:213b, first param eg. '-vendor_id 0abb'")
	productIdStr := flag.String("product_id", "", "Bus 002 Device 012: ID 0abb:213b, next param eg. '-product_id 213b'")
	flag.Parse()

	if driver == nil {
		log.Fatalf("required driver")
	}
	if device == nil {
		log.Fatalf("required device")
	}
	if vendorIdStr == nil {
		log.Fatalf("required vendor_id")
	}
	if productIdStr == nil {
		log.Fatalf("required product_id")
	}

	vendorId := parseVendorIDOrProductID(*vendorIdStr)
	productId := parseVendorIDOrProductID(*productIdStr)

	runtime.LockOSThread() // just in case

	log.SetLevel(log.DebugLevel) // TODO: control from cli
	log.Println("start")

	host := hostproxy.New()
	defer host.Close()
	host.Open(gousb.ID(vendorId), gousb.ID(productId))

	raw := rawproxy.New()
	defer raw.Close()
	raw.Open(*driver, *device)

	raw.SetHostProxy(host)
	host.SetRawProxy(raw)

	go raw.EP0Loop()

	tunnelEp0 := &rawproxy.WrapperUsbEP0{Self: raw}
	go customCopy(host, tunnelEp0)
	customCopy(tunnelEp0, host)

	log.Println("finished click ctrl+c")
	time.Sleep(time.Second * 30)
}
