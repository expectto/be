package psi_matchers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/amberpixels/abu/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega/format"
)

//
// KindMatcher:
//

type KindMatcher struct {
	kind *reflect.Kind

	matching types.BeMatcher

	*MixinMatcherGomock
}

var _ types.BeMatcher = &KindMatcher{}

func NewKindMatcher(args ...any) *KindMatcher {
	if len(args) == 0 {
		panic("KindMatcher requires a reflect.Kind or list of matcher be given")
	}

	matcher := &KindMatcher{}
	if len(args) == 1 && !IsMatcher(args[0]) {
		matcher.kind = new(reflect.Kind)
		*(matcher.kind) = cast.AsKind(args[0])
		return matcher
	}

	matcher.matching = Psi(args...)
	matcher.MixinMatcherGomock = NewMixinMatcherGomock(matcher, "Kind of")
	return matcher
}

func (matcher *KindMatcher) Match(actual any) (success bool, err error) {
	if actual == nil {
		return false, nil
	}

	if matcher.kind != nil {
		return *matcher.kind == reflect.TypeOf(actual).Kind(), nil
	}

	return matcher.matching.Match(reflect.TypeOf(actual).Kind())
}

func (matcher *KindMatcher) FailureMessage(actual any) string {
	if matcher.kind != nil {
		return format.Message(actual, fmt.Sprintf("to be kind of %s", matcher.kind.String()))
	}

	// Assuming that underlying message will be the same format
	// change:
	//		Expect [ ] to <>
	// into
	//		Expect [ ] kind to <>
	// Note: probably a weak solution: consider better
	failureMessage := matcher.matching.FailureMessage(actual)
	failureMessage = strings.Replace(failureMessage, "\nto", "\nkind to", 1)
	return failureMessage
}

func (matcher *KindMatcher) NegatedFailureMessage(actual any) string {
	if matcher.kind != nil {
		return format.Message(actual, fmt.Sprintf("not to be kind of %s", matcher.kind.String()))
	}

	failureMessage := matcher.matching.NegatedFailureMessage(actual)
	failureMessage = strings.Replace(failureMessage, "\nnot to", "\nkind not to", 1)
	return failureMessage
}
