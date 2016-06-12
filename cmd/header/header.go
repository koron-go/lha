package main

import (
	"flag"
	"fmt"
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
	for {
		h, err := r.ReadHeader()
		if err != nil {
			log.Fatal(err)
		}
		if h == nil {
			break
		}
		fmt.Printf("%+v\n", h)
		if _, err := r.Discard(int(h.PackedSize)); err != nil {
			log.Fatal(err)
		}
	}
}
