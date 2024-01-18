# ExpectTo/Be: Versatile Golang Matcher Library Designed for Testing with [Ginkgo](https://github.com/onsi/ginkgo) and [Gomock](https://github.com/uber-go/mock).

## Expect(👨🏼‍💻).To(Be(🚀))

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/expectto/be/blob/main/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/expectto/be.svg)](https://pkg.go.dev/github.com/expectto/be)

`expectto/be` is a Golang package that offers a substantial collection of `Be` matchers. Every `Be` matcher is compatible with both [Ginkgo](https://github.com/onsi/ginkgo)/[Gomega](https://github.com/onsi/gomega)
and [Gomock](https://github.com/uber-go/mock). Where possible, arguments of matchers can be either finite values or matchers (Be/Gomega/Gomock).<br>
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
    be_http.HavingMethod(http.MethodPost),
    
    // Matching request's context
    be_http.HavingCtx(be_ctx.Ctx(
        be_ctx.WithDeadline(be_time.LaterThan(time.Now().Add(30*time.Minute))),
        be_ctx.WithValue("foobar", 100),
    ))

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
            be_strings.Template("Bearer {{jwt}}",
                be_strings.MatchingPart("jwt",
                    be_jwt.Token(
                        be_jwt.BeingValid(),
                        be_jwt.HavingClaims("name", "John Doe"),
                    ),
                ),
            ),
        ),
    ),
))
```

## Matchers

### (core) be
<details>
  <summary>Be provides a set of core matchers for common testing scenarios:</summary>

| Matcher                      | Example Usage                                                 | Description                                                                           |
|------------------------------|---------------------------------------------------------------|---------------------------------------------------------------------------------------|
| `be.Always()`                | `Expect(anything()).To(be.Always())`                          | Always succeeds (passes).                                                             |
| `be.Never(err)`              | `Expect(anything()).To(be.Never(errors.New("custom error")))` | Never succeeds and always fails with a specified error                                |
| `be.All(ms ...any)`          | `Expect(m).To(be.All(HaveKey("foo"), HaveKey("bar"), ...))`   | Logical AND for multiple matchers. _Similar to Ginkgo's`And()`_                       |
| `be.Any(ms ...any)`          | `Expect(m).To(be.Any(HaveKey("foo"), HaveKey("bar"), ...)`    | Logical OR for multiple matchers. _Similar to Ginkgo's `Or()`_                        |
| `be.Eq(expected)`            | `Expect(v).To(be.Eq(expectedValue))`                          | Checks for equality. _Similar to Ginkgo's `Equal` _                                   |
| `be.Not(matcher)`            | `Expect(v).To(be.Not(anotherMatcher))`                        | Negates the result of another matcher. _Similar to Ginkgo's `Not()`_                  |
| `be.HaveLength(args ...any)` | `Expect(collection).To(be.HaveLength(lengthMatcher))`         | Matches the length of slices, arrays, strings, or maps. Supports matchers as argument |
</details>


### be_reflected

AssignableTo(), Implementing(), AsKind(), ...

### be_math

GreaterThan(), GreaterLessThan(), ...

### be_strings

EmptyString(), NonEmptyString(), Alpha(), ...

### be_time

LaterThan(), LaterThanEqual(), EarlierThan(), ...

### be_jwt

Token(), HavingClaims(), ...

### be_url

URL() HavingHost(), HavingHostname(), ...

### be_ctx

Ctx(), CtxWithValue(), CtxWithDeadline(), ...

### be_json

Matcher(), HaveKeyWithValue(), ...

### be_http

Request(), HavingMethod(), HavingUrl(), ...

# Contributing

`Be` welcomes contributions! Feel free to open issues, suggest improvements, or submit pull
requests. [Contribution guidelines for this project](CONTRIBUTING.md)

# License

This project is [licensed under the MIT License](LICENSE).