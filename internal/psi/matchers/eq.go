package psiMatchers

import (
	"bytes"
	"fmt"
	"github.com/onsi/gomega/format"
	"reflect"
)

type EqMatcher struct {
	Expected        any
	lastActualValue any
}

func (matcher *EqMatcher) Match(actual any) (success bool, err error) {
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
	matcher.lastActualValue = actual
	return reflect.DeepEqual(actual, matcher.Expected), nil
}

func (matcher *EqMatcher) FailureMessage(actual any) (message string) {
	actualString, actualOK := actual.(string)
	expectedString, expectedOK := matcher.Expected.(string)
	if actualOK && expectedOK {
		return format.MessageWithDiff(actualString, "to equal", expectedString)
	}

	return format.Message(actual, "to equal", matcher.Expected)
}

func (matcher *EqMatcher) NegatedFailureMessage(actual any) (message string) {
	return format.Message(actual, "not to equal", matcher.Expected)
}

func (matcher *EqMatcher) Matches(actual any) bool {
	res, _ := matcher.Match(actual)
	return res
}

// Todo: inaccurate behavior should be fixed
func (matcher *EqMatcher) String() string {
	mes := matcher.FailureMessage(matcher.lastActualValue)
	return mes
}
