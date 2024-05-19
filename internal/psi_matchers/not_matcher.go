package psi_matchers

import (
	. "github.com/expectto/be/internal/psi"
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

func (m *NotMatcher) Match(actual interface{}) (bool, error) {
	success, err := m.Matcher.Match(actual)
	if err != nil {
		return false, err
	}
	m.lastActualValue = actual
	return !success, nil
}

func (m *NotMatcher) FailureMessage(actual interface{}) (message string) {
	return m.Matcher.NegatedFailureMessage(actual) // works beautifully
}

func (m *NotMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return m.Matcher.FailureMessage(actual) // works beautifully
}

func (m *NotMatcher) Matches(actual any) bool {
	res, _ := m.Match(actual)
	return res
}

// Todo: inaccurate behavior should be fixed
func (m *NotMatcher) String() string {
	mes := m.FailureMessage(m.lastActualValue)
	return mes
}

// todo: MatchMayChangeInTheFuture
