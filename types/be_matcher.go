package types

// GomockMatcher stands for matcher that fits Gomock library
// We intentionally don't use Matcher interface from gomock library
// To avoid dependency on gomock library
type GomockMatcher interface {
	Matches(x any) bool
	String() string
}

// GomegaMatcher stands for matcher that fits Gomega library
// We intentionally don't use GomegaMatcher interface from gomega library.
// Although Gomega lib is anyway used as a dependency in other parts of `be`
// And will probably remain as a dep for now
type GomegaMatcher interface {
	Match(actual interface{}) (success bool, err error)
	FailureMessage(actual interface{}) (message string)
	NegatedFailureMessage(actual interface{}) (message string)
}

// BeMatcher is main matcher interface for `be` library
// It combines all types of Matchers that `be` supports:
// Currently it's GomegaMatcher and GomockMatcher.
type BeMatcher interface {
	GomegaMatcher
	GomockMatcher
}
