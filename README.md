# Be: Fluent Matchers for Golang Testing

## Expect(👨🏼‍💻).To(Be(🚀))

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/expectto/be/blob/main/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/expectto/be.svg)](https://pkg.go.dev/github.com/expectto/be)

Be is a Golang package that introduces a comprehensive set of fluent matchers, fully compatible with both the Gomega and
Gomock libraries, allowing seamless integration into various testing scenarios. With Be, you can express your
expectations in a clear and concise manner, making your tests more readable and maintainable. Let Be empower your
testing suite with its expressive matchers tailored for Golang testing.

## Table of Contents

- [Installation](#installation)
- [Core Matchers](#core-matchers)
- [Matchers for HTTP Requests](#matchers-for-http-requests)
    - [Example](#example)
- [Contributing](#contributing)
- [License](#license)

## Installation

To use Be in your Golang project, simply import it:

```go
import "github.com/expectto/be"
```

## Core Matchers

Be provides a set of core matchers for common testing scenarios:

| Matcher                      | Example Usage                                                 | Description                                                                           |
|------------------------------|---------------------------------------------------------------|---------------------------------------------------------------------------------------|
| `be.Always()`                | `Expect(anything()).To(be.Always())`                          | Always succeeds (passes).                                                             |
| `be.Never(err)`              | `Expect(anything()).To(be.Never(errors.New("custom error")))` | Never succeeds and always fails with a specified error                                |
| `be.All(ms ...any)`          | `Expect(m).To(be.All(HaveKey("foo"), HaveKey("bar"), ...))`   | Logical AND for multiple matchers. _Similar to Ginkgo's`And()`_                       |
| `be.Any(ms ...any)`          | `Expect(m).To(be.Any(HaveKey("foo"), HaveKey("bar"), ...)`    | Logical OR for multiple matchers. _Similar to Ginkgo's `Or()`_                        |
| `be.Eq(expected)`            | `Expect(v).To(be.Eq(expectedValue))`                          | Checks for equality. _Similar to Ginkgo's `Equal` _                                   |
| `be.Not(matcher)`            | `Expect(v).To(be.Not(anotherMatcher))`                        | Negates the result of another matcher. _Similar to Ginkgo's `Not()`_                  |
| `be.HaveLength(args ...any)` | `Expect(collection).To(be.HaveLength(lengthMatcher))`         | Matches the length of slices, arrays, strings, or maps. Supports matchers as argument |

# Matchers for HTTP Requests

Be provides powerful matchers specifically designed for testing HTTP requests, offering detailed and structured request
validation. These matchers seamlessly integrate with both Gomega and Gomock, providing a flexible and expressive way to
verify various aspects of your HTTP requests.

## Example

Consider the following example demonstrating the usage of Be's HTTP request matchers:

```go
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
		be_http.HavingMethod("POST"),
		
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

# Contributing

`Be` welcomes contributions! Feel free to open issues, suggest improvements, or submit pull requests. [Contribution guidelines for this project](CONTRIBUTING.md)

# License

This project is [licensed under the MIT License](LICENSE).