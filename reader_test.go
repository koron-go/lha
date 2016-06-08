package lha

import (
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
	if _, err := r.readHeader(); err != nil{
		t.Fatal(err)
	}
}
