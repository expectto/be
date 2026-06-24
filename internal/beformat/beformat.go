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

// maxCompactLen bounds how long a collapsed one-liner may be. Beyond it the
// original vertical layout is kept so large diffs stay readable.
const maxCompactLen = 120

// Compact strips gomega's "<type>:" tags and collapses its vertical layout onto a
// single line. For example:
//
//	Expected
//	    <int>: 3
//	to be >
//	    <int>: 5
//
// becomes "Expected 3 to be > 5".
//
// It only collapses scalar-ish messages: if the result would be long or contains
// composite values (maps/structs render with braces), the original multi-line,
// diff-friendly gomega formatting is preserved instead.
func Compact(msg string) string {
	oneLine := strings.TrimSpace(vertical.ReplaceAllString(typeTag.ReplaceAllString(msg, ""), " "))
	if len(oneLine) <= maxCompactLen && !strings.ContainsAny(oneLine, "{}") {
		return oneLine
	}
	return strings.TrimRight(msg, "\n")
}
