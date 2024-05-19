package psi_matchers

import (
	"fmt"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/onsi/gomega/format"
	"strings"
)

type JwtTokenMatcher struct {
	publicName string // e.g. HaveClaims
	cb         func(t *jwt.Token) any
	matching   types.BeMatcher

	// todo: adjust gomock methods work as intended
	*MixinMatcherGomock
}

var _ types.BeMatcher = &JwtTokenMatcher{}

func NewJwtTokenMatcher(publicName string, cb func(token *jwt.Token) any, args ...any) *JwtTokenMatcher {
	matcher := &JwtTokenMatcher{
		publicName: publicName,
		cb:         cb,
	}

	matcher.MixinMatcherGomock = NewMixinMatcherGomock(matcher, "Token field of")

	// No args means that this matcher succeeds when actual url will have any non-empty {field value}
	if len(args) > 0 {
		// compressing the args as list of matchers
		// or falling back to Equal matcher in case if len(args)==1
		// see types.Psi() for more details
		matcher.matching = Psi(args...)
	}

	// todo: pass fieldName as description to gomega
	return matcher
}

func (matcher *JwtTokenMatcher) Match(actual any) (success bool, err error) {
	if actual == nil {
		return false, fmt.Errorf("%s() expects actual value not to be nil", "jwt.Match")
	}

	actualUrl, ok := actual.(*jwt.Token)
	if !ok {
		return false, nil
	}

	if matcher.cb == nil {
		// we're just matching a valid URL
		return true, nil
	}

	v := matcher.cb(actualUrl)

	// If no inner matchers were given, then we simply validated if {field value} is not empty
	if matcher.matching == nil {
		return v != "" && v != nil && v != 0, nil
	}

	// simply allow underlying matchers to do their job
	return matcher.matching.Match(v)
}

func (matcher *JwtTokenMatcher) FailureMessage(actual any) string {
	v := matcher.cb(actual.(*jwt.Token))

	if matcher.matching == nil {
		return format.Message(v, fmt.Sprintf(`to be a non-empty %s`, "foo"))
	}
	return matcher.matching.FailureMessage(v)
}

func (matcher *JwtTokenMatcher) NegatedFailureMessage(actual any) string {
	// todo: not so accurate
	return strings.Replace(matcher.FailureMessage(actual), "\nto ", "\nnot to ", 1)
}
