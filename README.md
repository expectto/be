## Expect(👨🏼‍💻).To(Be(🚀))

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/expectto/be/blob/main/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/expectto/be.svg)](https://pkg.go.dev/github.com/expectto/be)

`expectto/be` is a Golang package that offers a substantial collection of `Be` matchers. Every `Be` matcher is
compatible with both [Ginkgo](https://github.com/onsi/ginkgo)/[Gomega](https://github.com/onsi/gomega)
and [Gomock](https://github.com/uber-go/mock). Where possible, arguments of matchers can be either finite values or
matchers (Be/Gomega/Gomock).<br>
Employing `expectto/be` matchers enables you to create straightforward, readable, and maintainable unit or
integration tests in Golang. Tasks such as testing HTTP requests, validating JSON responses, and more become remarkably
comprehensive and straightforward.

## Table of Contents

- [Installation](#installation)
- [Example](#example)
- [Matchers](#matchers)
    - [Be (core)](#core-be)
    - [Be Reflected](#be_reflected)
    - [Be Math](#be_math)
    - [Be Strings](#be_strings)
    - [Be Time](#be_time)
    - [Be JWT](#be_jwt)
    - [Be URL](#be_url)
    - [Be JSON](#be_json)
    - [Be HTTP](#be_http)

- [Contributing](#contributing)
- [License](#license)

## Installation

To use `Be` in your Golang project, simply import it:

```go
import "github.com/expectto/be"
```

## Example

Consider the following example demonstrating the usage of `expectto/be`'s HTTP request matchers:

```go
req, err := buildRequestForServiceFoo()
Expect(err).To(Succeed())

// Matching an HTTP request
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
    be_http.POST()

    // Matching request's context
    be_http.HavingCtx(be_ctx.Ctx(
        be_ctx.WithDeadline(be_time.LaterThan(time.Now().Add(30*time.Minute))),
        be_ctx.WithValue("foobar", 100),
    )),

    // Matching the request body using JSON matchers
    be_http.HavingBody(
        be_json.Matcher(
            be_json.JsonAsReader,
            be_json.HaveKeyValue("hello", "world"),
            be_json.HaveKeyValue("n", be_reflected.AsIntish()),
            be_json.HaveKeyValue("ids", be_reflected.AsSliceOf[string]),
            be_json.HaveKeyValue("details", And(
                be_reflected.AsObjects(),
                be.HaveLength(2),
                ContainElements(
                    be_json.HaveKeyValue("key", "foo"),
                    be_json.HaveKeyValue("key", "bar"),
                ),
            )),
        ),

        // Matching HTTP headers
        be_http.HavingHeader("X-Custom", "Hey-There"),
        be_http.HavingHeader("Authorization",
            be_strings.MatchTemplate("Bearer {{jwt}}",
                be_strings.Var("jwt",
                    be_jwt.Token(
                        be_jwt.Valid(),
                        be_jwt.HavingClaim("name", "John Doe"),
                    ),
                ),
            ),
        ),
    ),
))      
```

## Matchers

### Core Be

📦 `be` provides a set of core matchers for common testing scenarios.<br>[See detailed docs](core-be-matchers.md)

#### Core matchers:

`Always`, `Never`, `All`, `Any`, `Eq`, `Not`, `HaveLength`, `Dive`, `DiveAny`, `DiveFirst`

### be_reflected

📦 `be_reflected` provides Be matchers that use reflection, enabling expressive assertions on values' reflect kinds and
types.<br>[See detailed docs](be_reflected/README.md)

#### General Matchers based on reflect.Kind:

`AsKind`, `AsFunc`, `AsChan`, `AsPointer`, `AsFinalPointer`, `AsStruct`, `AsPointerToStruct`, `AsSlice`, `AsPointerToSlice`, `AsSliceOf`, `AsMap`, `AsPointerToMap`, `AsObject`, `AsObjects`, `AsPointerToObject`

#### Data Type Matchers based on reflect.Kind

`AsString`, `AsBytes`, `AsNumeric`, `AsNumericString`, `AsIntish`, `AsIntishString`, `AsFloatish`, `AsFloatishString`,

#### Interface Matchers based on reflect.Kind

`AsReader`,`AsStringer`

#### Matchers based on types compatibility:

`AssignableTo`, `Implementing`

### be_math

📦 `be_math` provides Be matchers for mathematical operations.<br>[See detailed docs](be_math/README.md)

#### Matchers on math:

`GreaterThan`, `GreaterThanEqual`, `LessThan`, `LessThanEqual`, `Approx`, `InRange`, `Odd`, `Even`, `Negative`, `Positive`, `Zero`, `Integral`, `DivisibleBy`

#### Shortcut aliases for math matchers:

`Gt`, `Gte`, `Lt`, `Lte`

### be_strings

📦 `be_strings` provides Be matchers for string-related assertions.<br>[See detailed docs](be_strings/README.md)

#### Matchers on strings

`NonEmptyString`, `EmptyString`, `Alpha`, `Numeric`, `AlphaNumeric`, `AlphaNumericWithDots`, `Float`, `Titled`, `LowerCaseOnly`, `MatchWildcard`, `ValidEmail`

#### Template matchers

`MatchTemplate`

### be_time

📦 `be_time` provides Be matchers on time.Time.<br>[See detailed docs](be_time/README.md)

#### Time Matchers

`LaterThan`, `LaterThanEqual`, `EarlierThan`, `EarlierThanEqual`, `Eq`, `Approx`, <br>
`SameExactMilli`, `SameExactSecond`, `SameExactMinute`, `SameExactHour`,  <br>
`SameExactDay`, `SameExactWeekday`, `SameExactWeek`, `SameExactMonth`, <br>
`SameSecond`, `SameMinute`, `SameHour`, `SameDay`, `SameYearDay`, <br>
`SameWeek`, `SameMonth`, `SameYear`, `SameTimzone`, `SameOffset`, `IsDST`

### be_jwt

📦 `be_jwt` provides Be matchers for handling JSON Web Tokens (JWT). It includes matchers for transforming and validating
JWT tokens. Matchers corresponds to specific
golang [jwt implementation](https://github.com/golang-jwt/jwt/v5).<br> [See detailed docs](be_jwt/README.md)

#### Transformers for JWT matching:

`TransformSignedJwtFromString`, `TransformJwtFromString`

#### Matchers on JWT:

`Token`, `Valid`, `HavingClaims`, `HavingClaim`, `HavingMethodAlg`, `SignedVia`

### be_url

📦 `be_url` provides Be matchers on url.URL.<br> [See detailed docs](be_jwt/README.md)

#### Transformers for URL Matchers:

`TransformUrlFromString`, `TransformSchemelessUrlFromString`

#### URL Matchers:

`URL`, `HavingHost`, `HavingHostname`, `HavingScheme`, `NotHavingScheme`, `WithHttps`, `WithHttp`, `HavingPort`, `NotHavingPort`, `HavingPath`, `HavingRawQuery`, `HavingSearchParam`, `HavingMultipleSearchParam`, `HavingUsername`, `HavingUserinfo`, `HavingPassword`

### be_ctx

📦 `be_ctx` provides Be matchers on context.Context.<br> [See detailed docs](be_ctx/README.md)

#### Context Matchers:

`Ctx`, `CtxWithValue`, `CtxWithDeadline`, `CtxWithError`

### be_json

📦 `be_json` provides Be matchers for expressive assertions on JSON.<br> [See detailed docs](be_json/README.md)

#### JSON Matchers:

`Matcher`, `HaveKeyValue`

### be_http

📦 `be_http` provides Be matchers for expressive assertions on http.Request.<br> [See detailed docs](be_http/README.md)

#### Matchers on HTTP:

`Request`, `HavingMethod`, <br>
`GET`, `HEAD`, `POST`, `PUT`, `PATCH`, `DELETE`, `OPTIONS`, `CONNECT`, `TRACE`, <br>
`HavingURL`, `HavingBody`, `HavingHost`, `HavingProto`, `HavingHeader`, `HavingHeaders`

# Contributing

`Be` welcomes contributions! Feel free to open issues, suggest improvements, or submit pull
requests. [Contribution guidelines for this project](CONTRIBUTING.md)

# License

This project is [licensed under the MIT License](LICENSE).