package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/koron-go/lha"
)

func extractTest(r *lha.Reader, h *lha.Header) error {
	name := filepath.Join(h.Dir, h.Name)
	n, err := r.Decode(ioutil.Discard)
	if err != nil {
		return err
	}
	fmt.Printf("  %s - %d bytes decoded\n", name, n)
	return nil
}

func extractLha(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Printf("%s - extact as lha\n", name)
	r := lha.NewReader(f)
	for {
		h, err := r.NextHeader()
		if err != nil {
			return err
		}
		if h == nil {
			break
		}
		err = extractTest(r, h)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	for _, arg := range flag.Args() {
		err := extractLha(arg)
		if err != nil {
			log.Fatal(err)
		}
	}
}
