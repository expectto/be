# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project aims to
follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

Work toward a stable **v1**: a framework-agnostic matcher core with opt-in drivers.

### Added (rc.7)
- **Error shortcuts at root** — `be.NoError(t, err)`, `be.Error(t, err)`,
  `be.ErrorIs(t, err, target)`: the testify `require` trio as hard (Fatalf)
  assertions. Soft spelling stays `be.AssertThat(t, err, be.Succeed())`.
- **`be.Eventually` / `be.Consistently`** — native async assertions: a be-owned
  poll loop against `TestingT` (not a gomega wrapper, so failures keep the
  compact native format). Actual may be a value, `func() T` or
  `func() (T, error)`; options: `be.WithTimeout`, `be.WithPolling`,
  `be.WithContext`.
- **`be.Zero` / `be.NonZero`** — general zero-value matchers (reflect-based,
  any type). The numeric-only `be_math.Zero` stays namespaced.
- **`be.HaveField` / `be.HaveFields`** — struct-field matching at root
  (wraps `gomega.HaveField`: dotted paths and `"Method()"` specs; the map form
  checks fields in sorted-key order for deterministic failures). The generic
  compile-time-typed `be_struct.HavingField[T]` remains for pinned struct types.
- **`be.MatchErrorAs[T]`** — the matcher spelling of `errors.As`, complementing
  `be.MatchError` (errors.Is / message / matcher).
- **Root numeric aliases** — `be.Gt/Gte/Lt/Lte`, `be.GreaterThan(Equal)`,
  `be.LessThan(Equal)`, `be.InRange`, `be.Positive`, `be.Negative` (aliases of
  `be_math`; rule: root alias only if collision-free and high-frequency —
  `be_time` is never aliased).
- **`be_string.HavingPrefix` / `HavingSuffix`** — matcher spellings of
  `strings.HasPrefix/HasSuffix`.
- **`x/belint`** — new module: a `go/analysis` linter (with `-fix` and
  golangci-lint module-plugin registration) that flags raw expressions wrapped
  in `be.True()/be.False()` and `be.Not(...)` compositions with a dedicated
  matcher (`be.Not(be.Nil())` → `be.NotNil()`, `be.True(len(xs) >= n)` →
  `be.HaveLength(be.Gte(n))`, ...).
- **`MATCHERS.md`** — generated flat catalog of every matcher across all
  packages, grouped by intent, with an "instead of" anti-pattern column
  (`internal/docgen`, wired into `generate-docs.sh`). Plus a root `doc.go`
  package comment (assertion spellings + anti-pattern table) and a
  doc-comment pass over the everyday matchers.

### Fixed (rc.7)
- `HaveLength(matcher)` negated failure message rendered the positive inner
  message ("length to be") instead of the negated one ("length not to be").
- Compact failure formatting swallowed a `<=`/`<` comparison operator standing
  right before a `<type>:` tag ("Expected 3 to be 2" instead of
  "Expected 3 to be <= 2").

### Added (rc.6)
- **`be.AssertThat` / `be.RequireThat`** — flat, testify-style spellings of
  `be.Expect(t, x).To(m)` / `be.Require(t, x).To(m)`, in the core module (no
  testify dependency). Lets testify users keep their `assert`/`require` calls and
  drop in a `be` matcher where it earns its keep:
  `be.AssertThat(t, got, be.Eq(want))`. The subject comes first and the expected
  value lives inside the matcher, so there's no want/got order to get wrong.

### Changed (rc.6)
- **`x/testify` module renamed to `x/mock`, scoped to mocks only.** Its
  `Assert`/`Require` are removed (superseded by the core `be.AssertThat` /
  `be.RequireThat`), and the mock adapter is renamed `Mock` → **`MatchedBy`** —
  the matcher equivalent of testify's `mock.MatchedBy`:
  `bemock.MatchedBy(be_math.GreaterThan(10))`. The module exists solely to keep
  testify out of the core dependency graph.

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
- **Testify mock / mockery support** — `bemock.MatchedBy(matcher)` (separate
  module `github.com/expectto/be/x/mock`) adapts a be matcher into a testify
  `mock.MatchedBy` argument matcher, keeping testify out of the core dependency
  graph. (Testify-style *assertions* are the core `be.AssertThat`/`be.RequireThat`
  above — no separate driver needed.)
- **`be_http.HavingCtx`** — match a request's context via `be_ctx` matchers.
- Test coverage for previously-untested packages: `be_url`, `be_http`, `be_ctx`,
  `be_jwt`, `be_json` (144 specs).

### Changed
- **Semantic contract**: matchers now return `(false, error)` for input that
  cannot be evaluated (invalid JSON, undecodable/wrong-signed JWT, unparseable
  URL) instead of a silent non-match; a value that was evaluated but did not
  match still returns `(false, nil)`.
- Core module no longer depends on **testify** (the mock adapter lives in the
  separate `x/mock` module) or **gomock** (removed a stray interface assertion).
  gomega remains as the internal matching engine.
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
