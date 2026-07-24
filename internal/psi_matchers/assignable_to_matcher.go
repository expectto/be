package psi_matchers

import (
	"reflect"

	"github.com/onsi/gomega/format"

	. "github.com/expectto/be/internal/psi" //nolint:staticcheck // should be moved to lintignore
	"github.com/expectto/be/types"
)

type AssignableToMatcher struct {
	*MixinMatcherGomock

	assignableTo reflect.Type
}

var _ types.BeMatcher = &AssignableToMatcher{}

func NewAssignableToMatcher[T any]() *AssignableToMatcher {
	t := reflect.TypeFor[T]()

	im := &AssignableToMatcher{assignableTo: t}
	im.MixinMatcherGomock = NewMixinMatcherGomock(im, "AssignableTo")

	return im
}

func (matcher *AssignableToMatcher) Match(actual any) (bool, error) {
	if actual == nil {
		return false, nil
	}
	actualT := reflect.TypeOf(actual)
	return actualT.AssignableTo(matcher.assignableTo), nil
}

func (matcher *AssignableToMatcher) FailureMessage(actual any) string {
	return format.Message(actual, "to be assignable to: "+matcher.assignableTo.String())
}

func (matcher *AssignableToMatcher) NegatedFailureMessage(actual any) string {
	return format.Message(actual, "not to be assignable to: "+matcher.assignableTo.String())
}
