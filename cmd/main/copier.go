package main

import (
	"io"
	"log"
)

func customCopy(dst io.Writer, src io.Reader) {
	for {
		buf := make([]byte, 1024)
		n, err := src.Read(buf)
		if err != nil {
			log.Printf("[panic] Read err=%#v\n", err)
			panic(err)
		}
		buf = buf[:n]

		_, err = dst.Write(buf)
		if err != nil {
			log.Printf("[panic] Write err=%#v\n", err)
			panic(err)
		}
	}
}
