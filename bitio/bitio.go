package bitio

import "errors"

var (
	// ErrTooMuchBits indicates required number of bits is too large.
	ErrTooMuchBits = errors.New("too much bits")

	// ErrTooLessBits means too less bits are available than required.
	ErrTooLessBits = errors.New("too less bits")
)

type bits struct {
	v uint64
	n uint
}

func (b *bits) skip(nbits uint) error {
	if nbits > b.n {
		return ErrTooMuchBits
	}
	b.v <<= nbits
	b.n -= nbits
	return nil
}

func (b *bits) peek(nbits uint) (uint64, error) {
	if nbits > b.n {
		return 0, ErrTooMuchBits
	}
	d := b.v >> (64 - nbits)
	return d, nil
}

func (b *bits) read(nbits uint) (uint64, error) {
	if nbits > b.n {
		return 0, ErrTooMuchBits
	}
	d := b.v >> (64 - nbits)
	b.v <<= nbits
	b.n -= nbits
	return d, nil
}

func (b *bits) write(d uint64, nbits uint) error {
	if nbits+b.n > 64 {
		return ErrTooMuchBits
	}
	b.v += d << (64 - b.n - nbits)
	b.n += nbits
	return nil
}

func (b *bits) get(n uint) (bool, error) {
	if n >= b.n {
		return false, ErrTooMuchBits
	}
	v := b.v & (1 << (63 - n)) != 0
	return v, nil
}

func (b *bits) set(p []byte) {
	if len(p) > 8 {
		return
	}
	b.v = 0
	b.n = uint(len(p)) * 8
	for i, v := range p {
		b.v += uint64(v) << uint(56-i*8)
	}
}
