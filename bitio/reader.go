package bitio

import (
	"io"
)

// Reader is a reader of bit stream.
type Reader struct {
	rd    io.Reader
	buf   [8]byte
	curr  bits
	ahead bits
}

// NewReader creates a bit stream reader.
func NewReader(rd io.Reader) *Reader {
	return &Reader{
		rd: rd,
	}
}

// NBits returns number of available bits without accessing underlying
// io.Reader.
func (r *Reader) NBits() uint {
	return r.curr.n + r.ahead.n
}

// ReadBits reads bits up to 64.
func (r *Reader) ReadBits(n uint) (uint64, error) {
	if n > 64 {
		return 0, ErrTooMuchBits
	}
	return r.readBits(n)
}

// ReadBits32 reads bits up to 32.
func (r *Reader) ReadBits32(n uint) (uint32, error) {
	if n > 32 {
		return 0, ErrTooMuchBits
	}
	d, err := r.readBits(n)
	return uint32(d), err
}

// ReadBits16 reads bits up to 16.
func (r *Reader) ReadBits16(n uint) (uint16, error) {
	if n > 16 {
		return 0, ErrTooMuchBits
	}
	d, err := r.readBits(n)
	return uint16(d), err
}

// ReadBits8 reads bits up to 8.
func (r *Reader) ReadBits8(n uint) (uint8, error) {
	if n > 8 {
		return 0, ErrTooMuchBits
	}
	d, err := r.readBits(n)
	return uint8(d), err
}

// ReadBit reads a bit.
func (r *Reader) ReadBit() (bool, error) {
	d, err := r.readBits(1)
	if err != nil {
		return false, err
	}
	return d != 0, nil
}

// PeekBit peeks a bit.
func (r *Reader) PeekBit() (bool, error) {
	d, err := r.PeekBits(1)
	if err != nil {
		return false, err
	}
	return d != 0, nil
}

// SkipBit skips a bit.
func (r *Reader) SkipBit() error {
	return r.SkipBits(1)
}

func (r *Reader) readBits(n uint) (uint64, error) {
	if n > r.curr.n {
		err := r.loadBits(n - r.curr.n)
		if err != nil {
			return 0, err
		}
	}
	return r.curr.read(n)
}

// PeekBits peeks some bits.
func (r *Reader) PeekBits(n uint) (uint64, error) {
	if n > r.curr.n {
		err := r.loadBits(n - r.curr.n)
		if err != nil {
			return 0, err
		}
	}
	return r.curr.peek(n)
}

// CountTrues counts continuous true bits.
func (r *Reader) CountTrues(n uint) (uint, error) {
	for i := uint(0); i < n; i++ {
		b, err := r.ReadBit()
		if err != nil {
			return 0, err
		}
		if !b {
			return i, nil
		}
	}
	return n, nil
}

// SkipBits skips some bits.
func (r *Reader) SkipBits(n uint) error {
	if n > r.curr.n {
		err := r.loadBits(n - r.curr.n)
		if err != nil {
			return err
		}
	}
	return r.curr.skip(n)
}

// loadBits loads at least n bits.
func (r *Reader) loadBits(nbits uint) error {
	if nbits > 64 {
		return ErrTooMuchBits
	}
	for {
		if nbits -= r.loadAhead(nbits); nbits == 0 {
			return nil
		}
		nbytes := int((nbits + 7) / 8)
		n, err := r.rd.Read(r.buf[:nbytes])
		if n > 0 {
			r.ahead.set(r.buf[:n])
			continue
		}
		if err != nil {
			if err == io.EOF && r.curr.n+r.ahead.n != 0 {
				//// supply zero bits if there are some bits remained.
				//err := r.ahead.write(0, nbits)
				//if err != nil {
				//	return err
				//}
				//continue
				return ErrTooLessBits
			}
			return err
		}
	}
}

func (r *Reader) loadAhead(nbits uint) uint {
	if r.ahead.n == 0 {
		return 0
	}
	n := nbits
	if n > r.ahead.n {
		n = r.ahead.n
	}
	d, _ := r.ahead.read(n)
	r.curr.write(d, n)
	return n
}
