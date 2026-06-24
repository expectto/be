# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project aims to
follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

Work toward a stable **v1**: a framework-agnostic matcher core with opt-in drivers.

### Changed (rc.5)
- **`be_string.MatchTemplate` is now literal + anchored.** Non-placeholder text is
  treated literally (regexp.QuoteMeta), so punctuation common in real strings —
  SQL parens, `?`, `.`, `|` — matches verbatim instead of as regex; and templates
  match the whole string (`^...$`), not a substring. Fixes silent false-positives.
- **`be_struct.HavingField` accepts a matcher as the expected value** —
  `HavingField[T]("Age", be_math.GreaterThan(18))` now works (previously the value
  was only deep-equal compared). Also removed a side-effect that corrupted the
  matcher's failure message on a wrong type parameter.

### Added (rc.5)
- `be.Identical`/`NotIdentical` (pointer identity, like testify's Same/NotSame).
- `be.Via(transform, matcher)` — project the actual through an accessor before
  matching (e.g. assert a context value via a public getter when the key is
  unexported: `be.Via(GetActor, be.Eq(want))`).
- `be_url.Values(...)` — match a `url.Values` directly with the Having* matchers
  (no need to build a `*url.URL`), and `be_url.NotHavingSearchParam` (absent vs
  present-but-empty).

### Added (rc.4)
- More everyday matchers surfaced by dogfooding on amberpixels/r3: `Ne` (not
  equal), `Empty`/`NotEmpty`, `ContainSubstring` (in core `be`, alongside the
  collection `ContainElement`), and `be_reflected.AsNumeric` (matches any
  integer or float — handy for JSON, where numbers decode to float64).

### Changed (rc.3)
- Failure messages: collapse gomega's output to a single line only for short,
  scalar comparisons; large slice/struct/map mismatches keep gomega's multi-line,
  diff-friendly layout (fixes an `rc.1`/`rc.2` readability regression).
- `Dive` now also dives over **map values** (not just slices/arrays).
- The native runner accepts an optional per-assertion message:
  `be.Expect(t, x).To(matcher, "context")` (also on `NotTo`/`ToNot`).

### Added (rc.2)
- **Everyday matchers** surfaced by real-project dogfooding: `Nil`/`NotNil`
  (typed-nil aware), `True`/`False`, `Succeed`/`HaveOccurred`/`MatchError`,
  `Panic`/`NotPanic`, `ContainElement`/`ContainElements`, `HaveKey`/
  `HaveKeyWithValue`. These cover the nil/bool/error/panic/collection assertions
  that idiomatic Go unit tests reach for most.

### Added
- **Native assertion runner** — `be.Expect(t, x).To(...)` (soft) and
  `be.Require(t, x).To(...)` (hard), bound to a minimal `TestingT` interface that
  `*testing.T` satisfies. No ginkgo/testify/gomega import required by the user.
- **Testify driver** as a separate module `github.com/expectto/be/x/testify`
  (`Assert`/`Require`), keeping testify out of the core dependency graph.
- **Testify mock / mockery support** — `x/testify.Mock(matcher)` adapts a be
  matcher into a testify `mock.MatchedBy` argument matcher.
- **`be_http.HavingCtx`** — match a request's context via `be_ctx` matchers.
- Test coverage for previously-untested packages: `be_url`, `be_http`, `be_ctx`,
  `be_jwt`, `be_json` (144 specs).

### Changed
- **Semantic contract**: matchers now return `(false, error)` for input that
  cannot be evaluated (invalid JSON, undecodable/wrong-signed JWT, unparseable
  URL) instead of a silent non-match; a value that was evaluated but did not
  match still returns `(false, nil)`.
- Core module no longer depends on **testify** (moved to `x/testify`) or
  **gomock** (removed a stray interface assertion). gomega remains as the
  internal matching engine.
- `go` directive bumped to 1.26; all dependencies updated.
- `be_ctx.Ctx(args...)` now enforces the actual is a `context.Context` before
  applying sub-matchers; placeholder error strings replaced with real messages.

### Fixed
- Match-time panics on unexpected `actual` values in `Dive` (non-slice) and
  `be_math` (`Integral`/`DivisibleBy` on non-numeric); `Dive` also now supports
  arrays.
- `be_ctx.CtxWithError(nil)` now asserts the context carries no error instead of
  matching any context.
- `be_http.HavingBody` no longer panics on a nil-body request.
- Broken symbols in the README showcase example.

## [0.2.4] and earlier

See the git history for pre-v1 releases.
