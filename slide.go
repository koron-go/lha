package lha

import "io"

type sliceWriter struct {
	wr  io.Writer
	cnt int
	crc crc16
	buf []byte
	loc int
}

func newSliceWriter(w io.Writer, bits uint) *sliceWriter {
	buf := make([]byte, 1<<bits)
	for i := range buf {
		buf[i] = ' '
	}
	return &sliceWriter{
		wr:  w,
		buf: buf,
	}
}

func (sw *sliceWriter) Write(p []byte) (int, error) {
	for i, b := range p {
		if err := sw.WriteByte(b); err != nil {
			return i, err
		}
	}
	return len(p), nil
}

func (sw *sliceWriter) WriteByte(b byte) error {
	sw.buf[sw.loc] = b
	sw.loc++
	if sw.loc == len(sw.buf) {
		err := sw.Flush()
		if err != nil {
			return err
		}
	}
	sw.cnt++
	return nil
}

func (sw *sliceWriter) Flush() error {
	if sw.loc == 0 {
		return nil
	}
	_, err := sw.wr.Write(sw.buf[:sw.loc])
	if err != nil {
		return err
	}
	sw.crc = sw.crc.update(sw.buf[:sw.loc])
	sw.loc = 0
	return nil
}

func (sw *sliceWriter) writeCopy(off, size int) (int, error) {
	var (
		st = (sw.loc - off - 1 + len(sw.buf)) % len(sw.buf)
		r  = len(sw.buf) - st
		nw = 0
	)
	for size > 0 {
		if size <= r {
			n, err := sw.Write(sw.buf[st : st+size])
			nw += n
			return nw, err
		}
		n, err := sw.Write(sw.buf[st : st+r])
		nw += n
		if err != nil {
			return nw, err
		}
		size -= n
		st = 0
		r = len(sw.buf)
		if size < r {
			r = size
		}
	}
	return nw, nil
}

func slideDecode(hd huffDecoder, w io.Writer, bits, adjust uint, size int) (crc16, error) {
	sw := newSliceWriter(w, bits)
	hd.start()
	defer hd.end()
	for sw.cnt < size {
		c, err := hd.decodeC()
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
		off, err := hd.decodeP()
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
