package psi_matchers

import (
	"fmt"
	"reflect"

	. "github.com/expectto/be/internal/psi" //nolint:staticcheck // should be moved to lintignore
	"github.com/expectto/be/types"
	"github.com/onsi/gomega/format"
)

type ImplementsMatcher struct {
	implements reflect.Type

	*MixinMatcherGomock
}

var _ types.BeMatcher = &ImplementsMatcher{}

func NewImplementsMatcher[T any]() *ImplementsMatcher {
	t := reflect.TypeOf((*T)(nil)).Elem()

	if t.Kind() != reflect.Interface {
		panic("ImplementsMatcher accepts interfaces to be given as T")
	}

	im := &ImplementsMatcher{implements: t}
	im.MixinMatcherGomock = NewMixinMatcherGomock(im, "Implements")

	return im
}

func (matcher *ImplementsMatcher) Match(actual any) (success bool, err error) {
	if actual == nil {
		return false, nil
	}
	actualT := reflect.TypeOf(actual)
	return actualT.Implements(matcher.implements), nil
}

func (matcher *ImplementsMatcher) FailureMessage(actual any) string {
	return format.Message(actual, fmt.Sprintf("to implement: %s", matcher.implements.String()))
}

func (matcher *ImplementsMatcher) NegatedFailureMessage(actual any) string {
	return format.Message(actual, fmt.Sprintf("not to implement: %s", matcher.implements.String()))
}
