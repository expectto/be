package matchers

import (
	"fmt"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega/format"
	"reflect"
)

type AssignableToMatcher struct {
	assignableTo reflect.Type

	*MixinMatcherGomock
}

var _ types.BeMatcher = &AssignableToMatcher{}

func NewAssignableToMatcher[T any]() *AssignableToMatcher {
	t := reflect.TypeOf((*T)(nil)).Elem()

	im := &AssignableToMatcher{assignableTo: t}
	im.MixinMatcherGomock = NewMixinMatcherGomock(im, "AssignableTo")

	return im
}

func (matcher *AssignableToMatcher) Match(actual any) (success bool, err error) {
	if actual == nil {
		return false, nil
	}
	actualT := reflect.TypeOf(actual)
	return actualT.AssignableTo(matcher.assignableTo), nil
}

func (matcher *AssignableToMatcher) FailureMessage(actual any) string {
	return format.Message(actual, fmt.Sprintf("to be assignable to: %s", matcher.assignableTo.String()))
}

func (matcher *AssignableToMatcher) NegatedFailureMessage(actual any) string {
	return format.Message(actual, fmt.Sprintf("not to be assignable to: %s", matcher.assignableTo.String()))
}
