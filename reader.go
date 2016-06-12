package lha

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"
)

const (
	commonHeaderSize = 21
)

var (
	errTooShortExtendedHeader = errors.New("too short extended header")
	errHeaderCRCMismatch      = errors.New("header CRC mismatch")
)

// Reader is LHA archive reader.
type Reader struct {
	raw io.Reader
	br  *bufio.Reader
	err error
	cnt uint64
	crc crc16

	curr *Header
}

// NewReader creates LHA archive reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		raw: r,
		br:  bufio.NewReader(r),
	}
}

// CRC16 returns current CRC16 value.
func (r *Reader) CRC16() uint16 {
	return uint16(r.crc)
}

// NextHeader reads a next file header.
func (r *Reader) NextHeader() (h *Header, err error) {
	if r.err != nil {
		return nil, r.err
	}
	if err := r.seekNext(); err != nil {
		return nil, err
	}
	lv, err := r.peekHeaderLevel()
	if err == io.EOF {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	proc, ok := headerReaders[lv]
	if !ok {
		return nil, fmt.Errorf("unknown header level: %d", lv)
	}
	r.cnt = 0
	r.crc = 0
	h, err = proc(r)
	if err != nil {
		return nil, err
	}
	if h.HeaderCRC != nil && *h.HeaderCRC != r.crc {
		return nil, errHeaderCRCMismatch
	}
	r.cnt = 0
	r.curr = new(Header)
	*r.curr = *h
	return h, nil
}

func (r *Reader) seekNext() error {
	if r.curr != nil && r.curr.PackedSize > r.cnt {
		remain := r.curr.PackedSize - r.cnt
		r.curr = nil
		// FIXME: consider 64bit length.
		_, r.err = r.br.Discard(int(remain))
		if r.err != nil {
			return r.err
		}
	}
	return nil
}

func (r *Reader) peekHeaderLevel() (lv byte, err error) {
	if r.err != nil {
		return 0, r.err
	}
	var d []byte
	d, r.err = r.br.Peek(1)
	if r.err != nil {
		return 0, r.err
	}
	if d[0] == 0 {
		r.err = io.EOF
		return 0, io.EOF
	}
	d, r.err = r.br.Peek(commonHeaderSize)
	if r.err != nil {
		return 0, r.err
	}
	return d[20], nil
}

func (r *Reader) skip(n int) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	var d []byte
	d, r.err = r.readBytes(n)
	if r.err != nil {
		return 0, r.err
	}
	return len(d), nil
}

func (r *Reader) readBytes(n int) ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	d := make([]byte, n)
	_, r.err = io.ReadFull(r.br, d)
	if r.err != nil {
		return nil, r.err
	}
	r.cnt += uint64(len(d))
	r.crc = r.crc.updateBytes(d)
	return d, nil
}

func (r *Reader) readStringN(n int) (string, error) {
	d, err := r.readBytes(n)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func (r *Reader) readUint8() (uint8, error) {
	if r.err != nil {
		return 0, r.err
	}
	b0, err := r.br.ReadByte()
	if err != nil {
		r.err = err
		return 0, r.err
	}
	r.cnt++
	r.crc = r.crc.updateByte(b0)
	return uint8(b0), nil
}

func (r *Reader) readUint16() (uint16, error) {
	if r.err != nil {
		return 0, r.err
	}
	b0, _ := r.br.ReadByte()
	b1, err := r.br.ReadByte()
	if err != nil {
		r.err = err
		return 0, r.err
	}
	r.cnt += 2
	r.crc = r.crc.update(b0, b1)
	return uint16(b1)<<8 + uint16(b0), nil
}

func (r *Reader) readUint16NoCRC() (uint16, error) {
	if r.err != nil {
		return 0, r.err
	}
	b0, _ := r.br.ReadByte()
	b1, err := r.br.ReadByte()
	if err != nil {
		r.err = err
		return 0, r.err
	}
	r.cnt += 2
	r.crc = r.crc.update(0, 0)
	return uint16(b1)<<8 + uint16(b0), nil
}

func (r *Reader) readUint32() (uint32, error) {
	if r.err != nil {
		return 0, r.err
	}
	var (
		b0, b1, b2, b3 byte
	)
	b0, _ = r.br.ReadByte()
	b1, _ = r.br.ReadByte()
	b2, _ = r.br.ReadByte()
	b3, r.err = r.br.ReadByte()
	if r.err != nil {
		return 0, r.err
	}
	r.cnt += 4
	r.crc = r.crc.update(b0, b1, b2, b3)
	return uint32(b3)<<24 + uint32(b2)<<16 + uint32(b1)<<8 + uint32(b0), nil
}

func (r *Reader) readUint64() (uint64, error) {
	if r.err != nil {
		return 0, r.err
	}
	d, err := r.readBytes(8)
	if err != nil {
		return 0, err
	}
	r.cnt += uint64(len(d))
	r.crc = r.crc.updateBytes(d)
	return binary.LittleEndian.Uint64(d), nil
}

func (r *Reader) readTime() (time.Time, error) {
	if r.err != nil {
		return time.Time{}, r.err
	}
	var n uint32
	n, r.err = r.readUint32()
	if r.err != nil {
		return time.Time{}, r.err
	}
	return time.Unix(int64(n), 0), nil
}

// Decode decode a file to w.  It returns decoded size and error.
func (r *Reader) Decode(w io.Writer) (decoded int, err error) {
	// TODO:
	return 0, nil
}
