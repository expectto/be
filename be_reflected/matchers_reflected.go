// Package be_reflected provides Be matchers that use reflection,
// enabling expressive assertions on values' reflect kinds and types.
// It consists of several "core" matchers e.g. AsKind / AssignableTo / Implementing
// And many other matchers that are made on-top on core ones. E.g. AsFunc / AsString / AsNumber / etc
package be_reflected

import (
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/expectto/be/internal/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	reflect2 "github.com/expectto/be/internal/reflect"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
)

// AsKind succeeds if actual is assignable to any of the specified kinds or matches the provided matchers.
func AsKind(args ...any) types.BeMatcher { return psi_matchers.NewKindMatcher(args...) }

// AssignableTo succeeds if actual is assignable to the specified type T.
func AssignableTo[T any]() types.BeMatcher { return psi_matchers.NewAssignableToMatcher[T]() }

// Implementing succeeds if actual implements the specified interface type T.
func Implementing[T any]() types.BeMatcher { return psi_matchers.NewImplementsMatcher[T]() }

// Following matchers below are nice syntax-sugar, pretty usages of core matchers above:

// AsFunc succeeds if actual is of kind reflect.Func.
func AsFunc() types.BeMatcher { return Psi(AsKind(reflect.Func), "be a func") }

// AsChan succeeds if actual is of kind reflect.Chan.
func AsChan() types.BeMatcher { return Psi(AsKind(reflect.Chan), "be a channel") }

// AsPointer succeeds if the actual value is a pointer.
func AsPointer() types.BeMatcher { return Psi(AsKind(reflect.Pointer), "be a pointer") }

// AsFinalPointer succeeds if the actual value is a final pointer, meaning it's a pointer to a non-pointer type.
func AsFinalPointer() types.BeMatcher {
	return Psi(psi_matchers.NewAllMatcher(
		AsPointer(),
		WithFallibleTransform(func(actual any) any {
			return reflect.ValueOf(actual).Elem()
		}, psi_matchers.NewNotMatcher(AsPointer())),
	), "be a final pointer")
}

// AsStruct succeeds if actual is of kind reflect.Struct.
func AsStruct() types.BeMatcher { return Psi(AsKind(reflect.Struct), "be a struct") }

// AsPointerToStruct succeeds if actual is a pointer to a struct.
func AsPointerToStruct() types.BeMatcher {
	return Psi(psi_matchers.NewAllMatcher(
		AsPointer(),
		WithFallibleTransform(func(actual any) any {
			return reflect.ValueOf(actual).Elem()
		}, AsStruct()),
	), "be a pointer to a struct")
}

// AsSlice succeeds if actual is of kind reflect.Slice.
func AsSlice() types.BeMatcher { return Psi(AsKind(reflect.Slice), "be a slice") }

// AsPointerToSlice succeeds if actual is a pointer to a slice.
func AsPointerToSlice() types.BeMatcher {
	return &psi_matchers.AllMatcher{Matchers: []types.BeMatcher{
		AsPointer(),
		WithFallibleTransform(func(actual any) any {
			return reflect.ValueOf(actual).Elem()
		}, AsSlice()),
	}}
}

// AsSliceOf succeeds if actual is of kind reflect.Slice and each element of the slice
// is assignable to the specified type T.
func AsSliceOf[T any]() types.BeMatcher {
	return Psi(psi_matchers.NewAllMatcher(
		AsKind(reflect.Slice),
		gomega.HaveEach(AssignableTo[T]()),
	), "be a slice of "+reflect2.TypeFor[T]().String())
}

// AsMap succeeds if actual is of kind reflect.Map.
func AsMap() types.BeMatcher { return Psi(AsKind(reflect.Map), "be a map") }

// AsPointerToMap succeeds if actual is a pointer to a map.
func AsPointerToMap() types.BeMatcher {
	return Psi(psi_matchers.NewAllMatcher(
		AsPointer(),
		WithFallibleTransform(func(actual any) any {
			return reflect.ValueOf(actual).Elem()
		}, AsMap()),
	), "be a pointer to a map")
}

// AsObject is more specific than AsMap. It checks if the given `actual` value is a map with string keys
// and values of any type. This is particularly useful in the context of BeJson matcher,
// where the term 'Object' aligns with JSON notation.
func AsObject() types.BeMatcher {
	return Psi(psi_matchers.NewAllMatcher(
		AsKind(reflect.Map), AssignableTo[map[string]any](),
	), "be an object")
}
func AsObjects() types.BeMatcher {
	return Psi(AsSliceOf[map[string]any](), "be objects")
}

// AsPointerToObject succeeds if actual is a pointer to a value that matches AsObject after applying dereference.
func AsPointerToObject() types.BeMatcher {
	return Psi(psi_matchers.NewAllMatcher(
		AsPointer(),
		WithFallibleTransform(func(actual any) any {
			return reflect.ValueOf(actual).Elem()
		}, AsObject()),
	), "be a pointer to an object")
}

// AsReader succeeds if actual implements the io.Reader interface.
func AsReader() types.BeMatcher {
	return Psi(Implementing[io.Reader](), "implement io.Reader interface")
}

// AsStringer succeeds if actual implements the fmt.Stringer interface.
func AsStringer() types.BeMatcher {
	return Psi(Implementing[fmt.Stringer](), "implement fmt.Stringer interface")
}

// AsString succeeds if actual is of kind reflect.String.
func AsString() types.BeMatcher { return Psi(AsKind(reflect.String), "be a string") }

// AsBytes succeeds if actual is assignable to a slice of bytes ([]byte).
func AsBytes() types.BeMatcher { return Psi(AssignableTo[[]byte](), "be bytes") }

// AsNumber succeeds if actual is a numeric value, supporting various
// integer kinds: reflect.Int, ... reflect.Int64,
// and floating-point kinds: reflect.Float32, reflect.Float64
func AsNumber() types.BeMatcher {
	return Psi(AsKind(
		gomega.BeNumerically(">=", reflect.Int),
		gomega.BeNumerically("<=", reflect.Float64),
	), "be a number")
}

// AsNumericString succeeds if actual is a string that can be parsed into a valid numeric value.
func AsNumericString() types.BeMatcher {
	return Psi(psi_matchers.NewAllMatcher(
		AsString(),
		WithFallibleTransform(func(actual any) any {
			_, err := strconv.ParseFloat(cast.AsString(actual), 64)
			return err == nil
		}, gomega.BeTrue()),
	), "be a numeric string")
}

// AsInteger succeeds if actual is a numeric value that represents an integer (from reflect.Int up to reflect.Uint64).
func AsInteger() types.BeMatcher {
	return Psi(AsKind(
		gomega.BeNumerically(">=", reflect.Int),
		gomega.BeNumerically("<=", reflect.Uint64),
	), "be an integer value")
}

// AsIntegerString succeeds if actual is a string that can be parsed into a valid integer value.
func AsIntegerString() types.BeMatcher {
	return Psi(psi_matchers.NewAllMatcher(
		AsString(),
		WithFallibleTransform(func(actual any) any {
			_, err := strconv.ParseInt(cast.AsString(actual), 10, 64)
			return err == nil
		}, gomega.BeTrue()),
	), "be an integer-ish string")
}

// AsFloat succeeds if actual is a numeric value that represents a floating-point value.
func AsFloat() types.BeMatcher {
	return Psi(AsKind(
		gomega.BeNumerically(">=", reflect.Float32),
		gomega.BeNumerically("<=", reflect.Float64),
	), "be a float value")
}

// AsFloatString succeeds if actual is a string that can be parsed into a valid floating-point value.
func AsFloatString() types.BeMatcher {
	return Psi(psi_matchers.NewAllMatcher(
		AsString(),
		WithFallibleTransform(func(actual any) any {
			_, err := strconv.ParseFloat(cast.AsString(actual), 64)
			return err == nil
		}, gomega.BeTrue()),
	), "be a float-ish string")
}
