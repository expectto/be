package psi

import (
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gcustom"
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

// WithCustomMessage is a wrapper for gcustom.MakeMatcher
// todo: make `v any` so we can check here if it's types.BeMatcher of match-func
func WithCustomMessage(v types.BeMatcher, message string) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(v.Match, message))
}
