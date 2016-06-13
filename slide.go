package lha

import "io"

type slideWriter struct {
	wr  io.Writer
	cnt int
	crc crc16
	buf []byte
	loc int
}

func newSlideWriter(w io.Writer, bits uint) *slideWriter {
	buf := make([]byte, 1<<bits)
	for i := range buf {
		buf[i] = ' '
	}
	return &slideWriter{
		wr:  w,
		buf: buf,
	}
}

func (sw *slideWriter) Write(p []byte) (int, error) {
	for i, b := range p {
		if err := sw.WriteByte(b); err != nil {
			return i, err
		}
	}
	return len(p), nil
}

func (sw *slideWriter) WriteByte(b byte) error {
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

func (sw *slideWriter) Flush() error {
	if sw.loc == 0 {
		return nil
	}
	_, err := sw.wr.Write(sw.buf[:sw.loc])
	if err != nil {
		return err
	}
	sw.crc = sw.crc.updateBytes(sw.buf[:sw.loc])
	sw.loc = 0
	return nil
}

func (sw *slideWriter) WriteCopy(off, size int) (int, error) {
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
