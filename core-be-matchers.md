<p align="center">
  <img src="logo.svg" alt="be" width="144">
</p>

<div align="center">

Part of [`expectto/be`](README.md) - composable test matchers for Go.

</div>

---

```go
import "github.com/expectto/be"
```

Package be provides fluent, composable test matchers for Go with a native,
dependency-free assertion runner. Assertions read as an English sentence —
"expect ... to be ..." — and every matcher argument can be a raw value or
another matcher (Be/Gomega/Gomock), so matchers compose:

    be.Expect(t, resp.Items).To(be.HaveLength(be.Gte(3)))

# Assertion spellings

    be.Expect(t, actual).To(matcher)           // soft: Errorf, test continues
    be.Require(t, actual).To(matcher)          // hard: Fatalf, test stops
    be.AssertThat(t, actual, matcher)          // flat soft spelling
    be.RequireThat(t, actual, matcher)         // flat hard spelling
    be.NoError(t, err)                         // hard error shortcuts:
    be.Error(t, err)                           //   the testify require trio
    be.ErrorIs(t, err, target)
    be.Eventually(t, poll, matcher)            // async: poll until it matches

All matchers also work inside gomega (Expect(x).To(be.Eq(y))) and as gomock
argument matchers.

# Say it with a matcher, not with be.True

Wrapping a raw expression in be.True hides everything from the failure message.
There is a matcher for almost every idiom:

    be.True(x == y)                  -> be.Eq(y)
    be.True(x >= n)                  -> be.Gte(n)
    be.Not(be.Nil())                 -> be.NotNil()
    be.HaveLength(0)                 -> be.Empty()
    be.Not(be.HaveLength(0))         -> be.NotEmpty()
    be.Not(be.Eq(0)), be.Ne(0)       -> be.NonZero()
    be.True(len(xs) >= n)            -> be.HaveLength(be.Gte(n))
    be.True(slices.Contains(xs, v))  -> be.ContainElement(v)
    be.True(strings.Contains(s, q))  -> be.ContainSubstring(q)
    be.True(strings.HasPrefix(s, p)) -> be_string.HavingPrefix(p)
    _, ok := m[k]; be.True(ok)       -> be.HaveKey(k)
    be.True(errors.Is(err, target))  -> be.MatchError(target)
    var e E; be.True(errors.As(err, &e)) -> be.MatchErrorAs[E]() (if e unused after)
    be.True(t1.Equal(t2))            -> be_time.SameExactSecond(t2)

The full catalog of matchers across all packages lives in MATCHERS.md at the
repository root: https://github.com/expectto/be/blob/main/MATCHERS.md

# Subpackages

The root package covers everyday matchers (equality, nil, errors, booleans,
collections, lengths, structs) plus root aliases for hot numeric matchers.
Specialized matchers live in subpackages:

    - be_math: numbers (Approx, Odd, Even, DivisibleBy, ...)
    - be_string: strings (HavingPrefix, MatchTemplate, MatchWildcard, ...)
    - be_time: time.Time (SameExactSecond, Approx, LaterThan, ...)
    - be_struct: typed struct fields (HavingField[T])
    - be_reflected: kind/type assertions (AsNumericString, AsKind, ...)
    - be_http, be_url, be_json, be_jwt, be_ctx: HTTP requests, URLs, JSON,
      JWT tokens and contexts

Temporal matchers are never aliased at root (their Eq, Approx, Day would
collide) — always reach for be_time explicitly.

## Usage

```go
var Ctx = be_ctx.Ctx
```
Ctx is an alias for be_ctx.Ctx

```go
var GreaterThan = be_math.GreaterThan
```
GreaterThan is an alias for be_math.GreaterThan (long spelling of Gt).

```go
var GreaterThanEqual = be_math.GreaterThanEqual
```
GreaterThanEqual is an alias for be_math.GreaterThanEqual (long spelling of
Gte).

```go
var Gt = be_math.Gt
```
Gt is an alias for be_math.Gt: succeeds if actual is numerically > arg. Prefer
be.Gt(n) over be.True(x > n) — the failure message shows both values.

```go
var Gte = be_math.Gte
```
Gte is an alias for be_math.Gte: succeeds if actual is numerically >= arg.

```go
var HttpRequest = be_http.Request
```
HttpRequest is an alias for be_http.Request matcher

```go
var InRange = be_math.InRange
```
InRange is an alias for be_math.InRange: succeeds if actual is within [from,
until] with configurable inclusivity.

```go
var JSON = be_json.Matcher
```
JSON is an alias for be_json.JSON matcher

```go
var JwtToken = be_jwt.Token
```
JwtToken is an alias for be_jwt.Token matcher

```go
var LessThan = be_math.LessThan
```
LessThan is an alias for be_math.LessThan (long spelling of Lt).

```go
var LessThanEqual = be_math.LessThanEqual
```
LessThanEqual is an alias for be_math.LessThanEqual (long spelling of Lte).

```go
var Lt = be_math.Lt
```
Lt is an alias for be_math.Lt: succeeds if actual is numerically < arg.

```go
var Lte = be_math.Lte
```
Lte is an alias for be_math.Lte: succeeds if actual is numerically <= arg.

```go
var Negative = be_math.Negative
```
Negative is an alias for be_math.Negative: succeeds if actual is < 0.

```go
var Positive = be_math.Positive
```
Positive is an alias for be_math.Positive: succeeds if actual is > 0.

```go
var StringAsTemplate = be_string.MatchTemplate
```
StringAsTemplate is an alias for be_string.MatchTemplate matcher

```go
var URL = be_url.URL
```
URL is an alias for be_url.URL matcher

#### func  All

```go
func All(ms ...any) types.BeMatcher
```
All is like gomega.And()

#### func  Always

```go
func Always() types.BeMatcher
```
Always does always match

#### func  Any

```go
func Any(ms ...any) types.BeMatcher
```
Any is like gomega.Or()

#### func  AssertThat

```go
func AssertThat(t TestingT, actual, matcher any, msgAndArgs ...any) bool
```
AssertThat is the flat, testify-style spelling of Expect(t, actual).To(matcher):
a soft assertion that reports via Errorf and lets the test continue. It is the
drop-in for testify's assert when you want a be matcher:

    assert.Equal(t, want, got)          // testify
    be.AssertThat(t, got, be.Eq(want))  // be — and now `got` can face any matcher

The subject (actual) comes first and the expected value lives inside the
matcher, so unlike testify's Equal there is no want/got order to get wrong. An
optional message provides failure context (see To). Returns true on success.

#### func  Consistently

```go
func Consistently(t TestingT, actual any, matcher any, opts ...EventuallyOption) bool
```
Consistently polls actual and requires it to satisfy the matcher on EVERY poll
for the whole duration (default 100ms, set via WithTimeout). The first mismatch
fails the test (softly, via Errorf) immediately. actual takes the same forms as
in Eventually. Returns true if the matcher held throughout.

    be.Consistently(t, queue.Len, be.Zero())

#### func  ContainElement

```go
func ContainElement(element any) types.BeMatcher
```
ContainElement succeeds if actual (a slice, array or map) contains an element
that matches the given value or matcher:

    be.Expect(t, ids).To(be.ContainElement(42))
    be.Expect(t, users).To(be.ContainElement(be.HaveField("Name", "Alice")))

Prefer this over be.True(slices.Contains(xs, v)) — the failure message shows the
collection. For substrings of a string use ContainSubstring.

#### func  ContainElements

```go
func ContainElements(elements ...any) types.BeMatcher
```
ContainElements succeeds if actual contains all of the given elements (each may
be a value or a matcher), in any order.

#### func  ContainSubstring

```go
func ContainSubstring(substr string) types.BeMatcher
```
ContainSubstring succeeds if actual is a string containing the given substring.
(For slices/arrays/maps use ContainElement.)

#### func  Dive

```go
func Dive(matcher any) types.BeMatcher
```
Dive applies the given matcher to each (every) element of a slice or array, or
to each value of a map. Note: Dive is very close to gomega.HaveEach

#### func  DiveAny

```go
func DiveAny(matcher any) types.BeMatcher
```
DiveAny applies the given matcher to each element and succeeds in case if it
succeeds at least at one item

#### func  DiveFirst

```go
func DiveFirst(matcher any) types.BeMatcher
```
DiveFirst applies the given matcher to the first element of the given slice

#### func  DiveNth

```go
func DiveNth(n int, matcher any) types.BeMatcher
```
DiveNth applies the given matcher to the nth element of the given slice

#### func  Empty

```go
func Empty() types.BeMatcher
```
Empty succeeds if actual is empty: a zero-length string, slice, array, map or
channel (like gomega.BeEmpty):

    be.Expect(t, errsList).To(be.Empty())

Prefer this over be.HaveLength(0) or be.True(len(xs) == 0).

#### func  Eq

```go
func Eq(expected any) types.BeMatcher
```
Eq succeeds if actual equals expected by VALUE (deep equality, like
gomega.Equal):

    be.Expect(t, got).To(be.Eq(want))

Footgun to know: two different pointers to equal structs satisfy be.Eq. When you
mean "the same instance" (pointer identity, Go's ==), use be.Identical instead.
For "unset / zero value" prefer be.Zero over be.Eq(0).

#### func  Error

```go
func Error(t TestingT, err error, msgAndArgs ...any) bool
```
Error fails the test immediately (Fatalf) if err is nil. It is the drop-in for
testify's require.Error:

    be.Error(t, err)

Equivalent to be.RequireThat(t, err, be.HaveOccurred()). For a soft check use
be.AssertThat(t, err, be.HaveOccurred()).

#### func  ErrorIs

```go
func ErrorIs(t TestingT, err, target error, msgAndArgs ...any) bool
```
ErrorIs fails the test immediately (Fatalf) unless errors.Is(err, target). It is
the drop-in for testify's require.ErrorIs:

    be.ErrorIs(t, err, io.EOF)

Equivalent to be.RequireThat(t, err, be.MatchError(target)). For a soft check
use be.AssertThat(t, err, be.MatchError(target)).

#### func  Eventually

```go
func Eventually(t TestingT, actual any, matcher any, opts ...EventuallyOption) bool
```
Eventually polls actual until it satisfies the matcher, failing the test
(softly, via Errorf) if it never does within the timeout. actual may be:

    - a plain value (matched repeatedly — useful for stateful matchers),
    - func() T — polled each interval,
    - func() (T, error) — a returned error means "not ready yet"; polling continues.

Example:

    be.Eventually(t, queue.Len, be.Gte(3))
    be.Eventually(t, fetchStatus, be.Eq("ready"), be.WithTimeout(5*time.Second))

The failure message reports the last mismatch in the same compact format as
be.Expect. Returns true on success.

#### func  False

```go
func False() types.BeMatcher
```
False succeeds if actual is the boolean false.

#### func  HaveField

```go
func HaveField(field string, value any) types.BeMatcher
```
HaveField succeeds if actual is a struct (or pointer to one) whose field — or
nil-safe method chain — matches the given value or matcher. The field spec
follows gomega.HaveField: a name, a dotted path, or a "Method()" call:

    be.Expect(t, user).To(be.HaveField("Name", "Alice"))
    be.Expect(t, user).To(be.HaveField("Address.City", be.NotEmpty()))
    be.Expect(t, user).To(be.HaveField("ID()", be.NonZero()))

This is the default struct-field matcher. The naming wobble with
be_struct.HavingField is deliberate: HavingField[T] is the generic,
compile-time-checked variant for when you want the struct type enforced.

#### func  HaveFields

```go
func HaveFields(fields map[string]any) types.BeMatcher
```
HaveFields succeeds if actual matches HaveField for every entry of the given map
(logical AND). Values may be raw values or matchers:

    be.Expect(t, user).To(be.HaveFields(map[string]any{
    	"Name":  "Alice",
    	"Email": be_string.ValidEmail(),
    }))

Fields are checked in sorted-key order, so failure output is deterministic.

#### func  HaveKey

```go
func HaveKey(key any) types.BeMatcher
```
HaveKey succeeds if actual (a map) has a key matching the given value or
matcher:

    be.Expect(t, headers).To(be.HaveKey("Authorization"))

Prefer this over `_, ok := m[k]` followed by be.True(ok).

#### func  HaveKeyWithValue

```go
func HaveKeyWithValue(key, value any) types.BeMatcher
```
HaveKeyWithValue succeeds if actual (a map) has the given key with a matching
value.

#### func  HaveLength

```go
func HaveLength(args ...any) types.BeMatcher
```
HaveLength succeeds if the actual value (string, slice, array, map or channel)
has a length matching the provided condition — either an exact count, or one or
more matchers applied to the length (unlike gomega.HaveLen, which only takes a
count):

    be.Expect(t, items).To(be.HaveLength(3))
    be.Expect(t, items).To(be.HaveLength(be.Gte(3)))          // composable form
    be.Expect(t, name).To(be.HaveLength(be.InRange(1, true, 64, true)))

Prefer be.HaveLength(be.Gte(n)) over be.True(len(xs) >= n). For zero / non-zero
length prefer be.Empty / be.NotEmpty.

#### func  HaveOccurred

```go
func HaveOccurred() types.BeMatcher
```
HaveOccurred succeeds if actual is a non-nil error.

#### func  Identical

```go
func Identical(expected any) types.BeMatcher
```
Identical succeeds if actual is identical to expected using Go's == operator
(pointer identity for pointers). Like gomega.BeIdenticalTo / testify's Same:

    be.Expect(t, gotPtr).To(be.Identical(wantPtr)) // same pointer

Footgun to know: be.Eq compares by VALUE (deep equality) — two different
pointers to equal structs satisfy be.Eq but not be.Identical. Use Identical when
"the same instance" is what you mean.

#### func  MatchError

```go
func MatchError(expected any) types.BeMatcher
```
MatchError succeeds if actual is an error matching expected. It is tri-mode —
expected may be:

    - a target error, compared with errors.Is (wrapping-aware):
      be.Expect(t, err).To(be.MatchError(io.EOF))
    - a string, compared against err.Error():
      be.Expect(t, err).To(be.MatchError("file not found"))
    - a matcher, applied to err.Error():
      be.Expect(t, err).To(be.MatchError(be.ContainSubstring("not found")))

Prefer be.MatchError(target) over be.True(errors.Is(err, target)), and
be.Not(be.MatchError(target)) over be.False(errors.Is(err, target)). To match by
error TYPE (errors.As), use MatchErrorAs.

#### func  MatchErrorAs

```go
func MatchErrorAs[T error]() types.BeMatcher
```
MatchErrorAs succeeds if actual is an error that matches type T via errors.As —
the matcher spelling of `var target T; errors.As(err, &target)`:

    be.Expect(t, err).To(be.MatchErrorAs[*fs.PathError]())

Prefer this over projecting through errors.As into be.True() — but only when the
target goes unused afterward: unlike errors.As, the matcher does NOT bind the
concrete error value, so keep errors.As when you need target later. A nil or
non-matching error fails; a non-error actual is an error (not a mismatch). To
match by errors.Is target, message or matcher, use MatchError.

#### func  Ne

```go
func Ne(expected any) types.BeMatcher
```
Ne succeeds if actual is NOT equal to expected (the negation of Eq):

    be.Expect(t, status).To(be.Ne("failed"))

Prefer this over be.Not(be.Eq(x)). To assert "not the zero value" use
be.NonZero() instead of be.Ne(0).

#### func  Never

```go
func Never(err error) types.BeMatcher
```
Never does never succeed (does always fail)

#### func  Nil

```go
func Nil() types.BeMatcher
```
Nil succeeds if actual is nil. It is typed-nil aware (a nil *T inside an
interface matches), unlike a bare `== nil` comparison.

#### func  NoError

```go
func NoError(t TestingT, err error, msgAndArgs ...any) bool
```
NoError fails the test immediately (Fatalf) if err is non-nil. It is the drop-in
for testify's require.NoError:

    be.NoError(t, err)
    be.NoError(t, err, "loading config %q", path)

Equivalent to be.RequireThat(t, err, be.Succeed()). For a soft check use
be.AssertThat(t, err, be.Succeed()).

#### func  NonZero

```go
func NonZero() types.BeMatcher
```
NonZero succeeds if actual is NOT the zero value for its type:

    be.Expect(t, userID).To(be.NonZero())

Prefer this over be.Not(be.Eq(0)) or be.Ne(0).

#### func  Not

```go
func Not(expected any) types.BeMatcher
```
Not is like gomega.Not()

#### func  NotEmpty

```go
func NotEmpty() types.BeMatcher
```
NotEmpty succeeds if actual is not empty:

    be.Expect(t, results).To(be.NotEmpty())

Prefer this over be.Not(be.HaveLength(0)) or be.True(len(xs) > 0).

#### func  NotIdentical

```go
func NotIdentical(expected any) types.BeMatcher
```
NotIdentical succeeds if actual is NOT identical to expected (the negation of
Identical). Like testify's NotSame.

#### func  NotNil

```go
func NotNil() types.BeMatcher
```
NotNil succeeds if actual is not nil:

    be.Expect(t, user).To(be.NotNil())

Prefer this over be.Not(be.Nil()).

#### func  NotPanic

```go
func NotPanic() types.BeMatcher
```
NotPanic succeeds if actual is a func() that does not panic when invoked.

#### func  Panic

```go
func Panic() types.BeMatcher
```
Panic succeeds if actual is a func() that panics when invoked.

#### func  RequireThat

```go
func RequireThat(t TestingT, actual, matcher any, msgAndArgs ...any) bool
```
RequireThat is the flat, testify-style spelling of Require(t,
actual).To(matcher): a hard assertion that stops the test on the first failure
via Fatalf. See AssertThat for the argument-order rationale. Returns true on
success.

#### func  Succeed

```go
func Succeed() types.BeMatcher
```
Succeed succeeds if actual is a nil error. Intended for error values:

    be.Expect(t, err).To(be.Succeed())

#### func  True

```go
func True() types.BeMatcher
```
True succeeds if actual is the boolean true.

#### func  Via

```go
func Via(transform any, matcher any) types.BeMatcher
```
Via applies the transform function to the actual value and matches the result
against the given matcher. Handy for projecting through a public accessor when
the underlying value can't be matched directly, e.g.:

    be.Expect(t, ctx).To(be.Via(GetActor, be.Eq(wantActor)))

transform must be a function of one argument returning one value (and optionally
an error).

#### func  Zero

```go
func Zero() types.BeMatcher
```
Zero succeeds if actual is the zero value for its type: 0, "", nil, false, a
zero struct, etc. (reflect-based, works for any type — like gomega.BeZero):

    be.Expect(t, count).To(be.Zero())
    be.Expect(t, cfg).To(be.Zero()) // zero struct

Prefer this over be.Eq(0) when you mean "unset". For the numeric-only spelling
(where a non-number is an error, not a mismatch) use be_math.Zero.

#### type EventuallyOption

```go
type EventuallyOption func(*asyncConfig)
```

EventuallyOption configures Eventually and Consistently.

#### func  WithContext

```go
func WithContext(ctx context.Context) EventuallyOption
```
WithContext bounds the poll loop by a context: when the context is done the
assertion fails immediately instead of waiting for the timeout.

#### func  WithPolling

```go
func WithPolling(d time.Duration) EventuallyOption
```
WithPolling sets the interval between polls (default 10ms).

#### func  WithTimeout

```go
func WithTimeout(d time.Duration) EventuallyOption
```
WithTimeout sets how long Eventually keeps polling before failing (default 1s),
or how long Consistently keeps verifying before succeeding (default 100ms).

#### type Expectation

```go
type Expectation struct {
}
```

Expectation is a TestingT-bound assertion produced by Expect or Require.

#### func  Expect

```go
func Expect(t TestingT, actual any) *Expectation
```
Expect begins a soft assertion: a failure is reported via Errorf and the test
continues (assert-style).

#### func  Require

```go
func Require(t TestingT, actual any) *Expectation
```
Require begins a hard assertion: the first failure stops the test via Fatalf
(require-style).

#### func (*Expectation) NotTo

```go
func (e *Expectation) NotTo(matcher any, msgAndArgs ...any) bool
```
NotTo asserts that actual does NOT satisfy the matcher. An optional message
provides failure context (see To).

#### func (*Expectation) To

```go
func (e *Expectation) To(matcher any, msgAndArgs ...any) bool
```
To asserts that actual satisfies the matcher. The matcher may be a be/gomega/
gomock matcher or a raw value (wrapped via Psi, like the rest of be). An
optional message — a format string plus args, or plain values — is prepended to
the failure output for context. Returns true on success.

#### func (*Expectation) ToNot

```go
func (e *Expectation) ToNot(matcher any, msgAndArgs ...any) bool
```
ToNot is an alias for NotTo.

#### type TestingT

```go
type TestingT interface {
	Helper()
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}
```

TestingT is the minimal subset of *testing.T the native driver needs. *testing.T
satisfies it; tests can supply a fake. Mirrors testify's approach so the runner
never imports the heavyweight `testing` package contract.
