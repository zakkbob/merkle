package merkle

import (
	"encoding/hex"
	"strings"
)

type Tree struct {
	root node
	// map containing leafs and their indexes
	// json encode/decode
}

func NewTree(leaves []string, hashFn func([]byte) []byte) Tree {
	nodes := make([]*node, len(leaves))
	for i, h := range leaves {
		nodes[i] = &node{
			left:  nil,
			right: nil,
			hash:  h,
		}
	}

	for len(nodes) != 1 {
		l := len(nodes)
		for i := range l / 2 {
			left, right := nodes[2*i], nodes[2*i+1]
			nodes[i] = newNode(left, right, hashFn)
		}
		if l%2 == 1 {
			nodes[l/2] = nodes[l-1]
		}
		nodes = nodes[:(l/2 + l%2)]
	}

	return Tree{
		root: *nodes[0],
	}
}

func (t *Tree) String() string {
	b := &strings.Builder{}
	t.root.buildString(0, []bool{}, false, b)
	return b.String()
}

func (t *Tree) Prove(hash string) []string {
	return t.root.prove(hash)
}

type node struct {
	left  *node
	right *node
	hash  string
}

func newNode(left *node, right *node, hashFn func([]byte) []byte) *node {
	return &node{
		left:  left,
		right: right,
		hash:  string(hashFn(append([]byte(left.hash), []byte(right.hash)...))),
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

func (n *node) prove(hash string) []string {
	if n.isLeaf() {
		if n.hash == hash {
			return []string{hash}
		}
		return nil
	}

	if p := n.left.prove(hash); p != nil {
		return append(p, n.right.hash)
	}

	if p := n.right.prove(hash); p != nil {
		return append(p, n.left.hash)
	}

	return nil
}
