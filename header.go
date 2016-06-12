package lha

import "time"

// Header is header of file in LHA archive.
type Header struct {
	Size         uint16
	Method       string
	PackedSize   uint64
	OriginalSize uint64
	Time         time.Time
	Attribute    uint8
	Level        uint8
	CRC          crc16
	OSID         uint8

	HeaderCRC *crc16
	Name      string
	Dir       string
	DOS       struct {
		Attr uint16
		Time uint64
	}
	UNIX struct {
		Perm  uint16
		GID   uint16
		UID   uint16
		Group string
		User  string
		Time  time.Time
	}
}
