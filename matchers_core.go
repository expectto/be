package be

import (
	"github.com/expectto/be/internal/psi"
	psiMatchers "github.com/expectto/be/internal/psi/matchers"
	"github.com/expectto/be/matchers"
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

// Not is like gomega.Not()
func Not(matcher any) types.BeMatcher {
	return &psiMatchers.NotMatcher{Matcher: psi.Psi(matcher)}
}

// HaveLength is like gomega.HaveLen()
// HaveLength succeeds if the actual value has a length that matches the provided conditions.
// It accepts either a count value or one or more Gomega matchers to specify the desired length conditions.
// Todo move to other file?
func HaveLength(args ...any) types.BeMatcher {
	return matchers.NewHaveLengthMatcher(args...)
}
