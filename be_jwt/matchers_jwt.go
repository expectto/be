package be_jwt

import (
	"fmt"
	"github.com/expectto/be/internal/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/onsi/gomega/gcustom"
)

var TransformJwtFromString = func(input string) (*jwt.Token, error) {
	return jwt.Parse(input, func(token *jwt.Token) (interface{}, error) {
		return []string{""}, nil
	})
}

var TransformSignedJwtFromString = func(secret string) func(string) (*jwt.Token, error) {
	return func(input string) (*jwt.Token, error) {
		return jwt.Parse(input, func(token *jwt.Token) (interface{}, error) {
			return []string{secret}, nil
		})
	}
}

// Token matches actual value to be a valid *jwt.Token corresponding to given inputs
// Possible inputs:
// 1. Nil args -> so actual value MUST be any valid *jwt.Token
// 2. Single arg <string>. Actual value MUST be a *jwt.Token, whose .String() compared against args[0]
// 3. Single arg <*jwt.Token>. Actual value MUST be a *jwt.Token
// 4. List of Omega/Gomock/Psi matchers, that are applied to *jwt.Token object
//   - TransformJwtFromString/TransformSignedJwtFromString(secret) transforms can be given as first argument,
//     so string->*jwt.Token transform is applied
func Token(args ...any) types.BeMatcher {
	if len(args) == 0 {
		return psi_matchers.NewUrlFieldMatcher("", "", nil)
	}

	if cast.IsString(args[0], cast.AllowCustomTypes(), cast.AllowPointers()) {
		if len(args) != 1 {
			panic("string arg must be a single arg")
		}

		return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
			token, ok := actual.(*jwt.Token)
			if !ok {
				return false, fmt.Errorf("actual must be a *jwt.Token")
			}
			return token.Raw == cast.AsString(args[0]), nil
		}))

		// match given string to whole url
		//return psi_matchers.NewUrlFieldMatcher("Url", "", func(u *url.URL) any {
		//	return u.String()
		//}, gomega.Equal(args[0]))
	}

	return Psi(args...)
}

func SignedVia(secret string) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		token, ok := actual.(*jwt.Token)
		if !ok {
			return false, fmt.Errorf("actual must be a *jwt.Token")
		}

		// actual token may be not signed, let's re-sign it with given signature

		//token.Raw

		return token.Valid, nil
	}))
}

func BeingValid() types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		token, ok := actual.(*jwt.Token)
		if !ok {
			return false, fmt.Errorf("actual must be a *jwt.Token")
		}

		return token.Valid, nil
	}))
}

func HavingClaims(args ...any) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		token, ok := actual.(*jwt.Token)
		if !ok {
			return false, fmt.Errorf("actual must be a *jwt.Token")
		}
		return Psi(args...).Match(token.Claims)
	}))
}
