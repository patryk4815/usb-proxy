package main

import (
	"log"
	"time"
)

func main() {
	log.Println("start")

	host := NewHost()
	host.Open()
	defer host.Close()

	raw := NewRawGadget()
	host.lolRaw = raw
	raw.lolHost = host
	go raw._WorkerDo()
	defer raw.Close()

	raw2 := &WrapperUsbRawControlIO{xRaw: raw}
	go customCopy(raw2, host)
	customCopy(host, raw2)

	time.Sleep(time.Second * 600)
}
