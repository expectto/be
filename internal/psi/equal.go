package psi

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/expectto/be/internal/tmp/format"
	"github.com/expectto/be/types"
)

// Equal uses reflect.DeepEqual to compare actual with expected.  Equal is strict about
// types when performing comparisons.
// It is an error for both actual and expected to be nil.  Use BeNil() instead.
func Equal(expected interface{}) types.BeMatcher {
	return &EqualMatcher{
		Expected: expected,
	}
}

type EqualMatcher struct {
	Expected interface{}
}

func (matcher *EqualMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil && matcher.Expected == nil {
		return false, fmt.Errorf("Refusing to compare <nil> to <nil>.\nBe explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized.")
	}
	// Shortcut for byte slices.
	// Comparing long byte slices with reflect.DeepEqual is very slow,
	// so use bytes.Equal if actual and expected are both byte slices.
	if actualByteSlice, ok := actual.([]byte); ok {
		if expectedByteSlice, ok := matcher.Expected.([]byte); ok {
			return bytes.Equal(actualByteSlice, expectedByteSlice), nil
		}
	}
	return reflect.DeepEqual(actual, matcher.Expected), nil
}

func (matcher *EqualMatcher) FailureMessage(actual interface{}) (message string) {
	actualString, actualOK := actual.(string)
	expectedString, expectedOK := matcher.Expected.(string)
	if actualOK && expectedOK {
		return format.MessageWithDiff(actualString, "to equal", expectedString)
	}

	return format.Message(actual, "to equal", matcher.Expected)
}

func (matcher *EqualMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to equal", matcher.Expected)
}

func (matcher *EqualMatcher) Matches(actual interface{}) (success bool) {
	v, _ := matcher.Match(actual)
	return v
}

func (matcher *EqualMatcher) String() string {
	return "Equal: todo"
}
