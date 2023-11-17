package be_url

import (
	"github.com/expectto/be/internal/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
	"net/url"
)

// todo: RawPath/EscapedPath matchers

// TransformUrlFromString returns string->*url.Url transform
var TransformUrlFromString = url.Parse

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

	if cast.IsString(args[0], cast.AllowCustomTypes(), cast.AllowPointers()) {
		if len(args) != 1 {
			panic("sting arg must be a single arg")
		}

		// match given string to whole url
		return psi_matchers.NewUrlFieldMatcher("Url", "", func(u *url.URL) any {
			return u.String()
		}, gomega.Equal(args[0]))
	}

	return Psi(args...)
}

func HavingHost(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingHost", "host",
		func(u *url.URL) any { return u.Host },
		args...,
	)
}

func HavingHostname(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingHostname", "hostname",
		func(u *url.URL) any { return u.Hostname() },
		args...,
	)
}

func HavingScheme(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingScheme", "scheme",
		func(u *url.URL) any { return u.Scheme },
		args...,
	)
}

func WithHttps() types.BeMatcher {
	return HavingScheme("https")
}
func WithHttp() types.BeMatcher {
	return HavingScheme("http")
}

func HavingPort(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingPort", "port",
		func(u *url.URL) any { return u.Port() },
		args...,
	)
}

func NotHavingPort(args ...any) types.BeMatcher {
	return Psi(gomega.Not(HavingPort(args...)))
}

func HavingPath(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingPath", "path",
		func(u *url.URL) any { return u.Path },
		args...,
	)
}

func HavingRawQuery(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingRawQuery", "rawQuery",
		func(u *url.URL) any { return u.RawQuery },
		args...,
	)
}

// todo: UrlHavingMultipleSearchParam -> the fact that param is found > 1 times
// todo: UrlHavingDistinctSearchParams -> meaning no search params are repeated

func HavingSearchParam(searchParamName string, args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingSearchParam", "searchParam",
		func(u *url.URL) any { return u.Query().Get(searchParamName) },
		args...,
	)
}

func HavingUsername(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingUsername", "username",
		func(u *url.URL) any { return u.User.Username() },
		args...,
	)
}

func HavingUserinfo(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingUserinfo", "userinfo",
		func(u *url.URL) any { return u.User.String() },
		args...,
	)
}

func HavingPassword(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"UrlHavingPassword", "password",
		func(u *url.URL) any { p, _ := u.User.Password(); return p },
		args...,
	)
}
