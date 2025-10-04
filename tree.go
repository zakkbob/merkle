package merkle

import (
	"encoding/hex"
	"strings"
)

type Tree struct {
	// map containing leafs and their indexes
	// root node
	// json encode/decode
	root Node
}

func NewTree(leaves []string, hashFn func([]byte) []byte) Tree {
	nodes := make([]*Node, len(leaves))
	for i, h := range leaves {
		nodes[i] = &Node{
			Left:  nil,
			Right: nil,
			Hash:  h,
		}
	}

	for len(nodes) != 1 {
		l := len(nodes)
		for i := range l / 2 {
			left, right := nodes[2*i], nodes[2*i+1]
			nodes[i] = NewNode(left, right, hashFn)
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
	t.root.BuildString(0, []bool{}, false, b)
	return b.String()
}

func (t *Tree) Prove(hash string) []string {
	return t.root.Prove(hash)
}

type Node struct {
	Left  *Node
	Right *Node
	Hash  string
}

func NewNode(left *Node, right *Node, hashFn func([]byte) []byte) *Node {
	return &Node{
		Left:  left,
		Right: right,
		Hash:  string(hashFn(append([]byte(left.Hash), []byte(right.Hash)...))),
	}
}

func (n *Node) BuildString(depth int, lines []bool, leftNode bool, builder *strings.Builder) {
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

	builder.WriteString(hex.EncodeToString([]byte(n.Hash)))
	builder.WriteByte('\n')

	if n.IsLeaf() {
		return
	}

	lines = append(lines, leftNode)

	n.Left.BuildString(depth+1, lines, true, builder)
	n.Right.BuildString(depth+1, lines, false, builder)
}

func (n *Node) IsLeaf() bool {
	return n.Left == nil || n.Right == nil
}

func (n *Node) Prove(hash string) []string {
	if n.IsLeaf() {
		if n.Hash == hash {
			return []string{hash}
		}
		return nil
	}

	if p := n.Left.Prove(hash); p != nil {
		return append(p, n.Right.Hash)
	}

	if p := n.Right.Prove(hash); p != nil {
		return append(p, n.Left.Hash)
	}

	return nil
}
