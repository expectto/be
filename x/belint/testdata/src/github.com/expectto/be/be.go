// Package be is a minimal stub of github.com/expectto/be for analysistest:
// belint only inspects call shapes and package paths, so signatures are all
// that matters here.
package be

type TestingT interface {
	Helper()
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}

type BeMatcher interface{ Match(actual any) (bool, error) }

type stub struct{}

func (stub) Match(any) (bool, error) { return true, nil }

type Expectation struct{}

func Expect(t TestingT, actual any) *Expectation  { return &Expectation{} }
func Require(t TestingT, actual any) *Expectation { return &Expectation{} }

func (e *Expectation) To(matcher any, msgAndArgs ...any) bool    { return true }
func (e *Expectation) NotTo(matcher any, msgAndArgs ...any) bool { return true }
func (e *Expectation) ToNot(matcher any, msgAndArgs ...any) bool { return true }

func AssertThat(t TestingT, actual, matcher any, msgAndArgs ...any) bool  { return true }
func RequireThat(t TestingT, actual, matcher any, msgAndArgs ...any) bool { return true }

func True() BeMatcher                     { return stub{} }
func False() BeMatcher                    { return stub{} }
func Nil() BeMatcher                      { return stub{} }
func NotNil() BeMatcher                   { return stub{} }
func Empty() BeMatcher                    { return stub{} }
func NotEmpty() BeMatcher                 { return stub{} }
func Zero() BeMatcher                     { return stub{} }
func NonZero() BeMatcher                  { return stub{} }
func Eq(expected any) BeMatcher           { return stub{} }
func Ne(expected any) BeMatcher           { return stub{} }
func Not(matcher any) BeMatcher           { return stub{} }
func HaveLength(args ...any) BeMatcher    { return stub{} }
func ContainElement(v any) BeMatcher      { return stub{} }
func ContainSubstring(s string) BeMatcher { return stub{} }
func MatchError(expected any) BeMatcher   { return stub{} }
func HaveKey(k any) BeMatcher             { return stub{} }
func Gt(v any) BeMatcher                  { return stub{} }
func Gte(v any) BeMatcher                 { return stub{} }
func Lt(v any) BeMatcher                  { return stub{} }
func Lte(v any) BeMatcher                 { return stub{} }
