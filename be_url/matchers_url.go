package be_url

import (
	"github.com/expectto/be/internal/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
	"net/url"
)

// TransformUrlFromString returns string->*url.Url transform
var TransformUrlFromString = url.Parse

// TransformSchemelessUrlFromString returns string->*url.Url transform
// It allows string to be a scheme-less url
var TransformSchemelessUrlFromString = func(rawURL string) (*url.URL, error) {
	result, err := url.Parse(rawURL)
	if err == nil && result.Scheme == "" {
		result, err = url.Parse("http://" + rawURL)
		if err == nil {
			result.Scheme = ""
		}
	}
	return result, err
}

// URL matches actual value to be a valid URL corresponding to given inputs
// Possible inputs:
// 1. Nil args -> so actual value MUST be any valid *url.URL
// 2. Single arg <string>. Actual value MUST be a *url.URL, whose .String() compared against args[0]
// 3. Single arg <*url.Url>. Actual value MUST be a *url.URL, whose .String() compared against args[0].String()
// 4. List of Omega/Gomock/Psi matchers, that are applied to *url.URL object
//   - TransformUrlFromString() transform can be given as first argument, so string->*url.URL transform is applied
func URL(args ...any) types.BeMatcher {
	if len(args) == 0 {
		return psi_matchers.NewUrlFieldMatcher("", "", nil)
	}

	if cast.IsString(args[0], cast.AllowCustomTypes(), cast.AllowPointers()) {
		if len(args) != 1 {
			panic("string arg must be a single arg")
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
		"HavingHost", "host",
		func(u *url.URL) any { return u.Host },
		args...,
	)
}

func HavingHostname(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingHostname", "hostname",
		func(u *url.URL) any { return u.Hostname() },
		args...,
	)
}

func HavingScheme(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingScheme", "scheme",
		func(u *url.URL) any { return u.Scheme },
		args...,
	)
}

func NotHavingScheme(args ...any) types.BeMatcher {
	return Psi(gomega.Not(HavingScheme(args...)))
}

func WithHttps() types.BeMatcher {
	return HavingScheme("https")
}
func WithHttp() types.BeMatcher {
	return HavingScheme("http")
}

func HavingPort(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingPort", "port",
		func(u *url.URL) any { return u.Port() },
		args...,
	)
}

func NotHavingPort(args ...any) types.BeMatcher {
	return Psi(gomega.Not(HavingPort(args...)))
}

func HavingPath(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingPath", "path",
		func(u *url.URL) any { return u.Path },
		args...,
	)
}

// todo: RawPath/EscapedPath matchers

func HavingRawQuery(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingRawQuery", "rawQuery",
		func(u *url.URL) any { return u.RawQuery },
		args...,
	)
}

// todo:"HavingMultipleSearchParam -> the fact that param is found > 1 times
// todo:"HavingDistinctSearchParams -> meaning no search params are repeated

func HavingSearchParam(searchParamName string, args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingSearchParam", "searchParam",
		func(u *url.URL) any { return u.Query().Get(searchParamName) },
		args...,
	)
}

func HavingUsername(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingUsername", "username",
		func(u *url.URL) any { return u.User.Username() },
		args...,
	)
}

func HavingUserinfo(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingUserinfo", "userinfo",
		func(u *url.URL) any { return u.User.String() },
		args...,
	)
}

func HavingPassword(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingPassword", "password",
		func(u *url.URL) any { p, _ := u.User.Password(); return p },
		args...,
	)
}
