package be

// matchers_be_common.go provides the everyday matchers that idiomatic Go unit
// tests reach for most — nil, booleans, errors, panics and collection membership.
// They wrap gomega's battle-tested implementations via Psi, exactly like Eq/Not.

import (
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"

	. "github.com/expectto/be/internal/psi" //nolint:staticcheck // dot-import is the established style here
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
)

// Nil succeeds if actual is nil. It is typed-nil aware (a nil *T inside an
// interface matches), unlike a bare `== nil` comparison.
func Nil() types.BeMatcher { return Psi(gomega.BeNil()) }

// NotNil succeeds if actual is not nil:
//
//	be.Expect(t, user).To(be.NotNil())
//
// Prefer this over be.Not(be.Nil()).
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

// MatchError succeeds if actual is an error matching expected. It is tri-mode —
// expected may be:
//   - a target error, compared with errors.Is (wrapping-aware):
//     be.Expect(t, err).To(be.MatchError(io.EOF))
//   - a string, compared against err.Error():
//     be.Expect(t, err).To(be.MatchError("file not found"))
//   - a matcher, applied to err.Error():
//     be.Expect(t, err).To(be.MatchError(be.ContainSubstring("not found")))
//
// Prefer be.MatchError(target) over be.True(errors.Is(err, target)), and
// be.Not(be.MatchError(target)) over be.False(errors.Is(err, target)).
// To match by error TYPE (errors.As), use MatchErrorAs.
func MatchError(expected any) types.BeMatcher { return Psi(gomega.MatchError(expected)) }

// MatchErrorAs succeeds if actual is an error that matches type T via
// errors.As — the matcher spelling of `var target T; errors.As(err, &target)`:
//
//	be.Expect(t, err).To(be.MatchErrorAs[*fs.PathError]())
//
// Prefer this over projecting through errors.As into be.True(). A nil or
// non-matching error fails; a non-error actual is an error (not a mismatch).
// To match by errors.Is target, message or matcher, use MatchError.
func MatchErrorAs[T error]() types.BeMatcher {
	typeName := reflect.TypeFor[T]().String()
	return Psi(func(actual any) (bool, error) {
		if actual == nil {
			return false, nil
		}
		err, ok := actual.(error)
		if !ok {
			return false, fmt.Errorf("MatchErrorAs matcher expects an error, got %T", actual)
		}
		var target T
		return errors.As(err, &target), nil
	}, fmt.Sprintf("match error as %s (via errors.As)", typeName))
}

// Panic succeeds if actual is a func() that panics when invoked.
func Panic() types.BeMatcher { return Psi(gomega.Panic()) }

// NotPanic succeeds if actual is a func() that does not panic when invoked.
func NotPanic() types.BeMatcher { return Psi(gomega.Not(gomega.Panic())) }

// ContainElement succeeds if actual (a slice, array or map) contains an element
// that matches the given value or matcher:
//
//	be.Expect(t, ids).To(be.ContainElement(42))
//	be.Expect(t, users).To(be.ContainElement(be.HaveField("Name", "Alice")))
//
// Prefer this over be.True(slices.Contains(xs, v)) — the failure message shows
// the collection. For substrings of a string use ContainSubstring.
func ContainElement(element any) types.BeMatcher { return Psi(gomega.ContainElement(element)) }

// ContainElements succeeds if actual contains all of the given elements (each may
// be a value or a matcher), in any order.
func ContainElements(elements ...any) types.BeMatcher {
	return Psi(gomega.ContainElements(elements...))
}

// HaveKey succeeds if actual (a map) has a key matching the given value or matcher:
//
//	be.Expect(t, headers).To(be.HaveKey("Authorization"))
//
// Prefer this over `_, ok := m[k]` followed by be.True(ok).
func HaveKey(key any) types.BeMatcher { return Psi(gomega.HaveKey(key)) }

// HaveKeyWithValue succeeds if actual (a map) has the given key with a matching value.
func HaveKeyWithValue(key, value any) types.BeMatcher {
	return Psi(gomega.HaveKeyWithValue(key, value))
}

// Ne succeeds if actual is NOT equal to expected (the negation of Eq):
//
//	be.Expect(t, status).To(be.Ne("failed"))
//
// Prefer this over be.Not(be.Eq(x)). To assert "not the zero value" use
// be.NonZero() instead of be.Ne(0).
func Ne(expected any) types.BeMatcher { return Psi(gomega.Not(gomega.Equal(expected))) }

// Zero succeeds if actual is the zero value for its type: 0, "", nil, false,
// a zero struct, etc. (reflect-based, works for any type — like gomega.BeZero):
//
//	be.Expect(t, count).To(be.Zero())
//	be.Expect(t, cfg).To(be.Zero()) // zero struct
//
// Prefer this over be.Eq(0) when you mean "unset". For the numeric-only
// spelling (where a non-number is an error, not a mismatch) use be_math.Zero.
func Zero() types.BeMatcher { return Psi(gomega.BeZero()) }

// NonZero succeeds if actual is NOT the zero value for its type:
//
//	be.Expect(t, userID).To(be.NonZero())
//
// Prefer this over be.Not(be.Eq(0)) or be.Ne(0).
func NonZero() types.BeMatcher { return Psi(gomega.Not(gomega.BeZero())) }

// Empty succeeds if actual is empty: a zero-length string, slice, array, map or
// channel (like gomega.BeEmpty):
//
//	be.Expect(t, errsList).To(be.Empty())
//
// Prefer this over be.HaveLength(0) or be.True(len(xs) == 0).
func Empty() types.BeMatcher { return Psi(gomega.BeEmpty()) }

// NotEmpty succeeds if actual is not empty:
//
//	be.Expect(t, results).To(be.NotEmpty())
//
// Prefer this over be.Not(be.HaveLength(0)) or be.True(len(xs) > 0).
func NotEmpty() types.BeMatcher { return Psi(gomega.Not(gomega.BeEmpty())) }

// ContainSubstring succeeds if actual is a string containing the given substring.
// (For slices/arrays/maps use ContainElement.)
func ContainSubstring(substr string) types.BeMatcher {
	return Psi(gomega.ContainSubstring(substr))
}

// Identical succeeds if actual is identical to expected using Go's == operator
// (pointer identity for pointers). Like gomega.BeIdenticalTo / testify's Same:
//
//	be.Expect(t, gotPtr).To(be.Identical(wantPtr)) // same pointer
//
// Footgun to know: be.Eq compares by VALUE (deep equality) — two different
// pointers to equal structs satisfy be.Eq but not be.Identical. Use Identical
// when "the same instance" is what you mean.
func Identical(expected any) types.BeMatcher { return Psi(gomega.BeIdenticalTo(expected)) }

// NotIdentical succeeds if actual is NOT identical to expected (the negation of
// Identical). Like testify's NotSame.
func NotIdentical(expected any) types.BeMatcher {
	return Psi(gomega.Not(gomega.BeIdenticalTo(expected)))
}

// HaveField succeeds if actual is a struct (or pointer to one) whose field —
// or nil-safe method chain — matches the given value or matcher. The field
// spec follows gomega.HaveField: a name, a dotted path, or a "Method()" call:
//
//	be.Expect(t, user).To(be.HaveField("Name", "Alice"))
//	be.Expect(t, user).To(be.HaveField("Address.City", be.NotEmpty()))
//	be.Expect(t, user).To(be.HaveField("ID()", be.NonZero()))
//
// This is the default struct-field matcher. The naming wobble with
// be_struct.HavingField is deliberate: HavingField[T] is the generic,
// compile-time-checked variant for when you want the struct type enforced.
func HaveField(field string, value any) types.BeMatcher {
	return Psi(gomega.HaveField(field, Psi(value)))
}

// HaveFields succeeds if actual matches HaveField for every entry of the given
// map (logical AND). Values may be raw values or matchers:
//
//	be.Expect(t, user).To(be.HaveFields(map[string]any{
//		"Name":  "Alice",
//		"Email": be_string.ValidEmail(),
//	}))
//
// Fields are checked in sorted-key order, so failure output is deterministic.
func HaveFields(fields map[string]any) types.BeMatcher {
	keys := slices.Sorted(maps.Keys(fields))
	ms := make([]any, len(keys))
	for i, k := range keys {
		ms[i] = HaveField(k, fields[k])
	}
	return All(ms...)
}

// Via applies the transform function to the actual value and matches the result
// against the given matcher. Handy for projecting through a public accessor when
// the underlying value can't be matched directly, e.g.:
//
//	be.Expect(t, ctx).To(be.Via(GetActor, be.Eq(wantActor)))
//
// transform must be a function of one argument returning one value (and
// optionally an error).
func Via(transform any, matcher any) types.BeMatcher {
	return Psi(gomega.WithTransform(transform, Psi(matcher)))
}
