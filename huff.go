package lha

import "io"

type huffDecoder struct {
	start   func(r *Reader)
	decodeC func(r *Reader) (uint16, error)
	decodeP func(r *Reader) (uint16, error)
}

func (hd huffDecoder) decode(r *Reader, w io.Writer, bits, adjust uint, size int) (crc16, error) {
	sw := newSliceWriter(w, bits)
	hd.start(r)
	for sw.cnt < size {
		c, err := hd.decodeC(r)
		if err != nil {
			return 0, err
		}
		if c < 256 {
			err := sw.WriteByte(byte(c))
			if err != nil {
				return 0, err
			}
			continue
		}
		off, err := hd.decodeP(r)
		if err != nil {
			return 0, nil
		}
		if _, err := sw.writeCopy(int(off), int(uint(c)-adjust)); err != nil {
			return 0, nil
		}
	}
	if err := sw.Flush(); err != nil {
		return 0, nil
	}
	return sw.crc, nil
}

func decodeST1Start(r *Reader) {
	// TODO:
}

func decodeST1C(r *Reader) (uint16, error) {
	// TODO:
	return 0, nil
}

func decodeST1P(r *Reader) (uint16, error) {
	// TODO:
	return 0, nil
}
