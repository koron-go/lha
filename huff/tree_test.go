package huff

import (
	"testing"

	"github.com/koron-go/lha/internal/assert"
)

func TestTreeAdd(t *testing.T) {
	f := func(data []int, expected *Tree) {
		actual := New(len(data))
		for i, d := range data {
			leaf, err := actual.Add(d)
			if err != nil {
				t.Errorf("%+v: Add(%d[#%d]) failed: %s", data, d, i, err)
				return
			}
			if d == 0 {
				if leaf != -1 {
					t.Errorf("%+v: Add(%d) should return -1: %d", data, d, leaf)
				}
				continue
			}
			if leaf != i {
				t.Errorf("%+v: Add(%d) unexpected leaf:%d expected:%d",
					data, d, leaf, i)
			}
		}
		assert.Equalf(t, actual, expected, "%+v: unmatched tree", data)
	}
	f([]int{1, 1}, &Tree{nodes: []int{-1, -2}, leaf: 2, node: 0})
	f([]int{1, 2, 2}, &Tree{nodes: []int{-1, 2, -2, -3}, leaf: 3, node: 2})
	f([]int{2, 2, 1}, &Tree{nodes: []int{-3, 2, -1, -2}, leaf: 3, node: 2})
	f([]int{2, 1, 2}, &Tree{nodes: []int{-2, 2, -1, -3}, leaf: 3, node: 2})
	f([]int{2, 2, 2, 2}, &Tree{
		nodes: []int{2, 4, -1, -2, -3, -4},
		leaf:  4,
		node:  4,
	})
	f([]int{1, 2, 3, 4, 5, 6, 7, 7}, &Tree{
		nodes: []int{-1, 2, -2, 4, -3, 6, -4, 8, -5, 10, -6, 12, -7, -8},
		leaf:  8,
		node:  12,
	})
	f([]int{7, 7, 6, 5, 4, 3, 2, 1}, &Tree{
		nodes: []int{-8, 2, -7, 4, -6, 6, -5, 8, -4, 10, -3, 12, -1, -2},
		leaf:  8,
		node:  12,
	})
	f([]int{0, 1, 1}, &Tree{nodes: []int{-2, -3, 0, 0}, leaf: 3, node: 0})
	f([]int{1, 0, 1}, &Tree{nodes: []int{-1, -3, 0, 0}, leaf: 3, node: 0})
}
