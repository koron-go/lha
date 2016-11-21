package huff

import (
	"io"

	"github.com/koron/go-lha/bitio"
)

const nc = 510

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

func (sd *staticDecoder) prepareBlock() error {
	nblock, err := sd.brd.ReadBits16(16)
	if err != nil {
		return err
	}
	// TODO:
	sd.nblock = int(nblock)
	return nil
}

func (sd *staticDecoder) DecodeC() (code uint16, err error) {
	if sd.nblock == 0 {
		if err := sd.prepareBlock(); err != nil {
			return 0, err
		}
	}
	sd.nblock--
	code, err = sd.brd.ReadBits16(12)
	if err != nil {
		return 0, err
	}
	if code < nc {
		// TODO: push back 12-r.lenC[code] bits.
		return code, nil
	}
	// TODO: docode C
	return 0, nil
}

func (sd *staticDecoder) DecodeP() (offset uint16, err error) {
	// TODO: decode P
	return 0, nil
}
