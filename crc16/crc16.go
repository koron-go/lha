package crc16

// Table is a 256-word table representing the polynomial for efficient
// processing.
type Table [256]uint16

const (
	// IBM is by far and away the most common CRC-16 polynomial.
	IBM = 0xA001
)

var (
	// IBMTable is the table for the IBM polynomial.
	IBMTable = MakeTable(IBM)
)

// MakeTable returns a Table constructed from the specified polynomial. The
// contents of this Table must not be modified.
func MakeTable(poly uint16) *Table {
	t := new(Table)
	for i := range t {
		r := uint16(i)
		for j := 0; j < 8; j++ {
			if r&1 != 0 {
				r = (r >> 1) ^ poly
			} else {
				r >>= 1
			}
		}
		t[i] = r
	}
	return t
}

// Update returns the result of adding the bytes in p to the crc.
func Update(crc uint16, tab *Table, p []byte) uint16 {
	for _, b := range p {
		crc = tab[(crc^uint16(b))&0xff] ^ (crc >> 8)
	}
	return crc
}

// New creates a new crc16.Hash16 computing the CRC-16 checksum using the
// polynomial represented by the Table.  out in big-endian byte order.
func New(tab *Table) Hash16 {
	return &hash16{tab: tab}
}

// NewIBM creates a new crc16.Hash16 computing the CRC-16 checksum using the
// IBM polynomial.
func NewIBM() Hash16 {
	return New(IBMTable)
}
