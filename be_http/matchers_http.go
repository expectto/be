// Package be_http provides matchers for url.Request
// TODO: more detailed documentation here is required
package be_http

import (
	"bytes"
	"io"
	"net/http"

	"github.com/expectto/be/be_json"
	"github.com/expectto/be/internal/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
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
		// ReqPropertyMatcher with empty args will simply check if `actual` is *http.Request
		return psi_matchers.NewReqPropertyMatcher("", "", nil)
	}

	if cast.IsString(args[0], cast.AllowCustomTypes(), cast.AllowPointers()) {
		if len(args) != 1 {
			panic("string arg must be a single arg")
		}

		// TODO: plan a feature, to improve the output of the failed part of the url
		//       This will be possible if instead of matching whole `req.URL.String()`
		//       we parse `req` into parts and construct combined matcher on them
		//       So then, e.g. failure message will strictly say: search-argument ?foo= is why we failed

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

// HavingMethod: Syntactic sugar
var (
	GET     = func() types.BeMatcher { return HavingMethod(http.MethodGet) }
	HEAD    = func() types.BeMatcher { return HavingMethod(http.MethodHead) }
	POST    = func() types.BeMatcher { return HavingMethod(http.MethodPost) }
	PUT     = func() types.BeMatcher { return HavingMethod(http.MethodPut) }
	PATCH   = func() types.BeMatcher { return HavingMethod(http.MethodPatch) }
	DELETE  = func() types.BeMatcher { return HavingMethod(http.MethodDelete) }
	OPTIONS = func() types.BeMatcher { return HavingMethod(http.MethodOptions) }
	CONNECT = func() types.BeMatcher { return HavingMethod(http.MethodConnect) }
	TRACE   = func() types.BeMatcher { return HavingMethod(http.MethodTrace) }
)

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
		func(req *http.Request) any {
			// TODO: do it in nicer form (Idea is to return a body but so it's still readable later)
			body, _ := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(body))

			return io.NopCloser(bytes.NewBuffer(body))
		},
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
// Key is a string key for a header, args can be nil or len(args)==1.
// Note: Golang's http.Header is `map[string][]string`, and matching is done on the FIRST value of the header
// in case if you have multiple-valued header that needs to be matched, use HavingHeaders() instead
//
// These are scenarios that can be handled here:
// (1) If no args are given, it simply matches a request with existed header by key.
// (2) If len(args) == 1 && args[0] is a stringish, it matches a request with header `Key: Args[0]`
// (3) if len(args) == 1 && args[0] is not stringish, it is considered to be matcher for header's value
// Examples:
// - HavingHeader("X-Header") matches request with non-empty X-Header header
// - HavingHeader("X-Header", "X-Value") matches request with X-Header: X-Value
// - HavingHeader("X-Header", HavePrefix("Bearer ")) matchers request with header(X-Header)'s value matching given HavePrefix matcher
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

	return psi_matchers.NewReqPropertyMatcher(
		"HavingHeader", "header",
		func(req *http.Request) any { return req.Header },
		be_json.HaveKeyValue(key, NewDiveMatcher(args[0], DiveModeFirst)),
	)
}

// HavingHeaders matches requests that have header with a given key.
// Key is a string key for a header, args can be nil or len(args)==1.
// Note: Matching is done on the list of header values.
// In case if you have single-valued header that needs to be matched, use HavingHeader() instead
//
// These are scenarios that can be handled here:
// (1) If no args are given, it simply matches a request with existed header by key.
// (2) If len(args) == 1 && args[0] is a stringish, it matches a request with header `Key: Args[0]`
// (3) if len(args) == 1 && args[0] is not stringish, it is considered to be matcher for header's value
// Examples:
// - HavingHeader("X-Header") matches request with non-empty X-Header header
// - HavingHeader("X-Header", "X-Value") matches request with X-Header: X-Value
// - HavingHeader("X-Header", Dive(HavePrefix("Foo "))) matchers request with multiple X-Header values, each of them having Foo prefix
func HavingHeaders(key string, args ...any) types.BeMatcher {
	if len(args) == 0 {
		// Behaves same way as HavingHeader(key)

		return psi_matchers.NewReqPropertyMatcher(
			"HavingHeaders", "header",
			func(req *http.Request) any { return req.Header },
			be_json.HaveKeyValue(key),
		)
	}
	if len(args) != 1 {
		panic("len(args) must be 0 or 1")
	}

	return psi_matchers.NewReqPropertyMatcher(
		"HavingHeader", "header",
		func(req *http.Request) any { return req.Header },
		be_json.HaveKeyValue(key, args[0]),
	)
}
