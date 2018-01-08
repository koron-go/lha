package lzhuff

import (
	"io"

	"github.com/koron/go-lha/bitio"
)

const (
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

	tC *tbl
	tP *tbl
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
	tt := newTbl(19, 256)
	err := tt.read(sd.brd, 5, 3)
	if err != nil {
		return err
	}
	tc := newTbl(nc, 4096)
	err = tt.read(sd.brd, cbits, -1)
	// TODO: customize tt.read() for C table.
	if err != nil {
		return err
	}
	sd.tC = tc
	return nil
}

func (sd *staticDecoder) prepareP() error {
	tp := newTbl(sd.pnum, 256)
	err := tp.read(sd.brd, sd.pbits, -1)
	if err != nil {
		return err
	}
	sd.tP = tp
	return nil
}

func (sd *staticDecoder) DecodeC() (code uint16, err error) {
	if sd.nblock == 0 {
		if err := sd.prepareBlock(); err != nil {
			return 0, err
		}
	}
	sd.nblock--
	n, err := sd.brd.ReadBits16(12)
	if err != nil {
		return 0, err
	}
	code = sd.tC.getV(int(n))
	if code < nc {
		err := sd.brd.SkipBits(uint(sd.tC.getL(int(code))))
		if err != nil {
			return 0, err
		}
		return code, nil
	}
	err = sd.brd.SkipBits(12)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (sd *staticDecoder) DecodeP() (offset uint16, err error) {
	n, err := sd.brd.PeekBits(8)
	if err != nil {
		return 0, nil
	}
	v := sd.tP.getV(int(n))
	if int(v) < sd.pnum {
		err := sd.brd.SkipBits(uint(sd.tP.getL(int(v))))
		if err != nil {
			return 0, nil
		}
	} else {
		err := sd.brd.SkipBits(8)
		if err != nil {
			return 0, nil
		}
		// TODO:
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
