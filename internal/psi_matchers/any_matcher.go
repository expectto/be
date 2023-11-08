package psi_matchers

import (
	"fmt"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/types"

	"github.com/onsi/gomega/format"
)

// AnyMatcher is the psi upgrade for gomega's OrMatcher
type AnyMatcher struct {
	Matchers []types.BeMatcher

	// state
	firstSuccessfulMatcher types.BeMatcher
}

func NewAnyMatcher(ms ...any) *AnyMatcher {
	matchers := []types.BeMatcher{}
	for _, m := range ms {
		matchers = append(matchers, AsMatcher(m))
	}

	return &AnyMatcher{Matchers: matchers}
}

func (m *AnyMatcher) Match(actual any) (success bool, err error) {
	m.firstSuccessfulMatcher = nil
	for _, matcher := range m.Matchers {
		success, err := matcher.Match(actual)
		if err != nil {
			return false, err
		}
		if success {
			m.firstSuccessfulMatcher = matcher
			return true, nil
		}
	}
	return false, nil
}

func (m *AnyMatcher) FailureMessage(actual any) (message string) {
	// not the most beautiful list of matchers, but not bad either...
	return format.Message(actual, fmt.Sprintf("To satisfy at least one of these matchers: %s", m.Matchers))
}

func (m *AnyMatcher) NegatedFailureMessage(actual any) (message string) {
	return m.firstSuccessfulMatcher.NegatedFailureMessage(actual)
}

// todo: MatchMayChangeInTheFuture

func (m *AnyMatcher) Matches(actual any) bool {
	m.firstSuccessfulMatcher = nil
	for _, matcher := range m.Matchers {
		if !matcher.Matches(actual) {
			m.firstSuccessfulMatcher = matcher
			return false
		}
	}
	return true
}

func (m *AnyMatcher) String() string {
	return m.firstSuccessfulMatcher.String()
}
