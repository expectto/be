package psi_matchers

import (
	"fmt"
	"strings"

	"github.com/amberpixels/abu/cast"
	"github.com/amberpixels/abu/reflectish"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega/format"
)

// HaveLengthMatcher is an Omega-format matcher that matches length of the given list
// in comparison to either a given int value, or matching to other given matchers
type HaveLengthMatcher struct {
	count    *int
	matching types.BeMatcher

	*MixinMatcherGomock
}

var _ types.BeMatcher = &HaveLengthMatcher{}

func NewHaveLengthMatcher(args ...any) *HaveLengthMatcher {
	if len(args) == 0 {
		panic("HaveLenMatcher requires an int or list of matcher be given")
	}

	matcher := &HaveLengthMatcher{}
	matcher.MixinMatcherGomock = NewMixinMatcherGomock(matcher, "HaveLen")

	// if just argument was given, and it's not a matcher, then it must be an integer
	if len(args) == 1 && !IsMatcher(args[0]) {
		matcher.count = new(int)
		*matcher.count = cast.AsInt(args[0])
		return matcher
	}

	// compress all given args as group of matchers
	matcher.matching = Psi(args...)
	return matcher
}

func (matcher *HaveLengthMatcher) Match(actual any) (success bool, err error) {
	length, ok := reflectish.LengthOf(actual)
	if !ok {
		return false, fmt.Errorf("HaveLen matcher expects a string/array/map/channel/slice.  Got:\n%s", format.Object(actual, 1))
	}

	if matcher.count != nil {
		return length == *matcher.count, nil
	}

	return matcher.matching.Match(length)
}

func (matcher *HaveLengthMatcher) FailureMessage(actual any) (message string) {
	if matcher.count != nil {
		return fmt.Sprintf("Expected\n%s\nto have length = %d", format.Object(actual, 1), *matcher.count)
	}

	// Assuming that underlying message will be the same format
	// change:
	//		Expect [ ] to <>
	// into
	//		Expect [ ] length to <>
	failureMessage := matcher.matching.FailureMessage(actual)
	failureMessage = strings.Replace(failureMessage, "\nto", "\nlength to", 1)
	return failureMessage
}

func (matcher *HaveLengthMatcher) NegatedFailureMessage(actual any) (message string) {
	if matcher.count != nil {
		return fmt.Sprintf("Expected\n%s\nnot to have length = %d", format.Object(actual, 1), *matcher.count)
	}

	// Assuming that underlying message will be the same format
	// change:
	//		Expect [ ] not to <>
	// into
	//		Expect [ ] length not to <>
	failureMessage := matcher.matching.FailureMessage(actual)
	failureMessage = strings.Replace(failureMessage, "\nnot to", "\nlength not to", 1)
	return failureMessage
}
