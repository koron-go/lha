package lha

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

// Reader is LHA archive reader.
type Reader struct {
	raw io.Reader
	rd  *bufio.Reader
}

// NewReader creates LHA archive reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		raw: r,
		rd:  bufio.NewReader(r),
	}
}

func (r *Reader) readHeader() (*Header, error) {
	const commonHeaderSize = 21
	b, err := r.rd.ReadByte()
	if err != nil && err != io.EOF {
		return nil, err
	} else if err == io.EOF || b == 0 {
		return nil, nil
	}
	d := make([]byte, commonHeaderSize)
	d[0] = b
	// FIXME: check read length
	if _, err = r.rd.Read(d[1:]); err != nil {
		return nil, err
	}
	switch lv := d[20]; lv {
	case 0:
		return r.readHeaderLv0(d)
	case 1:
		return r.readHeaderLv1(d)
	case 2:
		return r.readHeaderLv2(d)
	case 3:
		return r.readHeaderLv3(d)
	default:
		return nil, fmt.Errorf("unknown level header: %d", lv)
	}
}

func (r *Reader) readHeaderLv0(d []byte) (*Header, error) {
	log.Println("readHeader:", 0)
	// TODO:
	return nil, nil
}

func (r *Reader) readHeaderLv1(d []byte) (*Header, error) {
	log.Println("readHeader:", 1)
	// TODO:
	return nil, nil
}

func (r *Reader) readHeaderLv2(d []byte) (*Header, error) {
	log.Println("readHeader:", 2)
	// TODO:
	return nil, nil
}

func (r *Reader) readHeaderLv3(d []byte) (*Header, error) {
	log.Println("readHeader:", 3)
	// TODO:
	return nil, nil
}
