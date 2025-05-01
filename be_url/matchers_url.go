// Package be_url provides Be matchers on url.URL
package be_url

import (
	"net/url"

	"github.com/amberpixels/abu/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
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

// HavingHost succeeds if the actual value is a *url.URL and its Host matches the provided one (via direct value or matchers)
func HavingHost(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingHost", "host",
		func(u *url.URL) any { return u.Host },
		args...,
	)
}

// HavingHostname succeeds if the actual value is a *url.URL and its Hostname matches the provided one (via direct value or matchers)
func HavingHostname(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingHostname", "hostname",
		func(u *url.URL) any { return u.Hostname() },
		args...,
	)
}

// HavingScheme succeeds if the actual value is a *url.URL and its Scheme matches the provided one (via direct value or matchers)
func HavingScheme(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingScheme", "scheme",
		func(u *url.URL) any { return u.Scheme },
		args...,
	)
}

// NotHavingScheme succeeds if the actual value is a *url.URL and its Scheme negatively matches given value
// Example:  `Expect(u).To(NotHavingScheme())` matches url without a scheme
func NotHavingScheme(args ...any) types.BeMatcher {
	return Psi(gomega.Not(HavingScheme(args...)))
}

// WithHttps succeeds if the actual value is a *url.URL and its scheme is "https".
func WithHttps() types.BeMatcher {
	return HavingScheme("https")
}

// WithHttp succeeds if the actual value is a *url.URL and its scheme is "http".
func WithHttp() types.BeMatcher {
	return HavingScheme("http")
}

// HavingPort succeeds if the actual value is a *url.URL and its Port matches the provided one.
func HavingPort(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingPort", "port",
		func(u *url.URL) any { return u.Port() },
		args...,
	)
}

// NotHavingPort succeeds if the actual value is a *url.URL and its Port does not match the given one.
// Example:  `Expect(u).To(NotHavingPort())` matches port-less url
func NotHavingPort(args ...any) types.BeMatcher {
	return Psi(gomega.Not(HavingPort(args...)))
}

// HavingPath succeeds if the actual value is a *url.URL and its Path matches the given one.
func HavingPath(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingPath", "path",
		func(u *url.URL) any { return u.Path },
		args...,
	)
}

// HavingRawQuery succeeds if the actual value is a *url.URL and its RawQuery matches the given one.
func HavingRawQuery(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingRawQuery", "rawQuery",
		func(u *url.URL) any { return u.RawQuery },
		args...,
	)
}

// HavingSearchParam succeeds if the actual value is a *url.URL and
// its specified search parameter matches the provided arguments.
func HavingSearchParam(searchParamName string, args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingSearchParam", "searchParam",
		func(u *url.URL) any { return u.Query().Get(searchParamName) },
		args...,
	)
}

// HavingMultipleSearchParam succeeds if the actual value is a *url.URL and
// its specified search parameter (all its values via slice) matches the provided arguments.
func HavingMultipleSearchParam(searchParamName string, args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingMultipleSearchParam", "multipleSearchParam",
		func(u *url.URL) any { return u.Query()[searchParamName] },
		args...,
	)
}

// HavingUsername succeeds if the actual value is a *url.URL and its Username matches the provided one.
func HavingUsername(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingUsername", "username",
		func(u *url.URL) any { return u.User.Username() },
		args...,
	)
}

// HavingUserinfo succeeds if the actual value is a *url.URL and its User.String() matches the provided one.
func HavingUserinfo(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingUserinfo", "userinfo",
		func(u *url.URL) any { return u.User.String() },
		args...,
	)
}

// HavingPassword succeeds if the actual value is a *url.URL and its Password matches the provided one.
func HavingPassword(args ...any) types.BeMatcher {
	return psi_matchers.NewUrlFieldMatcher(
		"HavingPassword", "password",
		func(u *url.URL) any { p, _ := u.User.Password(); return p },
		args...,
	)
}

// todo: RawPath/EscapedPath matchers
// todo:"HavingDistinctSearchParam -> ensuring it has only a single search param and match it( fail if not)
//       Difference is that HavingSearchParam will not fail if given param is not single (and onl will match the first one)
