package be

// matchers_be.go contains Public callers for core psi matchers
// For advances matchers check out `be_*` packages

import (
	. "github.com/expectto/be/internal/psi" //nolint:staticcheck // should be moved to lintignore
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

// Eq succeeds if actual equals expected by VALUE (deep equality, like
// gomega.Equal):
//
//	be.Expect(t, got).To(be.Eq(want))
//
// Footgun to know: two different pointers to equal structs satisfy be.Eq.
// When you mean "the same instance" (pointer identity, Go's ==), use
// be.Identical instead. For "unset / zero value" prefer be.Zero over be.Eq(0).
func Eq(expected any) types.BeMatcher {
	return psi_matchers.NewEqMatcher(expected)
}

// Not is like gomega.Not()
func Not(expected any) types.BeMatcher {
	return psi_matchers.NewNotMatcher(Psi(expected))
}

// HaveLength succeeds if the actual value (string, slice, array, map or
// channel) has a length matching the provided condition — either an exact
// count, or one or more matchers applied to the length (unlike gomega.HaveLen,
// which only takes a count):
//
//	be.Expect(t, items).To(be.HaveLength(3))
//	be.Expect(t, items).To(be.HaveLength(be.Gte(3)))          // composable form
//	be.Expect(t, name).To(be.HaveLength(be.InRange(1, true, 64, true)))
//
// Prefer be.HaveLength(be.Gte(n)) over be.True(len(xs) >= n). For zero /
// non-zero length prefer be.Empty / be.NotEmpty.
func HaveLength(args ...any) types.BeMatcher {
	return psi_matchers.NewHaveLengthMatcher(args...)
}

// Dive applies the given matcher to each (every) element of a slice or array,
// or to each value of a map.
// Note: Dive is very close to gomega.HaveEach
func Dive(matcher any) types.BeMatcher { return NewDiveMatcher(matcher, DiveModeEvery) }

// DiveAny applies the given matcher to each element and succeeds in case if it succeeds at least at one item
func DiveAny(matcher any) types.BeMatcher { return NewDiveMatcher(matcher, DiveModeAny) }

// DiveFirst applies the given matcher to the first element of the given slice
func DiveFirst(matcher any) types.BeMatcher { return NewDiveMatcher(matcher, DiveModeFirst) }

// DiveNth applies the given matcher to the nth element of the given slice
func DiveNth(n int, matcher any) types.BeMatcher { return NewDiveMatcher(matcher, DiveModeNth, n) }
