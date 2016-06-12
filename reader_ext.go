package lha

import "os"

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
	*(*uint16)(&h.HeaderCRC), err = r.readUint16NoCRC()
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
	// TODO: parse Windows time.
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
