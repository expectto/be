package psi

import (
	"fmt"
	"reflect"

	"github.com/expectto/be/types"
)

// WithTransformMatcher is taken from gomega
// reference: https://github.com/onsi/gomega/blob/master/matchers/with_transform.go
type WithTransformMatcher struct {
	// input
	Transform interface{} // must be a function of one parameter that returns one value and an optional error
	Matcher   types.BeMatcher

	// cached value
	transformArgType reflect.Type

	// state
	transformedValue interface{}
}

// WithTransform applies the `transform` to the actual value and matches it against `matcher`.
// The given transform must be either a function of one parameter that returns one value or a
// function of one parameter that returns two values, where the second value must be of the
// error type.
//
//	var plus1 = func(i int) int { return i + 1 }
//	Expect(1).To(WithTransform(plus1, Equal(2))
//
//	 var failingplus1 = func(i int) (int, error) { return 42, "this does not compute" }
//	 Expect(1).To(WithTransform(failingplus1, Equal(2)))
//
// And(), Or(), Not() and WithTransform() allow matchers to be composed into complex expressions.
func WithTransform(transform any, matcher types.BeMatcher) types.GomegaMatcher {
	return NewWithTransformMatcher(transform, matcher)
}

// reflect.Type for error
var errorT = reflect.TypeOf((*error)(nil)).Elem()

func NewWithTransformMatcher(transform interface{}, matcher types.BeMatcher) *WithTransformMatcher {
	if transform == nil {
		panic("transform function cannot be nil")
	}
	txType := reflect.TypeOf(transform)
	if txType.NumIn() != 1 {
		panic("transform function must have 1 argument")
	}
	if numout := txType.NumOut(); numout != 1 {
		if numout != 2 || !txType.Out(1).AssignableTo(errorT) {
			panic("transform function must either have 1 return value, or 1 return value plus 1 error value")
		}
	}

	return &WithTransformMatcher{
		Transform:        transform,
		Matcher:          matcher,
		transformArgType: reflect.TypeOf(transform).In(0),
	}
}

func (m *WithTransformMatcher) Match(actual interface{}) (bool, error) {
	// prepare a parameter to pass to the Transform function
	var param reflect.Value
	if actual != nil && reflect.TypeOf(actual).AssignableTo(m.transformArgType) {
		// The dynamic type of actual is compatible with the transform argument.
		param = reflect.ValueOf(actual)

	} else if actual == nil && m.transformArgType.Kind() == reflect.Interface {
		// The dynamic type of actual is unknown, so there's no way to make its
		// reflect.Value. Create a nil of the transform argument, which is known.
		param = reflect.Zero(m.transformArgType)

	} else {
		return false, fmt.Errorf("Transform function expects '%s' but we have '%T'", m.transformArgType, actual)
	}

	// call the Transform function with `actual`
	fn := reflect.ValueOf(m.Transform)
	result := fn.Call([]reflect.Value{param})
	if len(result) == 2 {
		if !result[1].IsNil() {
			return false, fmt.Errorf("Transform function failed: %s", result[1].Interface().(error).Error())
		}
	}
	m.transformedValue = result[0].Interface() // expect exactly one value

	return m.Matcher.Match(m.transformedValue)
}

func (m *WithTransformMatcher) FailureMessage(_ interface{}) (message string) {
	return m.Matcher.FailureMessage(m.transformedValue)
}

func (m *WithTransformMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return m.Matcher.NegatedFailureMessage(m.transformedValue)
}

func (m *WithTransformMatcher) MatchMayChangeInTheFuture(_ interface{}) bool {
	panic("not implemented")
}
