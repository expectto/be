package be_jwt

import (
	"fmt"
	"github.com/expectto/be"
	"github.com/expectto/be/internal/cast"
	"github.com/expectto/be/internal/psi"
	"github.com/expectto/be/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/onsi/gomega/gcustom"
)

// Token TODO: transfer string => *jwt.Token, so other matchers must be token-level
func Token(args ...any) types.BeMatcher {
	return be.Always()
}

func BeingValidAndSignedWith(secretKey string) types.BeMatcher {
	return psi.Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		token, err := jwt.Parse(cast.AsString(actual), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
			}

			return []string{secretKey}, nil
		})
		if err != nil {
			return false, nil
		}
		return token.Valid, nil
	}))
}

func BeingValid() types.BeMatcher {
	return psi.Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		_, err := jwt.Parse(cast.AsString(actual), func(token *jwt.Token) (interface{}, error) {
			return []string{""}, nil
		})
		return err == nil, nil
	}))
}

func HavingClaims(args ...any) types.BeMatcher {
	return psi.Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		token, err := jwt.Parse(cast.AsString(actual), func(token *jwt.Token) (interface{}, error) {
			return []string{""}, nil
		})
		if err != nil {
			return false, nil
		}

		return psi.Psi(args...).Match(token.Claims)
	}))
}
