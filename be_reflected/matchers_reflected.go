// Package be_reflected offers Be matchers for reflection in Go,
// enabling expressive and versatile assertions on values' reflect kinds.
package be_reflected

import (
	"fmt"
	"github.com/expectto/be/internal/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
	"io"
	"reflect"
	"strconv"
)

// AsKind succeeds if actual is assignable to any of the specified kinds or matches the provided matchers.
func AsKind(args ...any) types.BeMatcher { return psi_matchers.NewKindMatcher(args...) }

// AssignableTo succeeds if actual is assignable to the specified type T.
func AssignableTo[T any]() types.BeMatcher { return psi_matchers.NewAssignableToMatcher[T]() }

// Implementing succeeds if actual implements the specified interface type T.
func Implementing[T any]() types.BeMatcher { return psi_matchers.NewImplementsMatcher[T]() }

// AsFunc succeeds if actual is of kind reflect.Func.
func AsFunc() types.BeMatcher { return AsKind(reflect.Func) }

// AsChan succeeds if actual is of kind reflect.Chan.
func AsChan() types.BeMatcher { return AsKind(reflect.Chan) }

// AsPointer succeeds if the actual value is a pointer.
func AsPointer() types.BeMatcher { return AsKind(reflect.Pointer) }

// AsFinalPointer succeeds if the actual value is a final pointer, meaning it's a pointer to a non-pointer type.
func AsFinalPointer() types.BeMatcher {
	return &psi_matchers.AllMatcher{Matchers: []types.BeMatcher{
		AsPointer(),
		WithFallibleTransform(func(actual any) any {
			return reflect.ValueOf(actual).Elem()
		}, psi_matchers.NewNotMatcher(AsPointer())),
	}}
}

// AsStruct succeeds if actual is of kind reflect.Struct.
func AsStruct() types.BeMatcher { return AsKind(reflect.Struct) }

// AsPointerToStruct succeeds if actual is a pointer to a struct.
func AsPointerToStruct() types.BeMatcher {
	return &psi_matchers.AllMatcher{Matchers: []types.BeMatcher{
		AsPointer(),
		WithFallibleTransform(func(actual any) any {
			return reflect.ValueOf(actual).Elem()
		}, AsStruct()),
	}}
}

// AsSlice succeeds if actual is of kind reflect.Slice.
func AsSlice() types.BeMatcher {
	return AsKind(reflect.Slice)
}

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
	return Psi(
		AsKind(reflect.Slice),
		gomega.HaveEach(AssignableTo[T]()),
	)
}

// AsMap succeeds if actual is of kind reflect.Map.
func AsMap() types.BeMatcher { return AsKind(reflect.Map) }

// AsPointerToMap succeeds if actual is a pointer to a map.
func AsPointerToMap() types.BeMatcher {
	return &psi_matchers.AllMatcher{Matchers: []types.BeMatcher{
		AsPointer(),
		WithFallibleTransform(func(actual any) any {
			return reflect.ValueOf(actual).Elem()
		}, AsMap()),
	}}
}

// AsObject is more specific than AsMap. It checks if the given `actual` value is a map with string keys
// and values of any type. This is particularly useful in the context of BeJson matcher,
// where the term 'Object' aligns with JSON notation.
func AsObject() types.BeMatcher {
	return &psi_matchers.AllMatcher{Matchers: []types.BeMatcher{
		AsKind(reflect.Map), AssignableTo[map[string]any](),
	}}
}
func AsObjects() types.BeMatcher {
	return AsSliceOf[map[string]any]()
}

// AsPointerToObject succeeds if actual is a pointer to a value that matches AsObject after applying dereference.
func AsPointerToObject() types.BeMatcher {
	return &psi_matchers.AllMatcher{Matchers: []types.BeMatcher{
		AsPointer(),
		WithFallibleTransform(func(actual any) any {
			return reflect.ValueOf(actual).Elem()
		}, AsObject()),
	}}
}

// AsReader succeeds if actual implements the io.Reader interface.
func AsReader() types.BeMatcher { return Implementing[io.Reader]() }

// AsStringer succeeds if actual implements the fmt.Stringer interface.
func AsStringer() types.BeMatcher { return Implementing[fmt.Stringer]() }

// AsString succeeds if actual is of kind reflect.String.
func AsString() types.BeMatcher { return AsKind(reflect.String) }

// AsBytes succeeds if actual is assignable to a slice of bytes ([]byte).
func AsBytes() types.BeMatcher { return AssignableTo[[]byte]() }

// AsNumeric succeeds if actual is a numeric value, supporting various
// integer kinds: reflect.Int, ... reflect.Int64,
// and floating-point kinds: reflect.Float32, reflect.Float64
func AsNumeric() types.BeMatcher {
	return AsKind(
		gomega.BeNumerically(">=", reflect.Int),
		gomega.BeNumerically("<=", reflect.Float64),
	)
}

// AsNumericString succeeds if actual is a string that can be parsed into a valid numeric value.
func AsNumericString() types.BeMatcher {
	return &psi_matchers.AllMatcher{Matchers: []types.BeMatcher{
		AsString(),
		WithFallibleTransform(func(actual any) any {
			_, err := strconv.ParseFloat(cast.AsString(actual), 64)
			return err == nil
		}, gomega.BeTrue()),
	}}
}

// AsIntish succeeds if actual is a numeric value that represents an integer (from reflect.Int up to reflect.Uint64).
func AsIntish() types.BeMatcher {
	return AsKind(
		gomega.BeNumerically(">=", reflect.Int),
		gomega.BeNumerically("<=", reflect.Uint64),
	)
}

// AsIntishString succeeds if actual is a string that can be parsed into a valid integer value.
func AsIntishString() types.BeMatcher {
	return &psi_matchers.AllMatcher{Matchers: []types.BeMatcher{
		AsString(),
		WithFallibleTransform(func(actual any) any {
			_, err := strconv.ParseInt(cast.AsString(actual), 10, 64)
			return err == nil
		}, gomega.BeTrue()),
	}}
}

// AsFloatish succeeds if actual is a numeric value that represents a floating-point value.
func AsFloatish() types.BeMatcher {
	return AsKind(
		gomega.BeNumerically(">=", reflect.Float32),
		gomega.BeNumerically("<=", reflect.Float64),
	)
}

// AsFloatishString succeeds if actual is a string that can be parsed into a valid floating-point value.
func AsFloatishString() types.BeMatcher {
	return &psi_matchers.AllMatcher{Matchers: []types.BeMatcher{
		AsString(),
		WithFallibleTransform(func(actual any) any {
			_, err := strconv.ParseFloat(cast.AsString(actual), 64)
			return err == nil
		}, gomega.BeTrue()),
	}}
}
