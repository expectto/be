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

func AsKind(args ...any) types.BeMatcher { return psi_matchers.NewKindMatcher(args...) }

func AssignableTo[T any]() types.BeMatcher { return psi_matchers.NewAssignableToMatcher[T]() }
func Implementing[T any]() types.BeMatcher { return psi_matchers.NewImplementsMatcher[T]() }

func AsFunc() types.BeMatcher   { return AsKind(reflect.Func) }
func AsChan() types.BeMatcher   { return AsKind(reflect.Chan) }
func AsStruct() types.BeMatcher { return AsKind(reflect.Struct) }
func AsSlice() types.BeMatcher  { return AsKind(reflect.Slice) }
func AsSliceOf[T any]() types.BeMatcher {
	return Psi(
		AsKind(reflect.Slice),
		gomega.HaveEach(AssignableTo[T]()),
	)
}

// AsMap matches any kind of map
func AsMap() types.BeMatcher { return AsKind(reflect.Map) }

// AsObject is more specific than Map. It checks if given `actual` value is a map[string]any-compatible value
// Is considered to be used in BeJson matcher (as term Object is more relevant to JSON notation)
func AsObject() types.BeMatcher {
	return &psi_matchers.AllMatcher{Matchers: []types.BeMatcher{
		AsKind(reflect.Map), AssignableTo[map[string]any](),
	}}
}
func AsObjects() types.BeMatcher {
	return AsSliceOf[map[string]any]()
}

// AsPointer succeeds if the actual value is a pointer.
func AsPointer() types.BeMatcher { return AsKind(reflect.Pointer) }

// AsFinalPointer succeeds if the actual value is a final pointer,
// meaning it's a pointer to a non-pointer type.
func AsFinalPointer() types.BeMatcher {
	return &psi_matchers.AllMatcher{Matchers: []types.BeMatcher{
		AsPointer(),
		WithFallibleTransform(func(actual any) any {
			return reflect.ValueOf(actual).Elem()
		}, psi_matchers.NewNotMatcher(AsPointer())),
	}}
}

func AsReader() types.BeMatcher   { return Implementing[io.Reader]() }
func AsStringer() types.BeMatcher { return Implementing[fmt.Stringer]() }

func AsString() types.BeMatcher { return AsKind(reflect.String) }
func AsBytes() types.BeMatcher  { return AssignableTo[[]byte]() }

func AsNumeric() types.BeMatcher {
	return AsKind(
		gomega.BeNumerically(">=", reflect.Int),
		gomega.BeNumerically("<=", reflect.Float64),
	)
}
func AsIntish() types.BeMatcher {
	return AsKind(
		gomega.BeNumerically(">=", reflect.Int),
		gomega.BeNumerically("<=", reflect.Uint64),
	)
}
func AsFloatish() types.BeMatcher {
	return AsKind(
		gomega.BeNumerically(">=", reflect.Float32),
		gomega.BeNumerically("<=", reflect.Float64),
	)
}

func AsNumericString() types.BeMatcher {
	return &psi_matchers.AllMatcher{Matchers: []types.BeMatcher{
		AsString(),
		WithFallibleTransform(func(actual any) any {
			_, err := strconv.ParseFloat(cast.AsString(actual), 64)
			return err == nil
		}, gomega.BeTrue()),
	}}
}
