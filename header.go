package lha

import "time"

// Header is header of file in LHA archive.
type Header struct {
	Size         uint16
	Sum          uint8 // for level 0
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
	DOS          struct {
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

type ExtendType uint8

const (
	ExtendGeneric ExtendType = ExtendType(0)

	ExtendUNIX   = 'U'
	ExtendMSDOS  = 'M'
	ExtendMACOS  = 'm'
	ExtendOS9    = '9'
	ExtendOS2    = '2'
	ExtendOS68K  = 'K'
	ExtendOS386  = '3'
	ExtendHuman  = 'H'
	ExtendCPM    = 'C'
	ExtendFLEX   = 'F'
	ExtendRUNSER = 'R'
)
