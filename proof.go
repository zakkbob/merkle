package merkle

import "encoding/json"

type Proof struct {
	path      []string
	treeDepth int
	leafIndex int
}

func (p *Proof) MarshalJSON() ([]byte, error) {
	v := struct {
		Path      []string `json:"path"`
		TreeDepth int      `json:"tree_depth"`
		LeafIndex int      `json:"leaf_index"`
	}{
		Path:      p.path,
		TreeDepth: p.treeDepth,
		LeafIndex: p.leafIndex,
	}

	return json.Marshal(v)
}

func (p *Proof) UnmarshalJSON(data []byte) error {
	v := struct {
		Path      []string `json:"path"`
		TreeDepth int      `json:"tree_depth"`
		LeafIndex int      `json:"leaf_index"`
	}{}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	p.path = v.Path
	p.treeDepth = v.TreeDepth
	p.leafIndex = v.LeafIndex
	return nil
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
