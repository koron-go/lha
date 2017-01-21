package huff

// Tree represents a huffman tree.
type Tree struct {
	Leaves []int32
	Nodes  []int
}

// New creates a empty huffman tree with size.
func New(n int) *Tree {
	return &Tree{
		Leaves: make([]int32, n),
		Nodes:  make([]int, (n-1)*2),
	}
}
