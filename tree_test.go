package merkle_test

import (
	"slices"
	"testing"

	"github.com/zakkbob/merkle"
)

func TestProve(t *testing.T) {
	// a
	// ├─ b
	// │  ├─ c
	// │  └─ d
	// └─ e

	e := merkle.Node{Hash: "e"}
	d := merkle.Node{Hash: "d"}
	c := merkle.Node{Hash: "c"}
	b := merkle.Node{Hash: "b", Left: &c, Right: &d}
	a := merkle.Node{Hash: "a", Left: &b, Right: &e}

	tests := []struct {
		Target   string
		Expected []string
	}{
		{
			Target:   "e",
			Expected: []string{"e", "b"},
		},
		{
			Target:   "d",
			Expected: []string{"d", "c", "e"},
		},
		{
			Target:   "c",
			Expected: []string{"c", "d", "e"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Target, func(t *testing.T) {
			p := a.Prove(tt.Target)
			if !slices.Equal(p, tt.Expected) {
				t.Fatalf("got %v; expected %v", p, tt.Expected)
			}
		})
	}

}
