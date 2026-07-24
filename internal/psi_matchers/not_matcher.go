package psi_matchers

import (
	. "github.com/expectto/be/internal/psi" //nolint:staticcheck // should be moved to lintignore
	"github.com/expectto/be/types"
)

type NotMatcher struct {
	Matcher         types.BeMatcher
	lastActualValue any
}

var _ types.BeMatcher = &NotMatcher{}

func NewNotMatcher(m any) *NotMatcher {
	return &NotMatcher{Matcher: AsMatcher(m)}
}

func (m *NotMatcher) Match(actual any) (bool, error) {
	success, err := m.Matcher.Match(actual)
	if err != nil {
		return false, err
	}
	m.lastActualValue = actual
	return !success, nil
}

func (m *NotMatcher) FailureMessage(actual any) string {
	return m.Matcher.NegatedFailureMessage(actual) // works beautifully
}

func (m *NotMatcher) NegatedFailureMessage(actual any) string {
	return m.Matcher.FailureMessage(actual) // works beautifully
}

func (m *NotMatcher) Matches(actual any) bool {
	res, _ := m.Match(actual)
	return res
}

// String returns the failure message for the last matched value.
// Todo: inaccurate behavior should be fixed.
func (m *NotMatcher) String() string {
	mes := m.FailureMessage(m.lastActualValue)
	return mes
}

// todo: MatchMayChangeInTheFuture
