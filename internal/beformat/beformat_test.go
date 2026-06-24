package beformat

import (
	"strings"
	"testing"
)

func TestCompactCollapsesScalars(t *testing.T) {
	in := "Expected\n    <int>: 3\nto be >\n    <int>: 5"
	if got := Compact(in); got != "Expected 3 to be > 5" {
		t.Errorf("scalar should collapse to one line, got %q", got)
	}
}

func TestCompactPreservesLargeMultiline(t *testing.T) {
	// A long slice mismatch must keep gomega's multi-line, diff-friendly layout
	// rather than collapsing into one unreadable line.
	in := "Expected\n    <[]int | len:30>: [" + strings.Repeat("1234567890, ", 20) + "]\nto equal\n    <[]int>: [...]"
	got := Compact(in)
	if !strings.Contains(got, "\n") {
		t.Errorf("large message should stay multi-line, got one line: %q", got)
	}
}

func TestCompactPreservesComposite(t *testing.T) {
	// Struct/map values render with braces and must not be flattened.
	in := "Expected\n    <main.T>: {Name: \"a\", Age: 1}\nto equal\n    <main.T>: {Name: \"b\", Age: 2}"
	if got := Compact(in); !strings.Contains(got, "\n") {
		t.Errorf("composite message should stay multi-line, got: %q", got)
	}
}
