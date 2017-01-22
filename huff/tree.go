package huff

// Tree represents a huffman tree.
type Tree struct {
	// nodes represents nodes of huffman tree. A node composed from two int.
	// First represents left(0) link, second is for right(1) link.  Link is, 0
	// means not assign, positive value means link to other node, negative
	// value means a leaf.
	nodes []int

	// next leaf to be allocated.
	leaf int

	// last allocated node.
	node int
}

// New creates a empty huffman tree with capacity.
func New(capacity int) *Tree {
	t := &Tree{
		nodes: make([]int, (capacity-1)*2),
	}
	return t
}

// Reset resets huffman tree.
func (t *Tree) Reset() *Tree {
	for i := range t.nodes {
		t.nodes[i] = 0
	}
	t.leaf = 0
	t.node = 0
	return t
}

// Add adds a leaf with bit length.
func (t *Tree) Add(length int) (leaf int, err error) {
	//log.Printf("Add(%d)", length)
	if length == 0 {
		t.leaf++
		return -1, nil
	}
	return t.add(length-1, 0)
}

func (t *Tree) add(length, idx int) (leaf int, err error) {
	//log.Printf("add(%d, %d)", length, idx)
	if length == 0 {
		return t.newLeaf(idx)
	}
	leaf, err = t.newNodeLeaf(length, idx)
	if err == nil {
		return leaf, nil
	}
	return t.newNodeLeaf(length, idx+1)
}

func (t *Tree) newNodeLeaf(length, idx int) (leaf int, err error) {
	//log.Printf("newNodeLeaf(%d, %d)", length, idx)
	v := t.nodes[idx]
	switch {
	case v < 0:
		return 0, ErrNoMoreLeaves
	case v == 0:
		nextNode := t.node + 2
		if nextNode >= len(t.nodes) {
			return 0, ErrNoMoreNodes
		}
		t.nodes[idx] = nextNode
		t.node = nextNode
	}
	return t.add(length-1, t.nodes[idx])
}

func (t *Tree) newLeaf(idx int) (leaf int, err error) {
	//log.Printf("newLeaf(%d)", idx)
	if t.nodes[idx] != 0 {
		if t.nodes[idx+1] != 0 {
			return -1, ErrNoMoreLeaves
		}
		if t.nodes[idx] > 0 {
			t.nodes[idx+1] = t.nodes[idx]
		} else {
			idx++
		}
	}
	leaf = t.leaf
	t.leaf++
	t.nodes[idx] = -leaf - 1
	return leaf, nil
}
