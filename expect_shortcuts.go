package be

// expect_shortcuts.go provides flat shortcuts for the highest-frequency assertion:
// checking errors. They exist because error checks make up roughly a third of all
// assertions in a typical test suite, and typing the full matcher sentence for
// them every time pushes people (and LLMs) back to raw `if err != nil` blocks.
//
// All shortcuts are HARD assertions (Fatalf, require-style): an unexpected error
// state almost always invalidates the rest of the test. For a soft error check
// spell the sentence out:
//
//	be.AssertThat(t, err, be.Succeed())
//
// Shortcuts are only added for names that have no matcher twin at be. root.
// There is deliberately no be.Equal(t, a, b): it would sit next to the matcher
// be.Eq(x) and blur the line between shortcuts and matchers.

// NoError fails the test immediately (Fatalf) if err is non-nil. It is the
// drop-in for testify's require.NoError:
//
//	be.NoError(t, err)
//	be.NoError(t, err, "loading config %q", path)
//
// Equivalent to be.RequireThat(t, err, be.Succeed()). For a soft check use
// be.AssertThat(t, err, be.Succeed()).
func NoError(t TestingT, err error, msgAndArgs ...any) bool {
	t.Helper()
	return RequireThat(t, err, Succeed(), msgAndArgs...)
}

// Error fails the test immediately (Fatalf) if err is nil. It is the drop-in
// for testify's require.Error:
//
//	be.Error(t, err)
//
// Equivalent to be.RequireThat(t, err, be.HaveOccurred()). For a soft check use
// be.AssertThat(t, err, be.HaveOccurred()).
func Error(t TestingT, err error, msgAndArgs ...any) bool {
	t.Helper()
	return RequireThat(t, err, HaveOccurred(), msgAndArgs...)
}

// ErrorIs fails the test immediately (Fatalf) unless errors.Is(err, target). It
// is the drop-in for testify's require.ErrorIs:
//
//	be.ErrorIs(t, err, io.EOF)
//
// Equivalent to be.RequireThat(t, err, be.MatchError(target)). For a soft check
// use be.AssertThat(t, err, be.MatchError(target)).
func ErrorIs(t TestingT, err, target error, msgAndArgs ...any) bool {
	t.Helper()
	return RequireThat(t, err, MatchError(target), msgAndArgs...)
}
