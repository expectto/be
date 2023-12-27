// Package be_jwt provides Be matchers for handling JSON Web Tokens (JWT).
// It includes matchers for transforming and validating JWT tokens.
// Matchers corresponds to specific golang jwt implementation: https://github.com/golang-jwt/jwt/v5
package be_jwt

import (
	"fmt"
	"github.com/expectto/be/internal/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/onsi/gomega"
	"strings"
)

// TransformSignedJwtFromString returns a transform function (string->*jwt.Token) for a given secret.
var TransformSignedJwtFromString = func(secret string) func(string) any {
	return func(input string) any {
		parsed, err := jwt.Parse(input, func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		if err != nil {
			return NewTransformError(fmt.Errorf("to parse jwt token (with secret=%s): %w", secret, err), input)
		}

		return parsed
	}
}

// TransformJwtFromString is a transform function (string->*jwt.Token) without a secret.
// It parses the input string as a JWT and returns the resulting *jwt.Token.
var TransformJwtFromString = func(input string) any {
	p := jwt.NewParser()

	t, parts, err := p.ParseUnverified(input, jwt.MapClaims{})
	if err != nil {
		return NewTransformError(err, input)
	}

	t.Signature, err = p.DecodeSegment(parts[2])
	if err != nil {
		return NewTransformError(fmt.Errorf("corrupted signature part: %w", err), input)
	}

	return t
}

// Token matches the actual value to be a valid *jwt.Token corresponding to given inputs.
// Possible inputs:
// 1. No args -> the actual value MUST be any valid *jwt.Token.
// 2. Single arg <string>. The actual value MUST be a *jwt.Token, whose .String() is compared against args[0].
// 3. Single arg <*jwt.Token>. The actual value MUST be a *jwt.Token.
// 4. List of Omega/Gomock/Psi matchers that are applied to *jwt.Token object.
//   - TransformJwtFromString/TransformSignedJwtFromString(secret) transforms can be given as the first argument,
//     so the string->*jwt.Token transform is applied.
func Token(args ...any) types.BeMatcher {
	if len(args) == 0 {
		return psi_matchers.NewJwtTokenMatcher("", nil)
	}

	// Given single string arg
	if cast.IsString(args[0], cast.AllowCustomTypes(), cast.AllowPointers()) {
		if len(args) != 1 {
			panic("string arg must be a single arg")
		}

		// match given string to raw token
		return psi_matchers.NewJwtTokenMatcher("Raw", func(t *jwt.Token) any {
			return t.Raw
		}, gomega.Equal(args[0]))
	}

	return Psi(args...)
}

// Valid checks if the jwt.Token is valid.
func Valid() types.BeMatcher {
	return psi_matchers.NewJwtTokenMatcher(
		"Valid",
		func(u *jwt.Token) any { return u.Valid },
		gomega.BeTrue(),
	)
}

// HavingClaims checks if the jwt.Token matches given claims.
func HavingClaims(args ...any) types.BeMatcher {
	return psi_matchers.NewJwtTokenMatcher(
		"Claims",
		func(u *jwt.Token) any { return u.Claims },
		Psi(args...),
	)
}

// HavingMethodAlg checks if the jwt.Token has a method and algorithm matching given arguments.
func HavingMethodAlg(args ...any) types.BeMatcher {
	return psi_matchers.NewJwtTokenMatcher(
		"Method.Alg()",
		func(u *jwt.Token) any { return u.Method.Alg() },
		Psi(args...),
	)
}

// SignedVia matches a valid & signed token (with a given secret).
// Token(TransformSignedJwtFromString(secret), Valid()) is the same as
// Token(TransformJwtFromString, SignedVia(secret)).
// It's useful when you already have matching against a secret-less token,
// and need the secret only for one specific matching.
func SignedVia(secret string) types.BeMatcher {
	return psi_matchers.NewJwtTokenMatcher(
		"Method.Verify()",
		func(u *jwt.Token) any {
			// text is parts[0]+parts[1]
			text := strings.Join(strings.Split(u.Raw, ".")[0:2], ".")
			return u.Method.Verify(text, u.Signature, []byte(secret))
		},
		gomega.BeNil(),
	)
}
