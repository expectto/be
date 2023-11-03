package psiMatchers

// AlwaysMatcher always matches
type AlwaysMatcher struct{}

func (m *AlwaysMatcher) Match(_ any) (bool, error)               { return true, nil }
func (m *AlwaysMatcher) FailureMessage(actual any) string        { return "" }
func (m *AlwaysMatcher) NegatedFailureMessage(actual any) string { return "" }
func (m *AlwaysMatcher) Matches(actual any) bool                 { return true }
func (m *AlwaysMatcher) String() string                          { return "" }
