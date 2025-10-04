package merkle_test

import (
	"crypto/sha256"
	"slices"
	"testing"

	"github.com/zakkbob/merkle"
)

func TestNewTree(t *testing.T) {
	tree := merkle.NewTree([]string{"a", "b", "c", "d", "e", "f"}, func(data []byte) []byte {
		b := sha256.Sum256(data)
		return b[:]
	})

	t.Log("\n" + tree.String())
}

func TestProve(t *testing.T) {
	// abc
	// ├─ ab
	// │  ├─ a
	// │  └─ b
	// └─ c

	tree := merkle.NewTree([]string{"a", "b", "c"}, func(b []byte) []byte {
		return b
	})

	tests := []struct {
		Target   string
		Expected []string
	}{
		{
			Target:   "a",
			Expected: []string{"a", "b", "c"},
		},
		{
			Target:   "b",
			Expected: []string{"b", "a", "c"},
		},
		{
			Target:   "c",
			Expected: []string{"c", "ab"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Target, func(t *testing.T) {
			p := tree.Prove(tt.Target)
			if !slices.Equal(p, tt.Expected) {
				t.Fatalf("got %v; expected %v", p, tt.Expected)
			}
		})
	}

}
