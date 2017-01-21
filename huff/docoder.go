package huff

import "github.com/koron/go-lha/bitio"

// Decoder decodes a value from bitstream
type Decoder interface {
	Decode(r *bitio.Reader) (int32, error)
}

// Decode decodes a value from bitstream.
func (t *Tree) Decode(r *bitio.Reader) (int32, error) {
	if len(t.Leaves) == 1 {
		return t.Leaves[0], nil
	}
	idx := 0
	for idx >= 0 {
		b, err := r.ReadBit()
		if err != nil {
			return 0, err
		}
		if b {
			idx++
		}
		idx = t.Nodes[idx]
	}
	return t.Leaves[1-idx], nil
}
