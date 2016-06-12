package lha

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"time"
)

var (
	errTooShortExtendedHeader = errors.New("too short extended header")
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

func (r *Reader) readHeader() (*Header, error) {
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
	switch lv := d[20]; lv {
	case 0:
		return r.readHeaderLv0()
	case 1:
		return r.readHeaderLv1()
	case 2:
		return r.readHeaderLv2()
	case 3:
		return r.readHeaderLv3()
	default:
		return nil, fmt.Errorf("unknown level header: %d", lv)
	}
}

func (r *Reader) readHeaderLv0() (*Header, error) {
	log.Println("readHeader:", 0)
	// TODO:
	return nil, nil
}

func (r *Reader) readHeaderLv1() (*Header, error) {
	log.Println("readHeader:", 1)
	// TODO:
	return nil, nil
}

func (r *Reader) readHeaderLv2() (*Header, error) {
	if r.err != nil {
		return nil, r.err
	}
	h := new(Header)
	h.Size, _ = r.readUint16()
	h.Method, _ = r.readStringN(5)
	packedSize, _ := r.readUint32()
	h.PackedSize = uint64(packedSize)
	originalSize, _ := r.readUint32()
	h.OriginalSize = uint64(originalSize)
	h.Time, _ = r.readTime()
	h.Attribute, _ = r.readUint8()
	h.Level, _ = r.readUint8()
	*(*uint16)(&h.CRC), _ = r.readUint16()
	h.OSID, _ = r.readUint8()
	nextSize, _ := r.readUint16()
	readAllExtendedHeaders(r, h, nextSize)
	if remain := int(h.Size) - r.cnt; remain > 0 {
		r.skip(remain)
	}
	if r.err != nil {
		return nil, r.err
	}
	return h, nil
}

func (r *Reader) readHeaderLv3() (*Header, error) {
	log.Println("readHeader:", 3)
	// TODO:
	return nil, nil
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
	r.cnt += len(d)
	r.crc = r.crc.updateBytes(d)
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
	r.cnt += 1
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
