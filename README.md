<p align="center">
  <img src="logo.svg" alt="be" width="230">
</p>

<div align="center">

### <samp>Expect(tests).To(Be(readable)).</samp>

A large collection of composable test matchers for Go -<br>
works with stdlib `testing`, [Ginkgo](https://github.com/onsi/ginkgo)/[Gomega](https://github.com/onsi/gomega), [Gomock](https://github.com/uber-go/mock), and [testify](https://github.com/stretchr/testify).

[![Go Reference](https://pkg.go.dev/badge/github.com/expectto/be.svg)](https://pkg.go.dev/github.com/expectto/be)
[![Go Version](https://img.shields.io/github/go-mod/go-version/expectto/be)](go.mod)
[![License: MIT](https://img.shields.io/badge/license-MIT-yellow.svg)](LICENSE)

</div>

---

`be` is a matcher library: instead of asserting on booleans and losing the failure
message, you say what a value should *be*. Matchers compose - almost any argument
that takes a value also takes another matcher:

```go
be.Expect(t, user.Email).To(be_string.ValidEmail())
be.Expect(t, items).To(be.HaveLength(be.Gte(3)))
```

Every matcher works everywhere: with the built-in stdlib runner shown above, as a
Gomega matcher inside Ginkgo, and as a mock argument matcher for Gomock and
testify/mockery. The core module imports no test framework.

> [!NOTE]
> `be` is at `v1.0.0-rc.*` - the API is stable and being hardened for the v1.0.0 release.

## Contents

- [Install](#install)
- [Quick Start](#quick-start)
- [Say It with a Matcher](#say-it-with-a-matcher)
- [Matching an HTTP Request](#matching-an-http-request)
- [Test Framework Integration](#test-framework-integration)
- [Matchers](#matchers) - flat catalog in [MATCHERS.md](MATCHERS.md)
- [Contributing](#contributing)
- [License](#license)

## Install

```sh
go get github.com/expectto/be
```

## Quick Start

```go
package user_test

import (
	"testing"

	"github.com/expectto/be"
	"github.com/expectto/be/be_string"
)

func TestNewUser(t *testing.T) {
	u, err := NewUser("john@tests.com")
	be.Require(t, err).To(be.Succeed()) // hard fail, require-style

	be.Expect(t, u.Email).To(be_string.ValidEmail()) // soft fail, assert-style
	be.Expect(t, u.ID).To(be.NonZero())
}
```

No Gomega, no testify - just `*testing.T` and matchers. If you already use
Ginkgo/Gomega or testify, see [Test Framework Integration](#test-framework-integration).

## Say It with a Matcher

Wrapping a raw expression in `be.True(...)` throws away the failure message -
all you learn is "expected true, got false". There is a matcher for almost
every idiom; this table is the cheat-sheet (and what [`belint`](x/belint)
flags automatically):

| Instead of | Use |
|---|---|
| `be.Not(be.Nil())` | `be.NotNil()` |
| `be.HaveLength(0)` | `be.Empty()` |
| `be.Not(be.HaveLength(0))` | `be.NotEmpty()` |
| `be.Not(be.Eq(0))`, `be.Ne(0)` | `be.NonZero()` |
| `x == y → be.True()` | `be.Eq(y)` |
| `x >= n → be.True()` | `be.Gte(n)` |
| `len(xs) >= n → be.True()` | `be.HaveLength(be.Gte(n))` |
| `slices.Contains(xs, v) → be.True()` | `be.ContainElement(v)` |
| `strings.Contains(s, q) → be.True()` | `be.ContainSubstring(q)` |
| `strings.HasPrefix(s, p) → be.True()` | `be_string.HavingPrefix(p)` |
| `_, ok := m[k]; ok → be.True()` | `be.HaveKey(k)` |
| `errors.Is(err, X) → be.True()/False()` | `be.MatchError(X)` / `be.Not(be.MatchError(X))` |
| `errors.As(err, &v) → be.True()` | `be.MatchErrorAs[V]()` - *only when `v` is unused afterward* |
| `t1.Equal(t2) → be.True()` | `be_time.SameExactSecond(t2)` / `be_time.Approx(...)` |

The full flat catalog of every matcher across all packages lives in
[MATCHERS.md](MATCHERS.md).

## Matching an HTTP Request

Composability is the point: matchers nest into matchers, so one assertion can
describe an entire HTTP request - URL, method, context, JSON body, headers,
even a JWT inside a header template:

```go
req, err := buildRequestForServiceFoo()
Expect(err).To(Succeed())

Expect(req).To(be_http.Request(
    // Matching the URL
    be_http.HavingURL(be_url.URL(
        be_url.WithHttps(),
        be_url.HavingHost("example.com"),
        be_url.HavingPath("/path"),
        be_url.HavingSearchParam("status", "active"),
        be_url.HavingSearchParam("v", be_reflected.AsNumericString()),
        be_url.HavingSearchParam("q", "Hello World"),
    )),

    // Matching the HTTP method
    be_http.POST(),

    // Matching the request's context
    be_http.HavingCtx(
        be_ctx.CtxWithDeadline(be_time.LaterThan(time.Now().Add(30*time.Minute))),
        be_ctx.CtxWithValue("foobar", 100),
    ),

    // Matching the request body using JSON matchers
    be_http.HavingBody(
        be.JSON(
            be_json.JsonAsReader,
            be_json.HaveKeyValue("hello", "world"),
            // NOTE: JSON numbers decode to float64, so use AsFloat (not AsInteger) here
            be_json.HaveKeyValue("n", be_reflected.AsFloat(), be_math.GreaterThan(10)),
            be_json.HaveKeyValue("ids", be_reflected.AsSliceOf[string]()),
            Not(be_json.HaveKeyValue("deleted_field")),

            be_json.HaveKeyValue("email", be_string.ValidEmail(), HaveSuffix("@tests.com")),

            // "details":[{"key":"foo"},{"key":"bar"}]
            be_json.HaveKeyValue("details", And(
                be_reflected.AsObjects(),
                be.HaveLength(be_math.GreaterThan(2)),
                ContainElements(
                    be_json.HaveKeyValue("key", "foo"),
                    be_json.HaveKeyValue("key", "bar"),
                ),
            )),
        ),
    ),

    // Matching HTTP headers
    be_http.HavingHeader("X-Custom", "Hey-There"),
    be_http.HavingHeader("Authorization",
        be_string.MatchTemplate("Bearer {{jwt}}",
            be_string.V("jwt",
                be_jwt.Token(
                    be_jwt.Valid(),
                    be_jwt.HavingClaim("name", "John Doe"),
                ),
            ),
        ),
    ),
))
```

## Test Framework Integration

`be` matchers are framework-agnostic. The core `github.com/expectto/be` module
imports **no** test framework - pick how you run assertions:

**Standard library (no extra deps).** Two equivalent spellings - a fluent one and
a flat, testify-style one - both backed by the same engine, both driven by the
stdlib `*testing.T`:

```go
import "github.com/expectto/be"

// fluent (ginkgo/gomega-flavored)
be.Expect(t, n).To(be_math.GreaterThan(10))    // soft fail (assert-style)
be.Require(t, n).To(be_math.GreaterThan(10))   // hard fail (require-style)
be.Expect(t, s).NotTo(be_string.EmptyString())

// flat (testify-flavored)
be.AssertThat(t, n, be_math.GreaterThan(10))     // soft fail (assert-style)
be.RequireThat(t, s, be_string.NonEmptyString()) // hard fail (require-style)
```

**Already on testify?** Keep your `assert`/`require` calls and reach for `be`
only where a matcher earns its keep - `be.AssertThat` / `be.RequireThat` are the
drop-in slots (no extra dependency):

```go
assert.Equal(t, want, got)          // testify, as usual
be.AssertThat(t, got, be.Eq(want))  // be - and now `got` can face any matcher,
                                    // e.g. be_url.URL(be_url.HavingHost("x"), ...)
```

> The **subject comes first** and the expected value lives inside the matcher
> (`be.Eq(want)`), so - unlike testify's `Equal(t, want, got)` - there's no
> want/got order to memorize or get wrong.

**Ginkgo / Gomega:** every `be` matcher already satisfies `gomega`'s matcher
interface, so use it directly inside `Expect(...).To(...)`.

### Mocking

`be` matchers also work as **mock argument matchers**:

**Gomock:** every `be` matcher already satisfies `gomock.Matcher`, so pass it directly:

```go
mockObj.EXPECT().Do(be_math.GreaterThan(10)).Return("ok")
```

**Testify mock / mockery:** wrap with `MatchedBy` (works for hand-written and
mockery-generated mocks) - the matcher equivalent of testify's own
`mock.MatchedBy`. This is the one place you need the separate `x/mock` module
(it's what keeps testify out of the core deps); install it with `@latest` (the
submodule shares version numbers with the core module, which confuses
`go get <pkg>@<version>`):

```sh
go get github.com/expectto/be/x/mock@latest
```

```go
import bemock "github.com/expectto/be/x/mock"

svc.On("Do", bemock.MatchedBy(be_math.GreaterThan(10))).Return("ok")
```

## Matchers

One package per domain; each has its own README with detailed docs. The flat
searchable catalog of everything is [MATCHERS.md](MATCHERS.md).

### Core Be

Core matchers for common testing scenarios. [Detailed docs](core-be-matchers.md)

- **Core:** `Always`, `Never`, `All`, `Any`, `Eq`, `Not`, `HaveLength`, `Dive`, `DiveAny`, `DiveFirst`
- **Everyday:** `Nil`, `NotNil`, `True`, `False`, `Eq`, `Ne`, `Zero`, `NonZero`, `Empty`, `NotEmpty`, `Identical`, `NotIdentical`, `Via`, `Succeed`, `HaveOccurred`, `MatchError`, `MatchErrorAs`, `Panic`, `NotPanic`, `ContainElement`, `ContainElements`, `ContainSubstring`, `HaveKey`, `HaveKeyWithValue`, `HaveField`, `HaveFields`
- **Numeric aliases at root** (from be_math): `Gt`, `Gte`, `Lt`, `Lte`, `GreaterThan`, `GreaterThanEqual`, `LessThan`, `LessThanEqual`, `InRange`, `Positive`, `Negative`
- **Assertion shortcuts & async:** `NoError`, `Error`, `ErrorIs` (hard, the testify `require` trio) · `Eventually`, `Consistently` (native poll loop, no gomega output leakage)

### be_reflected

Matchers on values' reflect kinds and types. [Detailed docs](be_reflected/README.md)

- **By reflect.Kind:** `AsKind`, `AsFunc`, `AsChan`, `AsPointer`, `AsFinalPointer`, `AsStruct`, `AsPointerToStruct`, `AsSlice`, `AsPointerToSlice`, `AsSliceOf`, `AsMap`, `AsPointerToMap`, `AsObject`, `AsObjects`, `AsPointerToObject`
- **Data types:** `AsString`, `AsBytes`, `AsNumeric`, `AsNumericString`, `AsInteger`, `AsIntegerString`, `AsFloat`, `AsFloatishString`
- **Interfaces:** `AsReader`, `AsStringer`
- **Type compatibility:** `AssignableTo`, `Implementing`

### be_math

Matchers for mathematical assertions. [Detailed docs](be_math/README.md)

- `GreaterThan`, `GreaterThanEqual`, `LessThan`, `LessThanEqual`, `Approx`, `InRange`, `Odd`, `Even`, `Negative`, `Positive`, `Zero`, `Integral`, `DivisibleBy`
- **Shortcuts:** `Gt`, `Gte`, `Lt`, `Lte`

### be_string

Matchers on strings. [Detailed docs](be_string/README.md)

- `NonEmptyString`, `EmptyString`, `Alpha`, `Numeric`, `AlphaNumeric`, `AlphaNumericWithDots`, `Float`, `Titled`, `LowerCaseOnly`, `MatchWildcard`, `ValidEmail`
- **Templates:** `MatchTemplate`

### be_time

Matchers on `time.Time`. [Detailed docs](be_time/README.md)

- `LaterThan`, `LaterThanEqual`, `EarlierThan`, `EarlierThanEqual`, `Eq`, `Approx`
- `SameExactMilli`, `SameExactSecond`, `SameExactMinute`, `SameExactHour`, `SameExactDay`, `SameExactWeekday`, `SameExactWeek`, `SameExactMonth`
- `SameSecond`, `SameMinute`, `SameHour`, `SameDay`, `SameYearDay`, `SameWeek`, `SameMonth`, `SameYear`, `SameTimezone`, `SameOffset`, `IsDST`

### be_jwt

Matchers on JSON Web Tokens (via [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt)). [Detailed docs](be_jwt/README.md)

- **Transformers:** `TransformSignedJwtFromString`, `TransformJwtFromString`
- **Matchers:** `Token`, `Valid`, `HavingClaims`, `HavingClaim`, `HavingMethodAlg`, `SignedVia`

### be_url

Matchers on `url.URL`. [Detailed docs](be_url/README.md)

- **Transformers:** `TransformUrlFromString`, `TransformSchemelessUrlFromString`
- **Matchers:** `URL`, `Values`, `HavingHost`, `HavingHostname`, `HavingScheme`, `NotHavingScheme`, `WithHttps`, `WithHttp`, `HavingPort`, `NotHavingPort`, `HavingPath`, `HavingRawQuery`, `HavingSearchParam`, `NotHavingSearchParam`, `HavingMultipleSearchParam`, `HavingUsername`, `HavingUserinfo`, `HavingPassword`

### be_ctx

Matchers on `context.Context`. [Detailed docs](be_ctx/README.md)

- `Ctx`, `CtxWithValue`, `CtxWithDeadline`, `CtxWithError`

### be_json

Matchers for expressive assertions on JSON. [Detailed docs](be_json/README.md)

- `Matcher`, `HaveKeyValue`

### be_struct

Matchers on struct fields. [Detailed docs](be_struct/README.md)

- `HavingField`

### be_http

Matchers on `http.Request`. [Detailed docs](be_http/README.md)

- `Request`, `HavingMethod`, `GET`, `HEAD`, `POST`, `PUT`, `PATCH`, `DELETE`, `OPTIONS`, `CONNECT`, `TRACE`
- `HavingURL`, `HavingBody`, `HavingHost`, `HavingProto`, `HavingCtx`, `HavingHeader`, `HavingHeaders`

## Contributing

Issues, ideas, and pull requests are welcome - see [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT - see [LICENSE](LICENSE).
