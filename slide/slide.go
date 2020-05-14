package slide

import (
	"io"

	"github.com/koron-go/lha/crc16"
)

// Writer provides slide window (buffered) writer.
type Writer struct {
	wr  io.Writer
	cnt int
	crc crc16.Hash16
	buf []byte
	loc int
}

// NewWriter create a slide window writer.
func NewWriter(w io.Writer, bits uint) *Writer {
	buf := make([]byte, 1<<bits)
	for i := range buf {
		buf[i] = ' '
	}
	return &Writer{
		wr:  w,
		buf: buf,
		crc: crc16.NewIBM(),
	}
}

// Flush flush all buffered data.
func (w *Writer) Flush() error {
	if w.loc == 0 {
		return nil
	}
	_, err := w.wr.Write(w.buf[:w.loc])
	if err != nil {
		return err
	}
	w.crc.Write(w.buf[:w.loc])
	w.loc = 0
	return nil
}

// WriteByte writes a byte.
func (w *Writer) WriteByte(b byte) error {
	w.buf[w.loc] = b
	w.loc++
	if w.loc == len(w.buf) {
		err := w.Flush()
		if err != nil {
			return err
		}
	}
	w.cnt++
	return nil
}

// Write writes data.
func (w *Writer) Write(p []byte) (int, error) {
	for i, b := range p {
		if err := w.WriteByte(b); err != nil {
			return i, err
		}
	}
	return len(p), nil
}

// WriteCopy writes datat which copied from window buffer.
func (w *Writer) WriteCopy(off, size int) (int, error) {
	var (
		st = (w.loc - off - 1 + len(w.buf)) % len(w.buf)
		r  = len(w.buf) - st
		nw = 0
	)
	for size > 0 {
		if size <= r {
			n, err := w.Write(w.buf[st : st+size])
			nw += n
			return nw, err
		}
		n, err := w.Write(w.buf[st : st+r])
		nw += n
		if err != nil {
			return nw, err
		}
		size -= n
		st = 0
		r = len(w.buf)
		if size < r {
			r = size
		}
	}
	return nw, nil
}

// CRC16 returns CRC-16 (IBM) for written bytes.
func (w *Writer) CRC16() uint16 {
	return w.crc.Sum16()
}

// Len returns written length of bytes.
func (w *Writer) Len() int {
	return w.cnt
}
