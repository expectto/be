package be

// matchers_be_common.go provides the everyday matchers that idiomatic Go unit
// tests reach for most — nil, booleans, errors, panics and collection membership.
// They wrap gomega's battle-tested implementations via Psi, exactly like Eq/Not.

import (
	. "github.com/expectto/be/internal/psi" //nolint:staticcheck // dot-import is the established style here
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
)

// Nil succeeds if actual is nil. It is typed-nil aware (a nil *T inside an
// interface matches), unlike a bare `== nil` comparison.
func Nil() types.BeMatcher { return Psi(gomega.BeNil()) }

// NotNil succeeds if actual is not nil.
func NotNil() types.BeMatcher { return Psi(gomega.Not(gomega.BeNil())) }

// True succeeds if actual is the boolean true.
func True() types.BeMatcher { return Psi(gomega.BeTrue()) }

// False succeeds if actual is the boolean false.
func False() types.BeMatcher { return Psi(gomega.BeFalse()) }

// Succeed succeeds if actual is a nil error. Intended for error values:
//
//	be.Expect(t, err).To(be.Succeed())
func Succeed() types.BeMatcher { return Psi(gomega.Succeed()) }

// HaveOccurred succeeds if actual is a non-nil error.
func HaveOccurred() types.BeMatcher { return Psi(gomega.HaveOccurred()) }

// MatchError succeeds if actual is an error matching expected, which may be:
//   - a target error (compared with errors.Is),
//   - a string (compared against err.Error()),
//   - a matcher applied to the error.
func MatchError(expected any) types.BeMatcher { return Psi(gomega.MatchError(expected)) }

// Panic succeeds if actual is a func() that panics when invoked.
func Panic() types.BeMatcher { return Psi(gomega.Panic()) }

// NotPanic succeeds if actual is a func() that does not panic when invoked.
func NotPanic() types.BeMatcher { return Psi(gomega.Not(gomega.Panic())) }

// ContainElement succeeds if actual (a slice, array or map) contains an element
// that matches the given value or matcher.
func ContainElement(element any) types.BeMatcher { return Psi(gomega.ContainElement(element)) }

// ContainElements succeeds if actual contains all of the given elements (each may
// be a value or a matcher), in any order.
func ContainElements(elements ...any) types.BeMatcher {
	return Psi(gomega.ContainElements(elements...))
}

// HaveKey succeeds if actual (a map) has a key matching the given value or matcher.
func HaveKey(key any) types.BeMatcher { return Psi(gomega.HaveKey(key)) }

// HaveKeyWithValue succeeds if actual (a map) has the given key with a matching value.
func HaveKeyWithValue(key, value any) types.BeMatcher {
	return Psi(gomega.HaveKeyWithValue(key, value))
}
