package lzhuff

import (
	"fmt"
	"io"

	"github.com/koron/go-lha/bitio"
)

const (
	nc    = 510
	tbits = 5
)

type staticDecoder struct {
	raw   io.Reader
	brd   *bitio.Reader
	pbits int
	pnum  int

	nblock int

	lenC   [nc]byte
	tableC [4096]byte
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

func (sd *staticDecoder) countTrues(max uint16) (count uint16, err error) {
	for count <= max {
		b, err := sd.brd.PeekBit()
		if err != nil {
			return 0, err
		}
		if !b {
			return count, nil
		}
		if err := sd.brd.SkipBit(); err != nil {
			return 0, err
		}
		count++
	}
	return 0, fmt.Errorf("too match trues, over %d", max)
}

func (sd *staticDecoder) readLengths(nbits uint) ([]uint16, error) {
	n, err := sd.brd.ReadBits16(nbits)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		c, err := sd.brd.ReadBits16(nbits)
		if err != nil {
			return nil, err
		}
		return []uint16{c}, nil
	}
	lengths := make([]uint16, n)
	for i := range lengths {
		l, err := sd.brd.ReadBits16(3)
		if err != nil {
			return nil, err
		}
		if l == 0x07 {
			c, err := sd.countTrues(12)
			if err != nil {
				return nil, err
			}
			l += c
		}
		lengths[i] = l
		// TODO:
	}
	return lengths, err
}

func (sd *staticDecoder) prepareBlock() error {
	nblock, err := sd.brd.ReadBits16(16)
	if err != nil {
		return err
	}
	sd.nblock = int(nblock)
	// TODO:
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
	code = uint16(sd.tableC[n])
	if code < nc {
		err := sd.brd.SkipBits(uint(sd.lenC[code]))
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
	// TODO: decode P
	return 0, nil
}
