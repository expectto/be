package be

import (
	"fmt"
	"github.com/expectto/be/internal/cast"
	"github.com/expectto/be/internal/psi"
	"github.com/expectto/be/matchers"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
	"io"
	"reflect"
	"strconv"
)

// Kind succeeds if the actual value has a kind that satisfies given conditions
// Parameter can be a raw reflect.Kind value, or one or more matchers
func Kind(args ...any) types.BeMatcher {
	return matchers.NewKindMatcher(args...)
}

func AssignableTo[T any]() types.BeMatcher { return matchers.NewAssignableToMatcher[T]() }
func Implement[T any]() types.BeMatcher    { return matchers.NewImplementsMatcher[T]() }

func Func() types.BeMatcher   { return Kind(reflect.Func) }
func Chan() types.BeMatcher   { return Kind(reflect.Chan) }
func Struct() types.BeMatcher { return Kind(reflect.Struct) }
func Slice() types.BeMatcher  { return Kind(reflect.Slice) }
func SliceOf[T any]() types.BeMatcher {
	return psi.Psi(
		Kind(reflect.Slice),
		gomega.HaveEach(AssignableTo[T]()),
	)
}

// Map matches any kind of map
func Map() types.BeMatcher { return Kind(reflect.Map) }

// Object is more specific than Map. It checks if given `actual` value is a map[string]any-compatible value
// Is considered to be used in BeJson matcher (as term Object is more relevant to JSON notation)
func Object() types.BeMatcher {
	return All(
		Kind(reflect.Map), AssignableTo[map[string]any](),
	)
}
func Objects() types.BeMatcher {
	return SliceOf[map[string]any]()
}

// Pointer succeeds if the actual value is a pointer.
func Pointer() types.BeMatcher { return Kind(reflect.Pointer) }

// FinalPointer succeeds if the actual value is a final pointer,
// meaning it's a pointer to a non-pointer type.
func FinalPointer() types.BeMatcher {
	return All(
		Pointer(),
		psi.WithFallibleTransform(func(actual any) any {
			return reflect.ValueOf(actual).Elem()
		}, gomega.Not(Pointer())),
	)
}

func Reader() types.BeMatcher   { return Implement[io.Reader]() }
func Stringer() types.BeMatcher { return Implement[fmt.Stringer]() }

func String() types.BeMatcher { return Kind(reflect.String) }
func Bytes() types.BeMatcher  { return AssignableTo[[]byte]() }

func Numeric() types.BeMatcher { return Kind(Gte(reflect.Int), Lte(reflect.Float64)) }
func Intish() types.BeMatcher  { return Kind(Gte(reflect.Int), Lte(reflect.Uint64)) }
func Floatish() types.BeMatcher {
	return Kind(Gte(reflect.Float32), Lte(reflect.Float64))
}

func NumericString() types.BeMatcher {
	return All(
		String(),
		psi.WithFallibleTransform(func(actual any) any {
			_, err := strconv.ParseFloat(cast.AsString(actual), 64)
			return err == nil
		}, gomega.BeTrue()),
	)
}
