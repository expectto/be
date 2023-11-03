package be

import (
	psiMatchers "github.com/expectto/be/internal/psi/matchers"
	"github.com/expectto/be/types"
)

// Always does always match
func Always() types.BeMatcher {
	return &psiMatchers.AlwaysMatcher{}
}

// Never does never succeed (does always fail)
func Never(err error) types.BeMatcher {
	return psiMatchers.NewNeverMatcher(err)
}

// All is like gomega.And()
func All(ms ...types.BeMatcher) types.BeMatcher {
	return psiMatchers.NewAllMatcher(ms...)
}

// Eq is like gomega.Equal()
func Eq(expected any) types.BeMatcher {
	return &psiMatchers.EqMatcher{Expected: expected}
}
