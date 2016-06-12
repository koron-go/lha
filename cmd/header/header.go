package main

import (
	"flag"
	"log"
	"os"
	"github.com/koron/go-lha"
)

func main() {
	flag.Parse()
	name := flag.Arg(0)
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := lha.NewReader(f)
	h, err := r.ReadHeader()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", h)
	log.Printf("Header CRC: %04x", r.CRC16())
}
