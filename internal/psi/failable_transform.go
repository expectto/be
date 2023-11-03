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
// todo: update docs here to be consistent wigh gomega's transforms docs
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

	numout := txType.NumOut()
	if numout == 1 {
		return true
	}
	if numout == 2 {
		return txType.Out(1).AssignableTo(reflect2.TypeFor[error]())
	}

	return false
}

// WithFallibleTransform creates a gomega transform matcher that can nicely handle failures
func WithFallibleTransform(transform any, matcher gomega.OmegaMatcher) types.BeMatcher {
	return Psi(gomega.WithTransform(transform, gomega.And(WithTransformError(), matcher)))
}

// transformErrorMatcher is actually a matcher
type transformErrorMatcher struct {
	actual any
	err    error
}

func WithTransformError() *transformErrorMatcher {
	return &transformErrorMatcher{}
}

func (matcher *transformErrorMatcher) Match(actual any) (success bool, err error) {
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

func (matcher *transformErrorMatcher) FailureMessage(actual any) string {
	return fmt.Sprintf("Expected\n%s\nto %s", format.Object(matcher.actual, 1), matcher.err)
}

func (matcher *transformErrorMatcher) NegatedFailureMessage(actual any) string {
	return fmt.Sprintf("Expected\n%s\nnot to %s", format.Object(matcher.actual, 1), matcher.err)
}

/*
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
*/
