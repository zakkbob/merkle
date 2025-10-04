package merkle_test

import (
	"crypto/sha256"
	"slices"
	"strconv"
	"testing"

	"github.com/zakkbob/merkle"
)

func TestNewTree(t *testing.T) {
	leaves := []string{}
	for i := range 100 {
		leaves = append(leaves, strconv.Itoa(i))
	}
	tree := merkle.NewTree(leaves, func(data []byte) []byte {
		b := sha256.Sum256(data)
		return b[:]
	})

	t.Log("\n" + tree.String())
}

func TestProve(t *testing.T) {
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

	tree := merkle.NewTree([]string{"a", "b", "c"}, hashFn)
	t.Log(tree.String())

	tests := []struct {
		Target   string
		Expected []string
	}{
		{
			Target:   "a",
			Expected: []string{"a", "b", "cc"},
		},
		{
			Target:   "b",
			Expected: []string{"b", "a", "cc"},
		},
		{
			Target:   "c",
			Expected: []string{"c", "c", "ab"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Target, func(t *testing.T) {
			p, err := tree.Prove(tt.Target)
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(p.Path(), tt.Expected) {
				t.Fatalf("got %v; expected %v", p, tt.Expected)
			}
			t.Log(p.Verify("abcc", tt.Target, hashFn))
		})
	}

}
