package lzhuff

import (
	"fmt"

	"github.com/koron-go/lha/bitio"
)

type tree struct {
	l []uint16
	v []uint16

	left  []uint16
	right []uint16
}

func newTree(nl int, nv int) *tree {
	return &tree{
		l: make([]uint16, nl),
		v: make([]uint16, nv),
	}
}

// readAsP reads stream as T or P table.
func (t *tree) readAsP(r *bitio.Reader, bits int, special int) error {
	ubits := uint(bits)
	n0, err := r.ReadBits16(ubits)
	if err != nil {
		return err
	}
	if n0 == 0 {
		return t.setup0(r, ubits)
	}

	var (
		n  = int(n0)
		nl = len(t.l)
		i  = 0
	)
	if n > nl {
		n = nl
	}
	for i < n {
		c, err := r.ReadBits16(3)
		if err != nil {
			return err
		}
		if c == 7 {
			c2, err := r.CountTrues(13)
			if err != nil {
				return err
			}
			c += uint16(c2)
		}

		t.l[i] = c
		i++
		if i == special {
			c, err := r.ReadBits16(2)
			if err != nil {
				return err
			}
			for c > 0 && i < n {
				t.l[i] = 0
				i++
				c--
			}
		}
	}

	// fill remain with 0
	for i < nl {
		t.l[i] = 0
		i++
	}
	return t.setupTree(8)
}

// readAsC reads stream as C table.
func (t *tree) readAsC(r *bitio.Reader, bits int, tmp *tree) error {
	ubits := uint(bits)
	n0, err := r.ReadBits16(ubits)
	if err != nil {
		return err
	}
	if n0 == 0 {
		return t.setup0(r, ubits)
	}

	var (
		n  = int(n0)
		nl = len(t.l)
		i  = 0
	)
	if n > nl {
		n = nl
	}
	for i < n {
		c, err := tmp.decode(r, 8)
		if err != nil {
			return err
		}

		if c > 2 {
			t.l[i] = c - 2
			i++
			continue
		}

		switch c {
		case 0:
			c = 1
		case 1:
			c, err = r.ReadBits16(4)
			if err != nil {
				return err
			}
			c += 3
		default:
			c, err = r.ReadBits16(ubits)
			if err != nil {
				return err
			}
			c += 20
		}
		for i < nl && c > 0 {
			t.l[i] = 0
			i++
			c--
		}
	}

	// fill remain with 0
	for i < nl {
		t.l[i] = 0
		i++
	}
	return t.setupTree(12)
}

func (t *tree) setup0(r *bitio.Reader, bits uint) error {
	c, err := r.ReadBits16(bits)
	if err != nil {
		return err
	}
	for i := range t.l {
		t.l[i] = 0
	}
	for i := range t.v {
		t.v[i] = c
	}
	return nil
}

func (t *tree) getL(n int) uint16 {
	return t.l[n]
}

func (t *tree) getV(n int) uint16 {
	return t.v[n]
}

func (t *tree) setupTree(bits int) error {
	var (
		count  = make([]uint16, 17)
		weight = make([]uint16, 17)
		start  = make([]uint16, 17)
	)
	for i := range weight {
		weight[i] = 1 << uint(16-i)
	}

	// count
	for i := range t.l {
		if int(t.l[i]) >= len(count) {
			return fmt.Errorf("bad tree, length overflow: %d", t.l[i])
		}
		count[t.l[i]]++
	}

	// calculate first code
	total := uint(0)
	for i := 1; i < len(start); i++ {
		start[i] = uint16(total)
		total += uint(weight[i] * count[i])
	}
	if total != 0x10000 {
		return fmt.Errorf("bad tree, total unexpected: %04x", total)
	}

	// shift data to make table
	m := uint(16 - bits)
	for i := 1; i <= bits; i++ {
		start[i] >>= m
		weight[i] >>= m
	}

	j := start[bits+1] >> m
	k := 1 << uint(bits)
	if k > 4096 {
		k = 4096
	}
	if j != 0 {
		for i := int(j); i < k; i++ {
			t.v[i] = 0
		}
	}

	// create tree
	var (
		avail = uint16(len(t.l))
		vp    *uint16
		left  = make([]uint16, 4096)
		right = make([]uint16, 4096)
	)
	for i, v := range t.l {
		if v == 0 {
			continue
		}
		c := start[v] + weight[v]
		if int(v) <= bits {
			// code in array
			if c > 4096 {
				c = 4096
			}
			for j := start[v]; j < c; j++ {
				t.v[j] = uint16(i)
			}
		} else {
			// code not in array
			x := start[v]
			if x>>m > 4096 {
				return fmt.Errorf("bad tree, big start: %d %d %d", i, v, x)
			}
			vp = &t.v[x>>m]
			x <<= uint(bits)
			for n := int(v) - bits; n > 0; n-- {
				if *vp == 0 {
					left[avail] = 0
					right[avail] = 0
					*vp = avail
					avail++
				}
				if (x & 0x8000) != 0 {
					vp = &t.right[*vp]
				} else {
					vp = &t.left[*vp]
				}
				x <<= 1
			}
			*vp = uint16(i)
		}

		start[v] = c
	}

	t.left = left
	t.right = right
	return nil
}

func (t *tree) decode(r *bitio.Reader, bits uint) (uint16, error) {
	c0, err := r.PeekBits(bits)
	if err != nil {
		return 0, err
	}
	if int(c0) >= len(t.v) {
		return 0, fmt.Errorf("c0 overflow: %d >= %d", c0, len(t.v))
	}
	c := t.v[c0]

	// when hits t.l array.
	if int(c) < len(t.l) {
		err := r.SkipBits(uint(t.l[c]))
		if err != nil {
			return 0, err
		}
		return c, nil
	}

	err = r.SkipBits(bits)
	if err != nil {
		return 0, nil
	}

	for int(c) >= len(t.l) {
		b, err := r.ReadBit()
		if err != nil {
			return 0, err
		}
		if b {
			if int(c) >= len(t.right) {
				return 0, fmt.Errorf("over right: %d >= %d", c, len(t.right))
			}
			c = t.right[c]
		} else {
			if int(c) >= len(t.left) {
				return 0, fmt.Errorf("over left: %d >= %d", c, len(t.left))
			}
			c = t.left[c]
		}
		if int(c) >= len(t.left) {
			return 0, fmt.Errorf("over left: %d >= %d", c, len(t.left))
		}
		if c == t.left[c] {
			break
		}
	}

	return c, nil
}
