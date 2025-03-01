// Package be_testified is an experimental package that allows to use all expectto/be matchers
// with testify/assert and testify/require packages.
//
// When stabilized, no `ginkgo/gomega` will be imported if you only need be_testified.
package be_testified

import (
	"testing"

	"github.com/expectto/be/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Assert uses a BeMatcher to assert that the actual value passes the matcher.
// It returns true if the match succeeds; otherwise, it fails the test using testify's assert.
func Assert(t *testing.T, actual any, matcher types.BeMatcher, msgAndArgs ...interface{}) bool {
	success, err := matcher.Match(actual)
	if err != nil {
		return assert.Fail(t, err.Error(), msgAndArgs...)
	}
	if !success {
		failureMessage := matcher.FailureMessage(actual)
		return assert.Fail(t, failureMessage, msgAndArgs...)
	}
	return true
}

// Require asserts that the actual value passes the matcher.
// It fails the test immediately using testify's require package if the match fails.
func Require(t *testing.T, actual any, matcher types.BeMatcher, msgAndArgs ...interface{}) {
	success, err := matcher.Match(actual)
	if err != nil {
		require.FailNow(t, err.Error(), msgAndArgs...)
	}
	if !success {
		failureMessage := matcher.FailureMessage(actual)
		require.FailNow(t, failureMessage, msgAndArgs...)
	}
}
