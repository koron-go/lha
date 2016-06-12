package lha

import (
	"log"
	"os"
	"testing"
)

func TestDQ(t *testing.T) {
	f, err := os.Open("tmp/dq2passwd.lzh")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r := NewReader(f)
	h, err := r.ReadHeader()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%+v", h)
	log.Printf("Header CRC: %04x", r.CRC16())
}
