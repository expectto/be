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

func HavingHeader(args ...any) types.BeMatcher {
	// todo: handle better:
	// here we consider args[0] is header key, and args[1] is header value (single) or matcher for it
	// otherwise we fallback to matching req.Header that is map[string][]string
	// so value is OK to be string for our cases, but required to be []string when matching req.Header

	// Syntax sugar: RequestHavingHeader("HeaderName)"
	//               or
	// 				 RequestHavingHeader("HeaderName", "HeaderValue")
	// todo: fixme, doesn't work for more than 1 header
	if len(args) == 2 && cast.IsStringish(args[0]) && cast.IsStringish(args[1]) {
		args = []any{
			be_json.HaveKeyValue(cast.AsString(args[0]), []string{cast.AsString(args[1])}),
		}
	} else if len(args) == 1 && cast.IsStringish(cast.AsString(args[0])) {
		args = []any{
			be_json.HaveKeyValue(cast.AsString(args[0])),
		}
	}

	return psi_matchers.NewReqPropertyMatcher(
		"HavingHeader", "header",
		func(req *http.Request) any { return req.Header },
		args...,
	)
}
