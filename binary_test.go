package merkle_test

import (
	"crypto/sha256"
	"encoding/json"
	"slices"
	"testing"

	"github.com/zakkbob/merkle"
)

func TestNewBinaryTree(t *testing.T) {
	leaves := [][]byte{}
	for i := range 100 {
		leaves = append(leaves, []byte{byte(i)})
	}
	tree := merkle.NewBinaryTree(leaves, func(data []byte) []byte {
		b := sha256.Sum256(data)
		return b[:]
	})

	t.Log("\n" + tree.String())
}

func TestBinaryTreeProof(t *testing.T) {
	// abcc
	// ├─ ab
	// │  ├─ a
	// │  └─ b
	// └─ cc
	//    ├─ c
	//    └─ c
	hashFn := func(b []byte) []byte {
		return b
	}

	tree := merkle.NewBinaryTree([][]byte{[]byte("a"), []byte("b"), []byte("c")}, hashFn)
	t.Log(tree.String())

	tests := []struct {
		Target   []byte
		Expected [][]byte
	}{
		{
			Target:   []byte("a"),
			Expected: [][]byte{[]byte("a"), []byte("b"), []byte("cc")},
		},
		{
			Target:   []byte("b"),
			Expected: [][]byte{[]byte("b"), []byte("a"), []byte("cc")},
		},
		{
			Target:   []byte("c"),
			Expected: [][]byte{[]byte("c"), []byte("c"), []byte("ab")},
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.Target), func(t *testing.T) {
			p, err := tree.Prove(tt.Target)
			if err != nil {
				t.Fatal(err)
			}

			b, err := p.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}

			var dp merkle.Proof
			err = json.Unmarshal(b, &dp)
			if err != nil {
				t.Fatal(err)
			}

			if len(dp.Path) != len(tt.Expected) {
				t.Fatalf("got %v; expected %v", dp.Path, tt.Expected)
			}

			for i := range len(dp.Path) {
				if !slices.Equal(dp.Path[i], tt.Expected[i]) {
					t.Fatalf("got %v; expected %v", dp.Path, tt.Expected)
				}
			}

			t.Log(dp.Verify([]byte("abcc"), tt.Target, hashFn))
		})
	}

}
