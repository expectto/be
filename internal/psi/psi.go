// Package psi contains helpers that extends gomega library
// Name psi stands for previous letter from Omega
// (as we want to have a name that is close to gomega, but not to be a gomega)
//
// Package psi is considered as internal package to be used only inside `be`
// It's a core functionality that upgrades any matcher to be a `be` matcher
// Also it contains some core matchers and upgraded things from `gomega/gcustom` package
package psi

import (
	"github.com/expectto/be/types"
	"github.com/onsi/gomega/gcustom"
	"strings"
)

// Psi is a main converter function that converts given input into a PsiMatcher
func Psi(args ...any) types.BeMatcher {
	if len(args) == 0 {
		return &allMatcher{} // will always match
	}
	if len(args) == 1 {
		if IsTransformFunc(args[0]) {
			// not sure, add more tests here
			return WithFallibleTransform(args[0], nil)
		}
		// If a Match Func was given: simply wrap it in a matcher
		if IsMatchFunc(args[0]) {
			return AsMatcher(gcustom.MakeMatcher(AsMatchFunc(args[0])))
		}

		return AsMatcher(args[0])
	}

	// Special edge case:
	// args = {MatchFunc, message}
	// OR
	// args = {Matcher, message}
	if len(args) == 2 {
		message, ok := args[1].(string)
		if ok {
			if IsMatchFunc(args[0]) {
				return AsMatcher(gcustom.MakeMatcher(AsMatchFunc(args[0]), message))
			} else if IsMatcher(args[0]) {
				return AsMatcher(gcustom.MakeMatcher(AsMatcher(args[0]).Match, message))
			}
		}
	}

	matchers := make([]types.BeMatcher, 0)

	// Cast each arg as:
	// 1. transform func: will be wrapped via WithFallibleTransform then
	// 2. Matcher (Gomega/Gomock/Psi)
	// 3. any raw value will be converted to EqualMatcher // TODO: this case must be eliminated. We should be stricter
	for i, arg := range args {
		if IsTransformFunc(arg) { // 1
			transformMatcher := WithFallibleTransform(arg, Psi(args[i+1:]...))
			matchers = append(matchers, Psi(transformMatcher))
			return &allMatcher{matchers: matchers}
		}

		matchers = append(matchers, Psi(arg)) // 2 or 3
	}

	return &allMatcher{matchers: matchers}
}

// allMatcher is declared here internally so we're not importing psi_matchers
// allMatcher matches if all given matchers were matched
// or when no matchers were given
type allMatcher struct {
	matchers []types.BeMatcher

	// state
	firstFailedMatcher types.BeMatcher
	allSucceedMatchers []types.BeMatcher
}

func (m *allMatcher) Match(actual any) (bool, error) {
	m.firstFailedMatcher = nil
	m.allSucceedMatchers = make([]types.BeMatcher, 0)
	for _, matcher := range m.matchers {
		success, err := matcher.Match(actual)
		if !success || err != nil {
			m.firstFailedMatcher = matcher
			return false, err
		} else {
			m.allSucceedMatchers = append(m.allSucceedMatchers, matcher)
		}
	}
	return true, nil
}

func (m *allMatcher) FailureMessage(actual any) (message string) {
	return m.firstFailedMatcher.FailureMessage(actual)
}

func (m *allMatcher) NegatedFailureMessage(actual any) (message string) {
	// todo: make it nicer
	messages := make([]string, len(m.allSucceedMatchers))
	for i, m := range m.allSucceedMatchers {
		messages[i] = m.NegatedFailureMessage(actual)
	}
	return strings.Join(messages, "\n and \n")
}

func (m *allMatcher) Matches(actual any) bool {
	m.firstFailedMatcher = nil
	for _, matcher := range m.matchers {
		if !matcher.Matches(actual) {
			m.firstFailedMatcher = matcher
			return false
		}
	}
	return true
}

func (m *allMatcher) String() string {
	return m.firstFailedMatcher.String()
}

// todo: allMatcher.MatchMayChangeInTheFuture
