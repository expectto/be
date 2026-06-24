// Package beformat renders matcher failure messages in a compact, framework-native
// form. be matchers use gomega internally, which formats failures vertically with
// "<type>:" tags; Compact reshapes that into a single line suitable for stdlib
// testing and testify output.
package beformat

import (
	"regexp"
	"strings"
)

var (
	// typeTag matches gomega's "<type>: " object annotations, e.g. "<int>: ".
	typeTag = regexp.MustCompile(`<[^>]*>:\s*`)
	// vertical matches a line break plus its surrounding indentation, which gomega
	// uses to lay failure messages out vertically.
	vertical = regexp.MustCompile(`[ \t]*\n[ \t]*`)
)

// Compact strips gomega's "<type>:" tags and collapses its vertical layout onto a
// single line. For example:
//
//	Expected
//	    <int>: 3
//	to be >
//	    <int>: 5
//
// becomes "Expected 3 to be > 5".
func Compact(msg string) string {
	msg = typeTag.ReplaceAllString(msg, "")
	msg = vertical.ReplaceAllString(msg, " ")
	return strings.TrimSpace(msg)
}
