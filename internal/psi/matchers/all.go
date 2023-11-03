package psiMatchers

import (
	"fmt"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega/format"
)

// AllMatcher is the same as Gomega's AndMatcher
// with a difference that we can track errors if using via gomock
type AllMatcher struct {
	Matchers []types.BeMatcher

	// state
	firstFailedMatcher types.BeMatcher
}

func NewAllMatcher(ms ...types.BeMatcher) *AllMatcher {
	return &AllMatcher{Matchers: ms}
}

func (m *AllMatcher) Match(actual any) (success bool, err error) {
	m.firstFailedMatcher = nil
	for _, matcher := range m.Matchers {
		success, err := matcher.Match(actual)
		if !success || err != nil {
			m.firstFailedMatcher = matcher
			return false, err
		}
	}
	return true, nil
}

func (m *AllMatcher) FailureMessage(actual any) (message string) {
	return m.firstFailedMatcher.FailureMessage(actual)
}

func (m *AllMatcher) NegatedFailureMessage(actual any) (message string) {
	// not the most beautiful list of matchers, but not bad either...
	return format.Message(actual, fmt.Sprintf("To not satisfy all of these matchers: %s", m.Matchers))
}

func (m *AllMatcher) Matches(actual any) bool {
	m.firstFailedMatcher = nil
	for _, matcher := range m.Matchers {
		if !matcher.Matches(actual) {
			m.firstFailedMatcher = matcher
			return false
		}
	}
	return true
}

func (m *AllMatcher) String() string {
	return m.firstFailedMatcher.String()
}

// todo: AllMatcher.MatchMayChangeInTheFuture
