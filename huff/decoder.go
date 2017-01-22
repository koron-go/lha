package huff

// BitReader implements bit stream.
type BitReader interface {
	ReadBit() (bool, error)
}

// Decoder decodes a value from bitstream
type Decoder interface {
	Decode(r BitReader) (int32, error)
}

// Decode decodes a value from bitstream.
func (t *Tree) Decode(r BitReader) (int, error) {
	if t.leaf == 0 {
		return 0, nil
	}
	idx := 0
	for {
		b, err := r.ReadBit()
		if err != nil {
			return 0, err
		}
		if b {
			idx++
		}
		idx = t.nodes[idx]
		switch {
		case idx == 0:
			return 0, ErrIncompleteTree
		case idx < 0:
			return 1 - idx, nil
		}
	}
}
