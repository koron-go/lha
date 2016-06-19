package huff

import "io"

type staticDecoder struct {
	raw   io.Reader
	pbits int
	pnum  int
}

// NewStaticDecoder creates a new static huffman decoder.
func NewStaticDecoder(rd io.Reader, pbits, pnum int) Decoder {
	return &staticDecoder{
		raw:   rd,
		pbits: pbits,
		pnum:  pnum,
	}
}

func (sd *staticDecoder) DecodeC() (code uint16, err error) {
	// TODO:
	return 0, nil
}

func (sd *staticDecoder) DecodeP() (offset uint16, err error) {
	// TODO:
	return 0, nil
}
