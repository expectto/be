package psi

import (
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
)

// IsMatcher returns true if given input is either Omega or Gomock or a Psi matcher
func IsMatcher(a any) bool {
	switch a.(type) {
	case types.GomegaMatcher, types.GomockMatcher:
		return true
	default:
		return false
	}
}

// AsMatcher returns BeMatcher that is made from given input
func AsMatcher(m any) types.BeMatcher {
	switch t := m.(type) {
	case types.BeMatcher:
		return t
	case types.GomegaMatcher:
		return FromGomega(t)
	case types.GomockMatcher:
		return FromGomock(t)
	default:
		return FromGomega(gomega.Equal(t))
	}
}

// IsMatchFunc returns true if given `m any` is a match func (from gomega matchers)
// todo: support typed match funcs (see gomega)
func IsMatchFunc(m any) bool {
	_, ok := m.(func(any) (bool, error))
	return ok
}

// AsMatchFunc returns given `m any` as match func (func (any) (bool, error))
// todo: support typed match funcs (see gomega)
func AsMatchFunc(m any) func(any) (bool, error) {
	v, ok := m.(func(any) (bool, error))
	if !ok {
		panic("match func must be func(any) (bool, error)")
	}
	return v
}
