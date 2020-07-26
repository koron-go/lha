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
