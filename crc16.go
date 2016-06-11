package lha

type crc16 uint16

func (c crc16) updateByte(b byte) crc16 {
	return crcTable[(c^crc16(b))&0xff] ^ (c >> 8)
}

func (c crc16) update(d []byte) crc16 {
	for _, b := range d {
		c = c.updateByte(b)
	}
	return c
}

const crcPoly = crc16(0xA001)

var crcTable = make([]crc16, 256)

func init() {
	for i := range crcTable {
		r := crc16(i)
		for j := 0; j < 8; j++ {
			if r&1 != 0 {
				r = (r >> 1) ^ crcPoly
			} else {
				r >>= 1
			}
		}
		crcTable[i] = r
	}
}
