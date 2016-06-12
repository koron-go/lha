package lha

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"
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
	cnt int
	crc crc16
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

// ReadHeader reads a hedader entry.
func (r *Reader) ReadHeader() (h *Header, err error) {
	const commonHeaderSize = 21
	b, err := r.br.Peek(1)
	if err != nil && err != io.EOF {
		return nil, err
	} else if err == io.EOF || b[0] == 0 {
		return nil, nil
	}
	d, err := r.br.Peek(commonHeaderSize)
	if err != nil {
		return nil, err
	}
	proc, ok := headerReaders[d[20]]
	if !ok {
		return nil, fmt.Errorf("unknown level header: %d", d[20])
	}
	r.crc = 0
	h, err = proc(r)
	if err != nil {
		return nil, err
	}
	if h.HeaderCRC != nil && *h.HeaderCRC != r.crc {
		return nil, errHeaderCRCMismatch
	}
	return h, nil
}

// Discard skips the next n bytes.
func (r *Reader) Discard(n int) (discarded int, err error) {
	return r.br.Discard(n)
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
	r.cnt += len(d)
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
	r.cnt += len(d)
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
