package huff

// Tree represents a huffman tree.
type Tree struct {
	// Leaves store
	Leaves []int32

	// Nodes store links to other nodes (0 or greater) or leaves (minus).
	Nodes  []int
}

// New creates a empty huffman tree with size.
func New(n int) *Tree {
	leaves := make([]int32, n)
	for i := range leaves {
		leaves[i] = int32(i)
	}

	nodes := make([]int, (n-1)*2)
	for i := 0; i < n-1; i++ {
		nodes[i] = -i - 1
		nodes[i+1] = (i + 1) * 2
	}
	nodes[len(nodes)-1] = -n - 1

	return &Tree{
		Leaves: leaves,
		Nodes:  nodes,
	}
}
