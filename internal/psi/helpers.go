package psi

import (
	psiMatchers "github.com/expectto/be/internal/psi/matchers"
	"github.com/expectto/be/types"
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
		return &psiMatchers.EqMatcher{Expected: m}
	}
}
