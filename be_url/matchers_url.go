package be_url

import (
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
	"net/url"
)

// todo: RawPath/EscapedPath matchers

// TransformUrlFromString returns string->*url.Url transform
func TransformUrlFromString() func(string) (*url.URL, error) {
	return url.Parse
}

// URL matches actual value to be a valid URL corresponding to given inputs
// Possible inputs:
// 1. Nil args -> so actual value MUST be any valid *url.URL
// 2. Single arg <string>. Actual value MUST be a *url.URL, whose .String() compared against args[0]
// 3. Single arg <*url.Url>. Actual value MUST be a *url.URL, whose .String() compared against args[0].String()
// 4. List of Omega/Gomock/Psi matchers, that are applied to *url.URL object
//   - AsStringUrl() transform can be given as first argument, so string->*url.URL transform is applied
func URL(args ...any) types.BeMatcher {
	if len(args) == 0 {
		return psi_matchers.NewUrlFieldMatcher("", "", nil)
	}

	// todo: support custom string types
	if strArg, ok := args[0].(string); ok {
		if len(args) != 1 {
			panic("sting arg must be a single arg")
		}

		// match given string to whole url
		return psi_matchers.NewUrlFieldMatcher("Url", "", func(u *url.URL) any {
			return u.String()
		}, gomega.Equal(strArg))
	}

	return Psi(args...)
}

func UrlHavingHost(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingHost", "host",
		func(u *url.URL) any { return u.Host },
		args...,
	)
}

func UrlHavingHostname(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingHostname", "hostname",
		func(u *url.URL) any { return u.Hostname() },
		args...,
	)
}

func UrlHavingScheme(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingScheme", "scheme",
		func(u *url.URL) any { return u.Scheme },
		args...,
	)
}

func UrlWithHttps() types.BeMatcher {
	return UrlHavingScheme("https")
}
func UrlWithHttp() types.BeMatcher {
	return UrlHavingScheme("http")
}

func UrlHavingPort(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingPort", "port",
		func(u *url.URL) any { return u.Port() },
		args...,
	)
}

func UrlNotHavingPort(args ...any) types.BeMatcher {
	return Psi(gomega.Not(UrlHavingPort(args...)))
}

func UrlHavingPath(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingPath", "path",
		func(u *url.URL) any { return u.Path },
		args...,
	)
}

func UrlHavingRawQuery(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingRawQuery", "rawQuery",
		func(u *url.URL) any { return u.RawQuery },
		args...,
	)
}

// todo: UrlHavingMultipleSearchParam -> the fact that param is found > 1 times
// todo: UrlHavingDistinctSearchParams -> meaning no search params are repeated

func UrlHavingSearchParam(searchParamName string, args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingSearchParam", "searchParam",
		func(u *url.URL) any { return u.Query().Get(searchParamName) },
		args...,
	)
}

func UrlHavingUsername(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingUsername", "username",
		func(u *url.URL) any { return u.User.Username() },
		args...,
	)
}

func UrlHavingUserinfo(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingUserinfo", "userinfo",
		func(u *url.URL) any { return u.User.String() },
		args...,
	)
}

func UrlHavingPassword(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingPassword", "password",
		func(u *url.URL) any { p, _ := u.User.Password(); return p },
		args...,
	)
}
