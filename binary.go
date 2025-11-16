package merkle

import (
	"encoding/hex"
	"slices"
	"strings"
)

type BinaryTree struct {
	depth  int
	node   node
	leaves map[string]int
	// json encode/decode
}

// Constructs a binary merkle tree from the provided data, using the provided hash function.
// If the number of elements provided is not a power of 2, the tree is padded by duplicating the last node.
func NewBinaryTree(bs [][]byte, hashFn func([]byte) []byte) BinaryTree {
	leafMap := make(map[string]int, len(bs))
	nodes := make([]*node, len(bs))
	for i, b := range bs {
		leafMap[string(b)] = i
		nodes[i] = &node{
			left:  nil,
			right: nil,
			hash:  hashFn(b),
		}
	}

	depth := 0
	for len(nodes) != 1 {
		l := len(nodes)
		for i := range l / 2 {
			left, right := nodes[2*i], nodes[2*i+1]
			nodes[i] = newNode(left, right, hashFn)
		}
		if l%2 == 1 {
			nodes[l/2] = newNode(nodes[l-1], nodes[l-1], hashFn)
		}
		nodes = nodes[:(l/2 + l%2)]
		depth++
	}

	return BinaryTree{
		node:   *nodes[0],
		leaves: leafMap,
		depth:  depth,
	}
}

func (t *BinaryTree) Root() []byte {
	return []byte(t.node.hash)
}

func (t *BinaryTree) String() string {
	b := &strings.Builder{}
	t.node.buildString(0, []bool{}, false, b)
	return b.String()
}

func (t *BinaryTree) Prove(b []byte) (Proof, error) {
	i, ok := t.leaves[string(b)]
	if !ok {
		return Proof{}, ErrLeafNotFound
	}
	return Proof{
		Path:      t.node.prove(b),
		LeafIndex: i,
		TreeDepth: t.depth,
	}, nil
}

type node struct {
	left  *node
	right *node
	hash  []byte
}

func newNode(left *node, right *node, hashFn func([]byte) []byte) *node {
	return &node{
		left:  left,
		right: right,
		hash:  hashFn(append(left.hash, right.hash...)),
	}
}

func (n *node) buildString(depth int, lines []bool, leftNode bool, builder *strings.Builder) {
	for i := range depth {
		if lines[i] {
			builder.WriteString("│  ")

		} else {
			builder.WriteString("   ")

		}
	}

	if leftNode {
		builder.WriteString("├─ ")
	} else {
		builder.WriteString("└─ ")
	}

	builder.WriteString(hex.EncodeToString([]byte(n.hash)))
	builder.WriteByte('\n')

	if n.isLeaf() {
		return
	}

	lines = append(lines, leftNode)

	n.left.buildString(depth+1, lines, true, builder)
	n.right.buildString(depth+1, lines, false, builder)
}

func (n *node) isLeaf() bool {
	return n.left == nil || n.right == nil
}

func (n *node) prove(b []byte) [][]byte {
	if n.isLeaf() {
		if slices.Equal(n.hash, b) {
			return [][]byte{b}
		}
		return nil
	}

	if p := n.left.prove(b); p != nil {
		return append(p, n.right.hash)
	}

	if p := n.right.prove(b); p != nil {
		return append(p, n.left.hash)
	}

	return nil
}
