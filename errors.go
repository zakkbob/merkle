package merkle

import "errors"

var (
	ErrLeafNotFound = errors.New("leaf not found in merkle tree")
)
