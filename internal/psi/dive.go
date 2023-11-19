package psi

import (
	"github.com/expectto/be/types"
)

// Dive does transformation so that given matcher is applied to each element of the slice.
//
// Todo: refactor: keep here only DiveMatcher, move Dive into be.go
// Deprecated
func Dive(matcher any) types.BeMatcher {
	// todo: implement
	return Psi(matcher)
}
