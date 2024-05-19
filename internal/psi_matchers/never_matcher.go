package psi_matchers

import "github.com/expectto/be/types"

// NeverMatcher never matches (always fails)
type NeverMatcher struct {
	err error
}

var _ types.BeMatcher = &NeverMatcher{}

func NewNeverMatcher(err error) *NeverMatcher {
	return &NeverMatcher{err: err}
}

func (m *NeverMatcher) Match(_ any) (bool, error)               { return false, nil }
func (m *NeverMatcher) FailureMessage(actual any) string        { return m.err.Error() }
func (m *NeverMatcher) NegatedFailureMessage(actual any) string { return m.err.Error() /* todo */ }
func (m *NeverMatcher) Matches(actual any) bool                 { return false }
func (m *NeverMatcher) String() string                          { return m.err.Error() }
