package lha

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/koron/go-lha/crc16"
)

const (
	commonHeaderSize = 21
)

var (
	errTooShortExtendedHeader = errors.New("too short extended header")
	errHeaderCRCMismatch      = errors.New("header CRC mismatch")
	errBodyCRCMismatch        = errors.New("body CRC mismatch")
)

// Reader is LHA archive reader.
type Reader struct {
	raw io.Reader
	br  *bufio.Reader
	err error
	cnt uint64
	crc crc16.Hash16

	curr *Header
}

// NewReader creates LHA archive reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		raw: r,
		br:  bufio.NewReader(r),
		crc: crc16.NewIBM(),
	}
}

// CRC16 returns current CRC16 value.
func (r *Reader) CRC16() uint16 {
	return r.crc.Sum16()
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
	r.crc.Reset()
	h, err = proc(r)
	if err != nil {
		return nil, err
	}
	if h.HeaderCRC != nil && *h.HeaderCRC != r.crc.Sum16() {
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
	r.crc.Write(d)
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
	r.crc.Write([]byte{b0})
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
	r.crc.Write([]byte{b0, b1})
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
	r.crc.Write([]byte{0, 0})
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
	r.crc.Write([]byte{b0, b1, b2, b3})
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
	r.crc.Write(d)
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
	m, err := getMethod(r.curr.Method)
	if err != nil {
		return 0, err
	}
	lr := &io.LimitedReader{
		R: r.br,
		N: int64(r.curr.PackedSize),
	}
	defer func() {
		// count read length.
		r.cnt += r.curr.PackedSize - uint64(lr.N)
	}()
	n, crc, err := m.decode(lr, w, int(r.curr.OriginalSize))
	if err != nil {
		return 0, err
	}
	if crc != r.curr.CRC {
		return 0, errBodyCRCMismatch
	}
	return n, nil
}
