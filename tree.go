package merkle

type Node struct {
	Left  *Node
	Right *Node
	Hash  string
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
