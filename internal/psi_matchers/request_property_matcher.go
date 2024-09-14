package psi_matchers

import (
	"fmt"
	"net/http"
	"strings"

	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega/format"
)

// ReqPropertyMatcher is a matcher for http.Request properties
type ReqPropertyMatcher struct {
	publicName string // e.g. RequestHavingMethod
	property   string // e.g. "method"

	// cb is a callback for extracting property value from http.Request
	cb func(r *http.Request) any

	// matching is a matcher for property value
	matching types.BeMatcher
}

// NewReqPropertyMatcher creates a new matcher for http.Request properties
func NewReqPropertyMatcher(publicName, fieldName string, cb func(r *http.Request) any, args ...any) types.BeMatcher {
	matcher := &ReqPropertyMatcher{publicName: publicName, property: fieldName, cb: cb}

	// No args means that this matcher succeeds when actual url will have any non-empty {field value}
	if len(args) > 0 {
		// compressing the args as list of matchers
		// or falling back to Equal matcher in case if len(args)==1
		// see types.Psi() for more details
		matcher.matching = Psi(args...)
	}

	// todo: pass fieldName as description to gomega
	return Psi(matcher)
}

func (matcher *ReqPropertyMatcher) Match(actual any) (success bool, err error) {
	if actual == nil {
		return false, fmt.Errorf("%s() expects actual value not to be nil", matcher.publicName)
	}

	actualReq, ok := actual.(*http.Request)
	if !ok {
		return false, fmt.Errorf("%s() expects actual value mast be a <*http.Request> received <%T>", matcher.publicName, actual)
	}

	if matcher.cb == nil {
		// we're just matching a valid URL
		return true, nil
	}

	v := matcher.cb(actualReq)

	// If no inner matchers were given, then we simply validated if {field value} is not empty
	if matcher.matching == nil {
		return v != "" && v != nil && v != 0, nil
	}

	// simply allow underlying matchers to do their job
	return matcher.matching.Match(v)
}

func (matcher *ReqPropertyMatcher) FailureMessage(actual any) string {
	v := matcher.cb(actual.(*http.Request))

	if matcher.matching == nil {
		return format.Message(v, fmt.Sprintf(`to be a non-empty %s`, matcher.property))
	}
	return matcher.matching.FailureMessage(v)
}

func (matcher *ReqPropertyMatcher) NegatedFailureMessage(actual any) string {
	// todo: not so accurate
	return strings.Replace(matcher.FailureMessage(actual), "\nto ", "\nnot to ", 1)
}
