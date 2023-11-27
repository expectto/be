package be_jwt

import (
	"fmt"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
	"github.com/golang-jwt/jwt/v5"
)

func Token(args ...any) types.BeMatcher {
	return psi_matchers.NewNeverMatcher(fmt.Errorf("todo: not implemented"))
}

func HavingClaims(args ...any) types.BeMatcher {
	return psi_matchers.NewNeverMatcher(fmt.Errorf("todo: not implemented"))
}

func BeingValid() types.BeMatcher {

	// todo: sandbox
	// sample token string taken from the New example
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.9lhSXK6s77Pa5Ha9OxdzIp5whsVj07Yv28lDgQTpjJg"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	if err != nil {
		fmt.Println("ERRS ", err)
	}
	fmt.Println("...", token.Claims)

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		fmt.Println(err)
	}

	return psi_matchers.NewNeverMatcher(fmt.Errorf("todo: not implemented"))
}
