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

func HavingMethod(args ...any) types.BeMatcher {
	return psi_matchers.NewReqPropertyMatcher(
		"HavingMethod", "method",
		func(req *http.Request) any { return req.Method },
		args...,
	)
}

func POST() types.BeMatcher {
	return psi_matchers.NewReqPropertyMatcher(
		"HavingMethod", "method",
		func(req *http.Request) any { return req.Method },
		"POST",
	)
}
func GET() types.BeMatcher {
	return psi_matchers.NewReqPropertyMatcher(
		"HavingMethod", "method",
		func(req *http.Request) any { return req.Method },
		"GET",
	)
}

// todo syntax sugar for specific http methods

func HavingURL(args ...any) types.BeMatcher {
	return psi_matchers.NewReqPropertyMatcher(
		"HavingURL", "url",
		func(req *http.Request) any { return req.URL },
		args...,
	)
}

func HavingBody(args ...any) types.BeMatcher {
	return psi_matchers.NewReqPropertyMatcher(
		"HavingBody", "body",
		// todo: re-stream body so it's available after matching
		func(req *http.Request) any { return req.Body },
		args...,
	)
}

func HavingHost(args ...any) types.BeMatcher {
	return psi_matchers.NewReqPropertyMatcher(
		"HavingHost", "host",
		func(req *http.Request) any { return req.Host },
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
