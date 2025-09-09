package lha

import (
	"errors"
	"os"
	"time"
)

// Header is header of file in LHA archive.
type Header struct {
	Size         uint16
	Sum          uint8 // for level 0, 1
	Method       string
	PackedSize   uint64
	OriginalSize uint64
	Time         time.Time
	Attribute    uint8
	Level        uint8
	CRC          uint16
	OSID         uint8

	Name string

	HeaderCRC    *uint16
	ExtendType   ExtendType
	MinorVersion uint8
	Dir          string

	ExtendedHeaderSize uint64

	DOS  HeaderDOS
	UNIX HeaderUNIX
}

// HeaderDOS is exntended header for DOS.
type HeaderDOS struct {
	Attr uint16
	Time uint64
}

// HeaderUNIX is exntended header for UNIX.
type HeaderUNIX struct {
	Perm  uint16
	GID   uint16
	UID   uint16
	Group string
	User  string
	Time  time.Time
}

// ExtendType is type of exntend part.
type ExtendType uint8

const (
	// ExtendGeneric is generic extend.
	ExtendGeneric ExtendType = 0
	// ExtendUNIX is extend for UNIX.
	ExtendUNIX ExtendType = 'U'
	// ExtendMSDOS is extend for MS-DOS.
	ExtendMSDOS ExtendType = 'M'
	// ExtendMACOS is extend for MacOS
	ExtendMACOS ExtendType = 'm'
	// ExtendOS9 is extend for OS/9
	ExtendOS9 ExtendType = '9'
	// ExtendOS2 is extend for OS/2
	ExtendOS2 ExtendType = '2'
	// ExtendOS68K is extend for OS68K
	ExtendOS68K ExtendType = 'K'
	// ExtendOS386 is extend for OS386
	ExtendOS386 ExtendType = '3'
	// ExtendHuman is extend for Human
	ExtendHuman ExtendType = 'H'
	// ExtendCPM is extend for CP/M
	ExtendCPM ExtendType = 'C'
	// ExtendFLEX is extend for FLEX
	ExtendFLEX ExtendType = 'F'
	// ExtendRUNSER is extend for RUNSER
	ExtendRUNSER ExtendType = 'R'
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
	// FIXME: verify check sum
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
	h := new(Header)
	headerSize, _ := r.readUint8()
	h.Size = uint16(headerSize)
	h.Sum, _ = r.readUint8()
	// FIXME: verify check sum
	h.Method, _ = r.readStringN(5)
	packedSize, _ := r.readUint32()
	h.PackedSize = uint64(packedSize)
	originalSize, _ := r.readUint32()
	h.OriginalSize = uint64(originalSize)
	dt, _ := r.readUint32()
	h.Time = fromDOSTimestamp(dt)
	h.Attribute, _ = r.readUint8() // 0x20 fixed
	h.Level, _ = r.readUint8()
	nameLen, _ := r.readUint8()
	h.Name, _ = r.readStringN(int(nameLen))
	*(*uint16)(&h.CRC), _ = r.readUint16()
	h.OSID, _ = r.readUint8()
	// FIXME: consider 64bit length.
	if remain := int(h.Size) - int(nameLen) - 25; remain > 0 {
		r.skip(remain)
	}
	nextSize, _ := r.readUint16()
	readAllExtendedHeaders(r, h, nextSize)
	if r.err != nil {
		return nil, r.err
	}
	return h, nil
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
	h := new(Header)
	_, _ = r.readUint16()
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
	headerSize, _ := r.readUint32()
	nextSize, _ := r.readUint32()
	readAllExtendedHeaders(r, h, uint16(nextSize))
	// FIXME: consider 64bit length.
	if remain := int(headerSize) - int(r.cnt); remain >= 0 {
		r.skip(remain)
	}
	if r.err != nil {
		return nil, r.err
	}
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
	h.ExtendedHeaderSize += uint64(size)
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
	for i, b := range d {
		if b == 0xff {
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
