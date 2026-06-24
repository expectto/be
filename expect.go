package be

// expect.go provides a native, dependency-free assertion runner so `be` matchers
// can be used with the standard library `testing` package alone — no ginkgo,
// gomega or testify import required. It preserves the brand sentence:
//
//	be.Expect(t, actual).To(be.GreaterThan(3))   // "expect ... to be greater than 3"
//
// The matcher core stays gomega-backed internally; only the failure message is
// reshaped to read natively (see internal/beformat).

import (
	"github.com/expectto/be/internal/beformat"
	"github.com/expectto/be/internal/psi"
)

// TestingT is the minimal subset of *testing.T the native driver needs.
// *testing.T satisfies it; tests can supply a fake. Mirrors testify's approach
// so the runner never imports the heavyweight `testing` package contract.
type TestingT interface {
	Helper()
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}

// Expectation is a TestingT-bound assertion produced by Expect or Require.
type Expectation struct {
	t      TestingT
	actual any
	fatal  bool
}

// Expect begins a soft assertion: a failure is reported via Errorf and the test
// continues (assert-style).
func Expect(t TestingT, actual any) *Expectation {
	return &Expectation{t: t, actual: actual}
}

// Require begins a hard assertion: the first failure stops the test via Fatalf
// (require-style).
func Require(t TestingT, actual any) *Expectation {
	return &Expectation{t: t, actual: actual, fatal: true}
}

// To asserts that actual satisfies the matcher. The matcher may be a be/gomega/
// gomock matcher or a raw value (wrapped via Psi, like the rest of be). It
// returns true on success.
func (e *Expectation) To(matcher any) bool {
	e.t.Helper()
	m := psi.Psi(matcher)
	ok, err := m.Match(e.actual)
	if err != nil {
		return e.fail(err.Error())
	}
	if !ok {
		return e.fail(beformat.Compact(m.FailureMessage(e.actual)))
	}
	return true
}

// NotTo asserts that actual does NOT satisfy the matcher.
func (e *Expectation) NotTo(matcher any) bool {
	e.t.Helper()
	m := psi.Psi(matcher)
	ok, err := m.Match(e.actual)
	if err != nil {
		return e.fail(err.Error())
	}
	if ok {
		return e.fail(beformat.Compact(m.NegatedFailureMessage(e.actual)))
	}
	return true
}

// ToNot is an alias for NotTo.
func (e *Expectation) ToNot(matcher any) bool {
	e.t.Helper()
	return e.NotTo(matcher)
}

func (e *Expectation) fail(msg string) bool {
	e.t.Helper()
	if e.fatal {
		e.t.Fatalf("%s", msg)
	} else {
		e.t.Errorf("%s", msg)
	}
	return false
}
