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

// todo:

func Request(args ...any) types.BeMatcher {
	if len(args) == 0 {
		return psi_matchers.NewReqPropertyMatcher("", "", nil)
	}

	// todo: support custom string types
	if strArg, ok := args[0].(string); ok {
		if len(args) != 1 {
			panic("string arg must be a single arg")
		}

		// match given string to whole url
		return psi_matchers.NewReqPropertyMatcher("Url", "", func(req *http.Request) any {
			return req.URL.String()
		}, gomega.Equal(strArg))
	}

	return Psi(args...)
}

func RequestHavingMethod(args ...any) types.BeMatcher {
	return psi_matchers.NewReqPropertyMatcher(
		"RequestHavingMethod", "method",
		func(req *http.Request) any { return req.Method },
		args...,
	)
}

// todo syntax sugar for specific http methods

func RequestHavingURL(args ...any) types.BeMatcher {
	return psi_matchers.NewReqPropertyMatcher(
		"RequestHavingURL", "url",
		func(req *http.Request) any { return req.URL },
		args...,
	)
}

func RequestHavingHost(args ...any) types.BeMatcher {
	return psi_matchers.NewReqPropertyMatcher(
		"RequestHavingHost", "host",
		func(req *http.Request) any { return req.Host },
		args...,
	)
}

func RequestHavingHeader(args ...any) types.BeMatcher {
	// Syntax sugar: RequestHavingHeader("HeaderName)"
	//               or
	// 				 RequestHavingHeader("HeaderName", "HeaderValue")
	if len(args) == 2 && cast.IsStringish(args[0]) && cast.IsStringish(args[1]) {
		args = []any{
			be_json.HaveKeyValue(args[0].(string), []string{args[1].(string)}),
		}
	} else if len(args) == 1 && cast.IsStringish(args[0]) {
		args = []any{
			be_json.HaveKeyValue(args[0].(string)),
		}
	}

	return psi_matchers.NewReqPropertyMatcher(
		"RequestHavingHeader", "header",
		func(req *http.Request) any { return req.Header },
		args...,
	)
}
