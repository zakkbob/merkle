package merkle

import (
	"encoding/hex"
	"errors"
	"strings"
)

var (
	ErrLeafNotFound = errors.New("leaf not found in merkle tree")
)

type Proof struct {
	path      []string
	treeDepth int
	leafIndex int
}

func (p *Proof) Verify(root string, leaf string, hashFn func([]byte) []byte) bool {
	if p.path[0] != leaf {
		return false
	}

	hash := []byte(p.path[0])
	direction := p.leafIndex

	for i := 1; i < len(p.path); i++ {
		leftNode := (direction & 1) == 0
		direction >>= 1

		if leftNode {
			b := append(hash, p.path[i]...)
			hash = hashFn(b)
		} else {
			b := append([]byte(p.path[i]), hash...)
			hash = hashFn(b)
		}

	}

	return string(hash) == root
}

func (p *Proof) Path() []string {
	path := make([]string, len(p.path))
	copy(path, p.path)
	return path
}

type Tree struct {
	depth  int
	node   node
	leaves map[string]int
	// json encode/decode
}

func NewTree(leaves []string, hashFn func([]byte) []byte) Tree {
	leafMap := make(map[string]int, len(leaves))
	nodes := make([]*node, len(leaves))
	for i, h := range leaves {
		leafMap[h] = i
		nodes[i] = &node{
			left:  nil,
			right: nil,
			hash:  h,
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

	return Tree{
		node:   *nodes[0],
		leaves: leafMap,
		depth:  depth,
	}
}

func (t *Tree) Root() []byte {
	return []byte(t.node.hash)
}

func (t *Tree) String() string {
	b := &strings.Builder{}
	t.node.buildString(0, []bool{}, false, b)
	return b.String()
}

func (t *Tree) Prove(leaf string) (Proof, error) {
	i, ok := t.leaves[leaf]
	if !ok {
		return Proof{}, ErrLeafNotFound
	}
	return Proof{
		path:      t.node.prove(leaf),
		leafIndex: i,
		treeDepth: t.depth,
	}, nil
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
