package merkle

import (
	"encoding/json"
	"slices"
)

type Proof struct {
	Path      [][]byte
	TreeDepth int
	LeafIndex int
}

func (p *Proof) MarshalJSON() ([]byte, error) {
	v := struct {
		Path      [][]byte `json:"path"`
		TreeDepth int      `json:"tree_depth"`
		LeafIndex int      `json:"leaf_index"`
	}{
		Path:      p.Path,
		TreeDepth: p.TreeDepth,
		LeafIndex: p.LeafIndex,
	}

	return json.Marshal(v)
}

func (p *Proof) UnmarshalJSON(data []byte) error {
	v := struct {
		Path      [][]byte `json:"path"`
		TreeDepth int      `json:"tree_depth"`
		LeafIndex int      `json:"leaf_index"`
	}{}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	p.Path = v.Path
	p.TreeDepth = v.TreeDepth
	p.LeafIndex = v.LeafIndex
	return nil
}

func (p *Proof) Verify(root []byte, leaf []byte, hashFn func([]byte) []byte) bool {
	if !slices.Equal(p.Path[0], leaf) {
		return false
	}

	hash := p.Path[0]
	direction := p.LeafIndex

	for i := 1; i < len(p.Path); i++ {
		leftNode := (direction & 1) == 0
		direction >>= 1

		if leftNode {
			b := append(hash, p.Path[i]...)
			hash = hashFn(b)
		} else {
			b := append([]byte(p.Path[i]), hash...)
			hash = hashFn(b)
		}

	}

	return slices.Equal(hash, root)
}
