package bitio

import (
	"bytes"
	"io"
	"testing"

	"github.com/koron-go/lha/internal/assert"
)

type readBits struct {
	nbits uint
	rval  uint64
	rerr  error
}

func TestReaderReadBits(t *testing.T) {
	f := func(t *testing.T, d []byte, p []readBits) *Reader {
		r := NewReader(bytes.NewReader(d))
		for i, q := range p {
			val, err := r.ReadBits(q.nbits)
			assert.Equalf(t, val, q.rval, "Reader.ReadBits() returned value for #%d", i)
			assert.Equalf(t, err, q.rerr, "Reader.ReadBits() returned error for #%d", i)
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
	// 1101 0010 0010 0000
	f(t, []byte{0xD2, 0x20}, []readBits{
		{1, 1, nil},
		{2, 2, nil},
		{3, 4, nil},
		{4, 8, nil},
		{5, 16, nil},
		{2, 0, ErrTooLessBits},
		{1, 0, nil},
		{1, 0, io.EOF},
	})
}
