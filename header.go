package lha

import "time"

// Header is header of file in LHA archive.
type Header struct {
	Size         uint16
	Method       string
	PackedSize   uint32
	OriginalSize uint32
	Time         time.Time
	Attribute    uint8
	Level        uint8
	CRC          crc16
	OSID         uint8
	NextSize     uint16
}
