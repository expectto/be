package psi

import (
	"github.com/expectto/be/types"
)

func FromGomock(m types.GomockMatcher) types.BeMatcher {
	return &upgradedGomockMatcher{GomockMatcher: m}
}

// upgradedGomockMatcher wraps GomockMatcher and GomegaMatcher
// Upgrade "Gomock => Psi" is done via attaching methods of GomegaMatcher
type upgradedGomockMatcher struct {
	types.GomockMatcher

	gomegaMatchFunc                 func(any) (bool, error)
	gomegaFailureMessageFunc        func(any) string
	gomegaNegatedFailureMessageFunc func(any) string
}

func (cm *upgradedGomockMatcher) Match(x any) (bool, error) {
	return cm.Matches(x), nil
}
func (cm *upgradedGomockMatcher) FailureMessage(actual any) string {
	// todo Expected <>: {expected} to equal <> {received}
	return cm.String()
}
func (cm *upgradedGomockMatcher) NegatedFailureMessage(actual any) string {
	// todo Expected <>: {expected} not to equal <> {received}
	return "not " + cm.String()
}
