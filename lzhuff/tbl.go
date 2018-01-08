package lzhuff

import "github.com/koron/go-lha/bitio"

type tbl struct {
	l []uint16
	v []uint16
}

func newTbl(nl int, nv int) *tbl {
	return &tbl{
		l: make([]uint16, nl),
		v: make([]uint16, nv),
	}
}

func (t *tbl) read(r *bitio.Reader, bits int, special int) error {
	ubits := uint(bits)
	n0, err := r.ReadBits16(ubits)
	if err != nil {
		return err
	}
	if n0 == 0 {
		return t.setup0(r, ubits)
	}

	var (
		n = int(n0)
		nl = len(t.l)
	)
	if n > nl {
		n = nl
	}
	var i = 0
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
	return t.setupV()
}

func (t *tbl) setup0(r *bitio.Reader, bits uint) error {
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

func (t *tbl) setupV() error {
	// TODO: setup t.v (table conttents)
	return nil
}

func (t *tbl) getL(n int) uint16 {
	return t.l[n]
}

func (t *tbl) getV(n int) uint16 {
	return t.v[n]
}
