package lzhuff

import (
	"io"

	"github.com/koron/go-lha/bitio"
)

const (
	nt    = 19
	nc    = 510
	tbits = 5
	cbits = 9
)

type staticDecoder struct {
	raw   io.Reader
	brd   *bitio.Reader
	pbits int
	pnum  int

	nblock int

	c *tree
	p *tree
}

// NewStaticDecoder creates a new static huffman decoder.
func NewStaticDecoder(rd io.Reader, pbits, pnum int) Decoder {
	return &staticDecoder{
		raw:   rd,
		brd:   bitio.NewReader(rd),
		pbits: pbits,
		pnum:  pnum,
	}
}

func (sd *staticDecoder) prepareBlock() error {
	nblock, err := sd.brd.ReadBits16(16)
	if err != nil {
		return err
	}
	sd.nblock = int(nblock)
	if err := sd.prepareC(); err != nil {
		return err
	}
	if err := sd.prepareP(); err != nil {
		return err
	}
	return nil
}

func (sd *staticDecoder) prepareC() error {
	tt := newTree(nt, 256)
	err := tt.readAsP(sd.brd, 5, 3)
	if err != nil {
		return err
	}
	tc := newTree(nc, 4096)
	err = tc.readAsC(sd.brd, cbits, tt)
	if err != nil {
		return err
	}
	sd.c = tc
	return nil
}

func (sd *staticDecoder) prepareP() error {
	tp := newTree(sd.pnum, 256)
	err := tp.readAsP(sd.brd, sd.pbits, -1)
	if err != nil {
		return err
	}
	sd.p = tp
	return nil
}

func (sd *staticDecoder) DecodeC() (code uint16, err error) {
	if sd.nblock == 0 {
		if err := sd.prepareBlock(); err != nil {
			return 0, err
		}
	}
	sd.nblock--
	v, err := sd.c.decode(sd.brd, 12)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (sd *staticDecoder) DecodeP() (offset uint16, err error) {
	v, err := sd.p.decode(sd.brd, 8)
	if err != nil {
		return 0, err
	}
	if v > 0 {
		w := v - 1
		d, err := sd.brd.ReadBits16(uint(w))
		if err != nil {
			return 0, nil
		}
		v = (1 << w) + d
	}
	return v, nil
}
