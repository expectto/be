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

| Matcher                      | Description                                            | Example Usage                          |
|------------------------------|--------------------------------------------------------|----------------------------------------|
| `be.Always()`                | Always succeeds                                        | `be.Always()`                          |
| `be.Never(err)`              | Never succeeds and always fails with a specified error | `be.Never(errors.New("custom error"))` |
| `be.All(ms ...any)`          | Logical AND for multiple matchers                      | `be.All(matcher1, matcher2, ...)`      |
| `be.Any(ms ...any)`          | Logical OR for multiple matchers                       | `be.Any(matcher1, matcher2, ...)`      |
| `be.Eq(expected)`            | Checks for equality                                    | `be.Eq(expectedValue)`                 |
| `be.Not(matcher)`            | Negates the result of another matcher                  | `be.Not(anotherMatcher)`               |
| `be.HaveLength(args ...any)` | Matches the length of slices, arrays, strings, or maps | `be.HaveLength(lengthMatcher)`         |


# Matchers for HTTP Requests

Be provides powerful matchers specifically designed for testing HTTP requests, offering detailed and structured request validation. These matchers seamlessly integrate with both Gomega and Gomock, providing a flexible and expressive way to verify various aspects of your HTTP requests.

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
Be welcomes contributions! Feel free to open issues, suggest improvements, or submit pull requests. See CONTRIBUTING.md for more details.

# License
This project is licensed under the MIT License.