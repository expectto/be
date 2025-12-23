// Package psi_matchers is a package that contains core matchers required
// Psi() to work properly
package psi_matchers

import (
	"fmt"

	. "github.com/expectto/be/internal/psi" //nolint:staticcheck // should be moved to lintignore
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

var _ types.BeMatcher = &AllMatcher{}

func NewAllMatcher(ms ...any) *AllMatcher {
	matchers := make([]types.BeMatcher, len(ms))
	for i, m := range ms {
		matchers[i] = AsMatcher(m)
	}

	return &AllMatcher{Matchers: matchers}
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

// todo: will be very nice if failure message will be slightly different
// depending on which one matcher inside AndGroup fails
