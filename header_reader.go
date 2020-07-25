package lha

import (
	"errors"
	"log"
	"os"
	"time"
)

type headerReader func(r *Reader) (*Header, error)

var headerReaders = map[byte]headerReader{
	0: readHeaderLv0,
	1: readHeaderLv1,
	2: readHeaderLv2,
	3: readHeaderLv3,
}

func readHeaderLv0(r *Reader) (*Header, error) {
	h := new(Header)
	headerSize, _ := r.readUint8()
	h.Size = uint16(headerSize)
	h.Sum, _ = r.readUint8()
	h.Method, _ = r.readStringN(5)
	packedSize, _ := r.readUint32()
	h.PackedSize = uint64(packedSize)
	originalSize, _ := r.readUint32()
	h.OriginalSize = uint64(originalSize)
	dt, _ := r.readUint32()
	h.Time = fromDOSTimestamp(dt)
	h.Attribute, _ = r.readUint8()
	h.Level, _ = r.readUint8()
	nameLen, _ := r.readUint8()
	h.Name, _ = r.readStringN(int(nameLen))

	extendSize := int(headerSize) + 2 - int(nameLen) - 24
	if extendSize < 0 {
		if extendSize == -2 {
			h.HeaderCRC = nil
			return h, nil
		}
		return nil, errors.New("unknown header")
	}

	*(*uint16)(&h.CRC), _ = r.readUint16()
	if extendSize == 0 {
		return h, nil
	}

	extendType, _ := r.readUint8()
	h.ExtendType = ExtendType(extendType)
	extendSize--

	if ExtendType(extendType) == ExtendUNIX {
		if extendSize >= 11 {
			h.MinorVersion, _ = r.readUint8()
			h.UNIX.Time, _ = r.readTime()
			h.UNIX.Perm, _ = r.readUint16()
			h.UNIX.UID, _ = r.readUint16()
			h.UNIX.GID, _ = r.readUint16()
		} else {
			h.ExtendType = ExtendGeneric
		}
	}

	if remain := int(h.Size) - int(r.cnt); remain > 0 {
		r.skip(remain)
	}
	if r.err != nil {
		return nil, r.err
	}
	return h, nil
}

func readHeaderLv1(r *Reader) (*Header, error) {
	log.Println("readHeader:", 1)
	// TODO: support header LV1
	return nil, nil
}

func readHeaderLv2(r *Reader) (*Header, error) {
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
	// FIXME: consider 64bit length.
	if remain := int(h.Size) - int(r.cnt); remain > 0 {
		r.skip(remain)
	}
	if r.err != nil {
		return nil, r.err
	}
	return h, nil
}

func readHeaderLv3(r *Reader) (*Header, error) {
	log.Println("readHeader:", 3)
	// TODO: support header LV3
	return nil, nil
}

type exHeaderReader func(r *Reader, h *Header, size int) (remain int, err error)

var exHeaderReaders = map[uint8]exHeaderReader{
	0x00: readHeaderCRC,
	0x01: readFilename,
	0x02: readDirectory,
	0x40: readDOSAttr,
	0x41: readWinTime,
	0x42: readWinSize,
	0x50: readUNIXPerm,
	0x51: readUNIXGIDUID,
	0x52: readUNIXGroup,
	0x53: readUNIXUser,
	0x54: readUNIXTime,
}

func readAllExtendedHeaders(r *Reader, h *Header, size uint16) error {
	if r.err != nil {
		return r.err
	}
	for size > 0 {
		if size < 3 {
			r.err = errTooShortExtendedHeader
			return r.err
		}
		size, r.err = readExtendedHeader(r, h, size)
		if r.err != nil {
			return r.err
		}
	}
	return nil
}

func readExtendedHeader(r *Reader, h *Header, size uint16) (uint16, error) {
	t, err := r.readUint8()
	if err != nil {
		return 0, err
	}
	proc, ok := exHeaderReaders[t]
	remain := int(size) - 3
	if ok {
		remain, err = proc(r, h, remain)
		if err != nil {
			return 0, err
		}
	}
	if remain > 0 {
		r.skip(remain)
	}
	return r.readUint16()
}

func readHeaderCRC(r *Reader, h *Header, size int) (remain int, err error) {
	crc := new(uint16)
	*crc, err = r.readUint16NoCRC()
	if err == nil {
		h.HeaderCRC = crc
	}
	return remain - 2, err
}

func readFilename(r *Reader, h *Header, size int) (remain int, err error) {
	h.Name, err = r.readStringN(size)
	return 0, err
}

func readDirectory(r *Reader, h *Header, size int) (remain int, err error) {
	d, err := r.readBytes(size)
	for i := range d {
		if d[i] == 0xff {
			d[i] = os.PathSeparator
		}
	}
	h.Dir = string(d)
	return 0, err
}

func readDOSAttr(r *Reader, h *Header, size int) (remain int, err error) {
	h.DOS.Attr, err = r.readUint16()
	return size - 2, err
}

func readWinTime(r *Reader, h *Header, size int) (remain int, err error) {
	// FIXME: parse Windows time.
	return size, r.err
}

func readWinSize(r *Reader, h *Header, size int) (remain int, err error) {
	h.PackedSize, _ = r.readUint64()
	h.OriginalSize, err = r.readUint64()
	return size - 16, err
}

func readUNIXPerm(r *Reader, h *Header, size int) (remain int, err error) {
	h.UNIX.Perm, err = r.readUint16()
	return size - 2, err
}

func readUNIXGIDUID(r *Reader, h *Header, size int) (remain int, err error) {
	h.UNIX.GID, _ = r.readUint16()
	h.UNIX.UID, err = r.readUint16()
	return size - 4, err
}

func readUNIXGroup(r *Reader, h *Header, size int) (remain int, err error) {
	h.UNIX.Group, err = r.readStringN(size)
	return 0, err
}

func readUNIXUser(r *Reader, h *Header, size int) (remain int, err error) {
	h.UNIX.User, err = r.readStringN(size)
	return 0, err
}

func readUNIXTime(r *Reader, h *Header, size int) (remain int, err error) {
	h.UNIX.Time, err = r.readTime()
	return size - 4, err
}

func fromDOSTimestamp(v uint32) time.Time {
	y := int(1980 + (v>>25)&0x7f)
	m := int((v >> 21) & 0x0f)
	d := int((v >> 16) & 0x1f)
	h := int((v >> 11) & 0x1f)
	mi := int((v >> 5) & 0x3f)
	s := int(v&0x1f) * 2
	return time.Date(y, time.Month(m), d, h, mi, s, 0, time.Local)
}
