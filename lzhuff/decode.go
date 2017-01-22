package lzhuff

import (
	"io"

	"github.com/koron/go-lha/slide"
)

// Decoder provides huffman decoder interface.
type Decoder interface {
	DecodeC() (uint16, error)
	DecodeP() (uint16, error)
}

// Decode decodes huffman encoding.
func Decode(d Decoder, w io.Writer, bits, adjust uint, size int) (n int, crc uint16, err error) {
	sw := slide.NewWriter(w, bits)
	for sw.Len() < size {
		c, err := d.DecodeC()
		if err != nil {
			return sw.Len(), 0, err
		}
		if c < 256 {
			err := sw.WriteByte(byte(c))
			if err != nil {
				return sw.Len(), 0, err
			}
			continue
		}
		off, err := d.DecodeP()
		if err != nil {
			return sw.Len(), 0, nil
		}
		if _, err := sw.WriteCopy(int(off), int(uint(c)-adjust)); err != nil {
			return sw.Len(), 0, nil
		}
	}
	if err := sw.Flush(); err != nil {
		return sw.Len(), 0, err
	}
	return sw.Len(), sw.CRC16(), nil
}
