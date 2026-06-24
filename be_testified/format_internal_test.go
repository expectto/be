package be_testified

import (
	"strings"
	"testing"

	"github.com/expectto/be/be_math"
	"github.com/expectto/be/be_string"
	"github.com/expectto/be/types"
)

// TestFailureMessageStripsGomegaDialect verifies that failure messages surfaced to
// testify are compact and free of gomega's vertical, type-tagged formatting.
func TestFailureMessageStripsGomegaDialect(t *testing.T) {
	for _, tc := range []struct {
		name    string
		actual  any
		matcher types.BeMatcher
	}{
		{"numeric", 3, be_math.GreaterThan(5)},
		{"string", "Hello World", be_string.LowerCaseOnly()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// gomega's contract: FailureMessage is only valid after Match.
			if ok, _ := tc.matcher.Match(tc.actual); ok {
				t.Fatalf("test setup: matcher unexpectedly matched %v", tc.actual)
			}
			msg := failureMessage(tc.actual, tc.matcher)
			if msg == "" {
				t.Fatal("failure message should not be empty")
			}
			if strings.Contains(msg, "\n") {
				t.Errorf("message should be a single line, got:\n%q", msg)
			}
			if gomegaTypeTag.MatchString(msg) {
				t.Errorf("message should not contain gomega <type>: tags, got: %q", msg)
			}
		})
	}
}
