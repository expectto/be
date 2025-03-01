package psi_matchers

import (
	"fmt"
	"net/url"
	"strings"

	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/tmp/format"
	"github.com/expectto/be/types"
)

// UrlFieldMatcher is a helper for matching url fields
// It's considered to be used in matchers_url.go
type UrlFieldMatcher struct {
	publicName string // e.g. UrlHavingHost
	fieldName  string // e.g. "host"
	cb         func(*url.URL) any

	matching types.BeMatcher

	// todo: adjust gomock methods work as intended
	*MixinMatcherGomock
}

var _ types.BeMatcher = &UrlFieldMatcher{}

func NewUrlFieldMatcher(publicName, fieldName string, cb func(*url.URL) any, args ...any) *UrlFieldMatcher {
	matcher := &UrlFieldMatcher{
		publicName: publicName,
		fieldName:  fieldName,
		cb:         cb,
	}

	matcher.MixinMatcherGomock = NewMixinMatcherGomock(matcher, "Url field of")

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

func (matcher *UrlFieldMatcher) Match(actual any) (success bool, err error) {
	if actual == nil {
		return false, fmt.Errorf("%s() expects actual value not to be nil", matcher.publicName)
	}

	actualUrl, ok := actual.(*url.URL)
	if !ok {
		return false, fmt.Errorf("%s() expects actual value mast be a <*url.URL> received <%T>", matcher.publicName, actual)
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

func (matcher *UrlFieldMatcher) FailureMessage(actual any) string {
	v := matcher.cb(actual.(*url.URL))

	if matcher.matching == nil {
		return format.Message(v, fmt.Sprintf(`to be a non-empty %s`, matcher.fieldName))
	}
	return matcher.matching.FailureMessage(v)
}

func (matcher *UrlFieldMatcher) NegatedFailureMessage(actual any) string {
	// todo: not so accurate
	return strings.Replace(matcher.FailureMessage(actual), "\nto ", "\nnot to ", 1)
}
