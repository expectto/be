# Plan: DX + LLM discoverability

**Status:** implemented in v1.0.0-rc.7 (§1–§3; §4 happens in the r3 repo).
Originally verified against working tree at `eda23f4`.

**Context:** the amberpixels/r3 dogfooding migration (~770 assertions) showed the
main defect is **discoverability, not missing matchers**: when a matcher isn't
visible at `be.` root, humans and LLMs fall back to `be.True(<raw expr>)`.
~70% of the fix is docs/surface, ~30% new code.

---

## 1. New code

### 1.1 Error shortcuts at root (`expect_shortcuts.go`) — highest impact
Error checks were ~33% of all assertions. All **hard** (Fatalf); soft spelling
stays `be.AssertThat(t, err, be.Succeed())`.

```go
func NoError(t TestingT, err error, msgAndArgs ...any) bool  // RequireThat(t, err, Succeed(), ...)
func Error(t TestingT, err error, msgAndArgs ...any) bool    // RequireThat(t, err, HaveOccurred(), ...)
func ErrorIs(t TestingT, err, target error, msgAndArgs ...any) bool // RequireThat(t, err, MatchError(target), ...)
```

- No `be.Equal(t, a, b)` — it would sit next to matcher `be.Eq(x)` and create
  exactly the confusion we're fixing. Shortcuts only for names with no matcher twin.
- Tests: pass/fail table with fake `TestingT` (exists in `expect_test.go`);
  assert Fatalf + message passthrough.

### 1.2 `Zero` / `NonZero` at root
```go
func Zero() types.BeMatcher    { return Psi(gomega.BeZero()) }
func NonZero() types.BeMatcher { return Psi(gomega.Not(gomega.BeZero())) }
```
General zero-value (reflect-based, no cast). `be_math.Zero` (numeric) stays
namespaced — NOT aliased at root (name taken).

### 1.3 `HaveField` / `HaveFields` at root
```go
func HaveField(name string, matcher any) types.BeMatcher // wraps gomega.HaveField; "Field" or "Method()"
func HaveFields(fields map[string]any) types.BeMatcher   // AND of HaveField per entry
```
- Sort map keys → deterministic failure output (test this).
- Keep `be_struct.HavingField[T]` for compile-time field checks; document
  `HaveField` as default. Naming wobble (Have vs Having) is deliberate — note in
  both doc-comments.
- Tests: field mismatch names the field; `"Foo()"` method form; nesting.

### 1.4 `MatchErrorAs[T]` at root
```go
func MatchErrorAs[T error]() types.BeMatcher // errors.As(actual, new(T))
```
Custom `types.BeMatcher` (no gomega equivalent). Failure message names `T`.
Tests: wrapped typed error matches; unrelated error fails; nil fails.

### 1.5 `Eventually` (+ `Consistently`)
```go
func Eventually(t TestingT, actual any, matcher any, opts ...EventuallyOption) bool
// actual: value | func() T | func() (T, error)
```
- **Native poll loop** against `TestingT`, failure via same path as
  `Expectation.To` (so `beformat.Compact` applies). Do NOT wrap gomega's
  `Eventually` — it would leak gomega-formatted output.
- Options: `WithTimeout` (1s), `WithPolling` (10ms), `WithContext`.
- Tests: flips after N polls (counter, not wall-clock); timeout failure;
  func-returning-error.

### 1.6 `be_string.HavingPrefix` (+ suffix)
Verified absent. Add `HavingPrefix(p)` / `HavingSuffix(s)`.

### 1.7 Root numeric aliases (`matchers.go` pattern: `var X = be_math.X`)
`Gt, Gte, Lt, Lte, GreaterThan, GreaterThanEqual, LessThan, LessThanEqual,
InRange, Positive, Negative`.

**Rule:** alias at root only if collision-free AND hot. `be_time` never — `Eq`,
`Approx`, `Day`, `Month` collide; temporal matchers stay namespaced.

### 1.8 `HaveLength(matcher)` — already works, just untested/undocumented
`NewHaveLengthMatcher` accepts matchers; the test for it is **commented out** at
`internal/psi_matchers/have_length_matcher_test.go:44`. Restore that test, add a
`be.HaveLength(be.Gte(n))` test + doc example. No production change.

---

## 2. Docs (the primary LLM surface — agents read `go doc`, not the README)

### 2.1 Doc-comments
Every matcher gets: one-line semantics + copy-pasteable example + "prefer this
over `be.Not(be.X())`" where relevant. Priority (what the migration got wrong):
`Empty`, `NotEmpty`, `NotNil`, `Ne`, `Zero`, `ContainElement`, `HaveKey`,
`HaveLength` (composable form), `MatchError` (spell out tri-mode: errors.Is
target / string / matcher), `Eq` vs `Identical` (pointer-identity footgun).

### 2.2 Anti-pattern table (single source of truth — drives README, doc.go, catalog, belint)

| Instead of | Use |
|---|---|
| `be.Not(be.Nil())` | `be.NotNil()` |
| `be.HaveLength(0)` | `be.Empty()` |
| `be.Not(be.HaveLength(0))` | `be.NotEmpty()` |
| `be.Not(be.Eq(0))` | `be.Zero()` (or `be.Ne(x)`) |
| `slices.Contains(xs, v) → be.True()` | `be.ContainElement(v)` |
| `m[k] → be.True()` | `be.HaveKey(k)` |
| `x >= n → be.True()` | `be.Gte(n)` |
| `len(xs) >= n → be.True()` | `be.HaveLength(be.Gte(n))` |
| `errors.Is(err, X) → be.True()/False()` | `be.MatchError(X)` / `be.Not(be.MatchError(X))` |
| `errors.As(err, &v) → be.True()` | `be.MatchErrorAs[V]()` |
| `t1.Equal(t2) → be.True()` | `be_time.SameExactSecond(t2)` / `be_time.Approx(...)` |
| `strings.HasPrefix(s, p) → be.True()` | `be_string.HavingPrefix(p)` |

### 2.3 Root `doc.go`
`go doc github.com/expectto/be` shows the package comment first — it's the first
thing an agent sees. Contents: brand sentence + assertion spellings, condensed
anti-pattern table, pointer to MATCHERS.md and subpackages.

### 2.4 Flat `MATCHERS.md` catalog
One generated file over ALL packages, grouped by intent (equality, errors,
collections, numbers, time, strings, structs, async), with the anti-pattern
column. Extend `generate-docs.sh`; link from README top. (Per-package READMEs
don't work — you must already know `be_math` exists to look there.)

Skip `llms.txt` — hosted-docs convention, never discovered for a source-consumed
library.

### 2.5 CONTRIBUTING: matcher definition of done
Doc-comment with example + catalog entry + (if it supersedes a raw idiom) a
belint rule.

---

## 3. `belint` — analyzer flagging the anti-pattern table, with `-fix`

- Standard `analysis.Analyzer` under `x/belint` + **golangci-lint plugin
  registration from day one** (that's the lint entry point repos and agents
  actually run; a standalone vettool re-creates the discoverability problem).
- AST shape: the raw expression lives in the **actual** position —
  `be.AssertThat(t, x >= 0, be.True())` — match on that when the matcher-arg is
  `be.True()`/`be.False()`; also `be.Not(be.Nil())` etc. in matcher position.
- Point-of-use feedback is what makes LLM agents self-correct without reading docs.

---

## 4. Close the loop in r3

1. Drop `x/testify` (renamed to `x/mock` in `ecdeb9f`): `betestify.Assert/Require`
   → `be.AssertThat/RequireThat`.
2. `belint -fix` as mechanical first pass, then hand-review the ~77
   `be.True(rawExpr)` sites and ~250 error checks (→ 1.1 shortcuts).
3. Success metric: near-zero `be.True(<raw comparison/contains/len>)` and
   `be.Not(be.Nil()|be.HaveLength(0)|be.Eq(zero))` left.

---

## Sequencing

1. §2.1–2.3 doc-comments + anti-pattern table + doc.go (cheap, primary surface)
2. §1.1, 1.2, 1.7 (shortcuts, Zero, aliases — small, high-frequency)
3. §1.3, 1.4, 1.6, 1.8
4. §1.5 Eventually
5. §2.4 catalog generation
6. §3 belint → then §4 r3 pass
