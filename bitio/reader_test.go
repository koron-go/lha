package bitio

import (
	"bytes"
	"io"
	"testing"
)

type readBits struct {
	nbits uint
	rval  uint64
	rerr  error
}

func TestReaderReadBits(t *testing.T) {
	f := func(t *testing.T, d []byte, p []readBits) *Reader {
		r := NewReader(bytes.NewReader(d))
		for _, q := range p {
			val, err := r.ReadBits(q.nbits)
			assertEquals(t, val, q.rval, "Reader.ReadBits() returned value")
			assertEquals(t, err, q.rerr, "Reader.ReadBits() returned error")
		}
		return r
	}
	f(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, []readBits{
		{64, 0, nil},
		{64, 0, io.EOF},
	})
	f(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, []readBits{
		{63, 0, nil},
		{1, 0, nil},
		{1, 0, io.EOF},
	})
	f(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, []readBits{
		{63, 0, nil},
		{2, 0, ErrTooLessBits},
	})
}
