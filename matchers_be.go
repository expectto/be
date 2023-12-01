package be

// matchers_be.go contains Public callers for core psi matchers
// For advances matchers check out `be_*` packages

import (
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
)

// Always does always match
func Always() types.BeMatcher {
	return psi_matchers.NewAlwaysMatcher()
}

// Never does never succeed (does always fail)
func Never(err error) types.BeMatcher {
	return psi_matchers.NewNeverMatcher(err)
}

// All is like gomega.And()
func All(ms ...any) types.BeMatcher {
	return psi_matchers.NewAllMatcher(Psi(ms...))
}

// Any is like gomega.Or()
func Any(ms ...any) types.BeMatcher {
	return psi_matchers.NewAnyMatcher(ms...)
}

// Eq is like gomega.Equal()
func Eq(expected any) types.BeMatcher {
	return psi_matchers.NewEqMatcher(expected)
}

// Not is like gomega.Not()
func Not(expected any) types.BeMatcher {
	return psi_matchers.NewNotMatcher(Psi(expected))
}

// HaveLength is like gomega.HaveLen()
// HaveLength succeeds if the actual value has a length that matches the provided conditions.
// It accepts either a count value or one or more Gomega matchers to specify the desired length conditions.
func HaveLength(args ...any) types.BeMatcher {
	return psi_matchers.NewHaveLengthMatcher(args...)
}
