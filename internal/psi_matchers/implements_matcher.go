package psi_matchers

import (
	"reflect"

	"github.com/onsi/gomega/format"

	. "github.com/expectto/be/internal/psi" //nolint:staticcheck // should be moved to lintignore
	"github.com/expectto/be/types"
)

type ImplementsMatcher struct {
	*MixinMatcherGomock

	implements reflect.Type
}

var _ types.BeMatcher = &ImplementsMatcher{}

func NewImplementsMatcher[T any]() *ImplementsMatcher {
	t := reflect.TypeFor[T]()

	if t.Kind() != reflect.Interface {
		panic("ImplementsMatcher accepts interfaces to be given as T")
	}

	im := &ImplementsMatcher{implements: t}
	im.MixinMatcherGomock = NewMixinMatcherGomock(im, "Implements")

	return im
}

func (matcher *ImplementsMatcher) Match(actual any) (bool, error) {
	if actual == nil {
		return false, nil
	}
	actualT := reflect.TypeOf(actual)
	return actualT.Implements(matcher.implements), nil
}

func (matcher *ImplementsMatcher) FailureMessage(actual any) string {
	return format.Message(actual, "to implement: "+matcher.implements.String())
}

func (matcher *ImplementsMatcher) NegatedFailureMessage(actual any) string {
	return format.Message(actual, "not to implement: "+matcher.implements.String())
}
