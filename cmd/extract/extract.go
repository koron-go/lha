package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/koron/go-lha"
)

func extract(r *lha.Reader, h *lha.Header) error {
	if h.Dir != "" {
		err := os.MkdirAll(h.Dir, 0777)
		if err != nil {
			return err
		}
	}
	name := filepath.Join(h.Dir, h.Name)
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	n, err := r.Decode(f)
	f.Close()
	if err != nil {
		os.Remove(name)
		return err
	}
	fmt.Printf("%s - %d bytes decoded\n", name, n)
	return nil
}

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
		h, err := r.NextHeader()
		if err != nil {
			log.Fatal(err)
		}
		if h == nil {
			break
		}
		err = extract(r, h)
		if err != nil {
			log.Fatal(err)
		}
	}
}
