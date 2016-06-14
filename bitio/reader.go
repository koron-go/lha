package bitio

import "io"

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

func (r *Reader) readBits(n uint) (uint64, error) {
	if n > r.curr.n {
		err := r.loadBits(n - r.curr.n)
		if err != nil {
			return 0, err
		}
	}
	return r.curr.read(n)
}

// loadBits loads at least n bits.
func (r *Reader) loadBits(nbits uint) error {
	if nbits -= r.loadAhead(nbits); nbits == 0 {
		return nil
	}
	nbytes := int((nbits + 7) / 8)
	if nbytes > 8 {
		return ErrTooMuchBits
	}
	for {
		n, err := r.rd.Read(r.buf[:nbytes])
		if n > 0 {
			r.ahead.set(r.buf[:n])
			nbytes -= n
		}
		if nbits -= r.loadAhead(nbits); nbits == 0 {
			return nil
		}
		if err == io.EOF {
			return ErrTooLessBits
		} else if err != nil {
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
