package types

// GomockMatcher represents a matcher compatible with the Gomock library.
// We intentionally refrain using the Matcher interface from the Gomock library
// to prevent direct dependency on the Gomock library.
// Interface source: https://github.com/uber-go/mock/blob/main/gomock/matchers.go
type GomockMatcher interface {
	Matches(x any) bool
	String() string
}

// GomegaMatcher represents a matcher compatible with the Gomega library.
// We intentionally refrain from using the GomegaMatcher interface from the Gomega library.
// Although the Gomega library is still utilized as a dependency in other parts of `be`,
// and is likely to remain as a dependency for the time being.
// Interface source: https://github.com/onsi/gomega/blob/master/types/types.go#L37
type GomegaMatcher interface {
	Match(actual interface{}) (success bool, err error)
	FailureMessage(actual interface{}) (message string)
	NegatedFailureMessage(actual interface{}) (message string)
}

// BeMatcher is the main matcher interface for the `be` library.
// It consolidates all types of matchers supported by `be`,
// which currently includes GomegaMatcher and GomockMatcher.
type BeMatcher interface {
	GomegaMatcher
	GomockMatcher
}
