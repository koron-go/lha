package crc16

import "io"

// Hash16 is the common interface implemented by all 16-bit hash functions.
type Hash16 interface {
	io.Writer

	Reset()

	Sum16() uint16
}

type hash16 struct {
	tab *Table
	crc uint16
}

func (h *hash16) Write(p []byte) (int, error) {
	h.crc = Update(h.crc, h.tab, p)
	return len(p), nil
}

func (h *hash16) Reset() {
	h.crc = 0
}

func (h *hash16) Sum16() uint16 {
	return h.crc
}
