package be_http

import (
	"github.com/expectto/be/be_json"
	"github.com/expectto/be/internal/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
	"net/http"
)

// Request matches an actual value to be a valid *http.Request corresponding to given inputs.
// Possible inputs:
// 1. Nil args -> so actual value MUST be any valid *http.Request.
// 2. Single arg <string>. Actual value MUST be a *http.Request, whose .URL.String() is compared against args[0].
// 3. List of Omega/Gomock/Psi matchers, that are applied to *http.Request object.
//   - Supports matching http.Request properties like method, URL, body, host, proto, and headers.
//   - Additional arguments can be used for matching specific headers, e.g., WithHeader("Content-Type", "application/json").
func Request(args ...any) types.BeMatcher {
	if len(args) == 0 {
		return psi_matchers.NewReqPropertyMatcher("", "", nil)
	}

	if cast.IsString(args[0], cast.AllowCustomTypes(), cast.AllowPointers()) {
		if len(args) != 1 {
			panic("string arg must be a single arg")
		}

		// match given string to whole url
		return psi_matchers.NewReqPropertyMatcher("Url", "", func(req *http.Request) any {
			return req.URL.String()
		}, gomega.Equal(cast.AsString(args[0])))
	}

	return Psi(args...)
}

// HavingMethod succeeds if the actual value is a *http.Request and its HTTP method matches the provided arguments.
func HavingMethod(args ...any) types.BeMatcher {
	return psi_matchers.NewReqPropertyMatcher(
		"HavingMethod", "method",
		func(req *http.Request) any { return req.Method },
		args...,
	)
}

// GET returns a matcher that succeeds if the actual *http.Request has a method "GET".
func GET() types.BeMatcher {
	return HavingMethod(http.MethodGet)
}

// POST returns a matcher that succeeds if the actual *http.Request has a method "POST".
func POST() types.BeMatcher {
	return HavingMethod(http.MethodPost)
}

// PUT returns a matcher that succeeds if the actual *http.Request has a method "PUT".
func PUT() types.BeMatcher {
	return HavingMethod(http.MethodPut)
}

// PATCH returns a matcher that succeeds if the actual *http.Request has a method "PATCH".
func PATCH() types.BeMatcher {
	return HavingMethod(http.MethodPatch)
}

// DELETE returns a matcher that succeeds if the actual *http.Request has a method "DELETE".
func DELETE() types.BeMatcher {
	return HavingMethod(http.MethodDelete)
}

// OPTIONS returns a matcher that succeeds if the actual *http.Request has a method "OPTIONS".
func OPTIONS() types.BeMatcher {
	return HavingMethod(http.MethodOptions)
}

// HavingURL succeeds if the actual value is a *http.Request and its URL matches the provided arguments.
func HavingURL(args ...any) types.BeMatcher {
	return psi_matchers.NewReqPropertyMatcher(
		"HavingURL", "url",
		func(req *http.Request) any { return req.URL },
		args...,
	)
}

// HavingBody succeeds if the actual value is a *http.Request and its body matches the provided arguments.
// Note: The body is not re-streamed, so it's not available after matching.
func HavingBody(args ...any) types.BeMatcher {
	return psi_matchers.NewReqPropertyMatcher(
		"HavingBody", "body",
		// todo: re-stream body so it's available after matching
		func(req *http.Request) any { return req.Body },
		args...,
	)
}

// HavingHost succeeds if the actual value is a *http.Request and its Host matches the provided arguments.
func HavingHost(args ...any) types.BeMatcher {
	return psi_matchers.NewReqPropertyMatcher(
		"HavingHost", "host",
		func(req *http.Request) any { return req.Host },
		args...,
	)
}

// HavingProto succeeds if the actual value is a *http.Request and its Proto matches the provided arguments.
func HavingProto(args ...any) types.BeMatcher {
	return psi_matchers.NewReqPropertyMatcher(
		"HavingProto", "proto",
		func(req *http.Request) any { return req.Proto },
		args...,
	)
}

// HavingHeader matches requests that have header with a given key.
// (1) If no args are given, it simply matches a request with existed header by key.
// (2) If len(args) == 1 && args[0] is a stringish, it matches a request with header `Key: Args[0]`
// (3) if len(args) == 1 && args[0] is not stringish, it is considered to be matcher for header's value
// Examples:
// - HavingHeader("X-Header") matches request with non-empty X-Header header
// - HavingHeader("X-Header", "X-Value") matches request with X-Header: X-Value
// - HavingHeader("X-Header", HavePrefix("Bearer ")) matchers request with header(X-Header)'s value matching given HavePrefix matcher
// -
// todo: support multiple header values
// todo: fixme I'm ugly for now
func HavingHeader(key string, args ...any) types.BeMatcher {
	if len(args) == 0 {
		return psi_matchers.NewReqPropertyMatcher(
			"HavingHeader", "header",
			func(req *http.Request) any { return req.Header },
			be_json.HaveKeyValue(key),
		)
	}
	if len(args) != 1 {
		panic("len(args) must be 0 or 1")
	}

	var headerValue []string
	if cast.IsStringish(args[0]) {
		headerValue = []string{cast.AsString(args[0])}
	} else if cast.IsStrings(args[0]) {
		headerValue = cast.AsStrings(args[0])
	}
	if headerValue != nil {
		return psi_matchers.NewReqPropertyMatcher(
			"HavingHeader", "header",
			func(req *http.Request) any { return req.Header },
			be_json.HaveKeyValue(key, headerValue),
		)
	}

	return psi_matchers.NewReqPropertyMatcher(
		"HavingHeader", "header[key]",
		func(req *http.Request) any { return req.Header[key][0] },
		args[0],
	)
}
