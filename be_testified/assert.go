// Package be_testified lets you use all expectto/be matchers with the
// testify/assert and testify/require packages.
//
// The be matchers use gomega internally as their matching engine, but that is an
// implementation detail: this adapter never exposes gomega's API, and it reshapes
// gomega's vertical, type-tagged failure messages into compact, testify-native
// one-liners (see failureMessage).
package be_testified

import (
	"regexp"
	"strings"
	"testing"

	"github.com/expectto/be/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Assert uses a BeMatcher to assert that the actual value passes the matcher.
// It returns true if the match succeeds; otherwise, it fails the test using testify's assert.
func Assert(t *testing.T, actual any, matcher types.BeMatcher, msgAndArgs ...any) bool {
	t.Helper()
	success, err := matcher.Match(actual)
	if err != nil {
		return assert.Fail(t, err.Error(), msgAndArgs...)
	}
	if !success {
		return assert.Fail(t, failureMessage(actual, matcher), msgAndArgs...)
	}
	return true
}

// Require asserts that the actual value passes the matcher.
// It fails the test immediately using testify's require package if the match fails.
func Require(t *testing.T, actual any, matcher types.BeMatcher, msgAndArgs ...any) {
	t.Helper()
	success, err := matcher.Match(actual)
	if err != nil {
		require.FailNow(t, err.Error(), msgAndArgs...)
	}
	if !success {
		require.FailNow(t, failureMessage(actual, matcher), msgAndArgs...)
	}
}

var (
	// gomegaTypeTag matches gomega's "<type>: " object annotations, e.g. "<int>: ".
	gomegaTypeTag = regexp.MustCompile(`<[^>]*>:\s*`)
	// gomegaVertical matches a line break plus its surrounding indentation,
	// which gomega uses to lay failure messages out vertically.
	gomegaVertical = regexp.MustCompile(`[ \t]*\n[ \t]*`)
)

// failureMessage converts a be matcher's gomega-flavored FailureMessage into a
// compact, testify-friendly one-liner.
//
// gomega renders failures vertically with type tags, e.g.:
//
//	Expected
//	    <int>: 3
//	to be >
//	    <int>: 5
//
// which reads as foreign inside testify's report. We strip the "<type>:" tags and
// collapse the vertical layout onto a single line: "Expected 3 to be > 5".
// (Multi-line actual values such as large structs are collapsed too, which is the
// idiomatic testify trade-off in favour of compact messages.)
func failureMessage(actual any, matcher types.BeMatcher) string {
	msg := matcher.FailureMessage(actual)
	msg = gomegaTypeTag.ReplaceAllString(msg, "")
	msg = gomegaVertical.ReplaceAllString(msg, " ")
	return strings.TrimSpace(msg)
}
