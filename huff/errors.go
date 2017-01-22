package huff

import "errors"

var (
	// ErrIncompleteTree represents incomplete huffman tree.
	ErrIncompleteTree = errors.New("incomplete huffman tree")

	// ErrNoMoreLeaves represents failed to allocate to a leaf.
	ErrNoMoreLeaves = errors.New("no more leaves")

	// ErrNoMoreNodes represents failed to allocate to a node.
	ErrNoMoreNodes = errors.New("no more nodes")
)
