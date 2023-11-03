package types

type GomegaMatcher interface {
	Match(actual interface{}) (success bool, err error)
	FailureMessage(actual interface{}) (message string)
	NegatedFailureMessage(actual interface{}) (message string)
}

type GomockMatcher interface {
	Matches(x any) bool
	String() string
}

// BeMatcher currently stands for matcher that fits both Gomega and Gomock libraries
// todo: that's a draft yet, and a subject to be changed:
//
//	we want to be more flexible: not to have gomock as a MUST probably
//	but let to have testify be suported as well
//	Though `gomega` will probably remain as a MUST dependency
type BeMatcher interface {
	GomegaMatcher
	GomockMatcher
}
