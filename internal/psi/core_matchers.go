package psi

import (
	psiMatchers "github.com/expectto/be/internal/psi/matchers"
	"github.com/expectto/be/types"
)

// Always does always match
func Always() *psiMatchers.AlwaysMatcher {
	return &psiMatchers.AlwaysMatcher{}
}

// Never does never succeed (does always fail)
func Never(err error) *psiMatchers.NeverMatcher {
	return psiMatchers.NewNeverMatcher(err)
}

// All is like gomega.And()
func All(ms ...types.BeMatcher) *psiMatchers.AllMatcher {
	return psiMatchers.NewAllMatcher(ms...)
}

// Eq is like gomega.Equal()
func Eq(expected any) *psiMatchers.EqMatcher {
	return &psiMatchers.EqMatcher{Expected: expected}
}
