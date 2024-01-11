package psi

import (
	"fmt"
	reflect2 "github.com/expectto/be/internal/reflect"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"reflect"
)

// IsTransformFunc checks if given thing is a Gomega-compatible transform
// For v to be a transform it must be a function of one parameter that returns one value and an optional error
func IsTransformFunc(v any) bool {
	if v == nil {
		return false
	}
	txType := reflect.TypeOf(v)
	if txType.Kind() != reflect.Func {
		return false
	}
	if txType.NumIn() != 1 {
		return false
	}

	numOut := txType.NumOut()
	if numOut == 1 {
		return true
	}
	if numOut == 2 {
		return txType.Out(1).AssignableTo(reflect2.TypeFor[error]())
	}

	return false
}

// WithFallibleTransform creates a gomega transform matcher that can nicely handle failures
// Also it allows to have nil matcher, meaning that we're OK unless transform failed
func WithFallibleTransform(transform any, matcher gomega.OmegaMatcher) types.BeMatcher {
	if matcher != nil {
		matcher = gomega.And(WithTransformError(), matcher)
	} else {
		matcher = WithTransformError()
	}

	return Psi(gomega.WithTransform(transform, matcher))
}

// TransformErrorMatcher is actually a matcher
type TransformErrorMatcher struct {
	actual any
	err    error
}

func WithTransformError() *TransformErrorMatcher {
	return &TransformErrorMatcher{}
}

func (matcher *TransformErrorMatcher) Match(actual any) (success bool, err error) {
	if err, ok := actual.(error); ok {
		matcher.err = err
	}

	// Fill in actual value for future messages
	if h, ok := actual.(interface {
		Actual() any
	}); ok {
		matcher.actual = h.Actual()
	}

	return matcher.err == nil, nil
}

func (matcher *TransformErrorMatcher) FailureMessage(actual any) string {
	return fmt.Sprintf("Expected\n%s\nto %s", format.Object(matcher.actual, 1), matcher.err)
}

func (matcher *TransformErrorMatcher) NegatedFailureMessage(actual any) string {
	return fmt.Sprintf("Expected\n%s\nnot to %s", format.Object(matcher.actual, 1), matcher.err)
}

// TransformError is used to store error + actual value which caused the error
type TransformError struct {
	error
	actual any
}

func NewTransformError(err error, actual any) *TransformError {
	return &TransformError{error: err, actual: actual}
}

func (terr *TransformError) Actual() any {
	return terr.actual
}
