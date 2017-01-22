package huff

// Tree represents a huffman tree.
type Tree struct {
	// Leaves store
	Leaves []int32

	// Nodes store links to other nodes (0 or greater) or leaves (minus).
	Nodes []int
}

// New creates a empty huffman tree with size.
func New(n int) *Tree {
	t := &Tree{
		Leaves: make([]int32, n),
		Nodes:  make([]int, (n-1)*2),
	}
	return t.Reset()
}

// Reset resets huffman tree.
func (t *Tree) Reset() *Tree {
	n := len(t.Leaves)
	for i := range t.Leaves {
		t.Leaves[i] = int32(i)
	}
	for i := 0; i < n-1; i++ {
		t.Nodes[i] = -i - 1
		t.Nodes[i+1] = (i + 1) * 2
	}
	t.Nodes[len(t.Nodes)-1] = -n - 1
	return t
}
