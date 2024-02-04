package psi

import (
	"fmt"
	"github.com/expectto/be/internal/cast"
	"github.com/expectto/be/types"
)

type DiveMode string

const (
	DiveModeEvery DiveMode = "every"
	DiveModeAny   DiveMode = "any"
	DiveModeFirst DiveMode = "first"
)

type DiveMatcher struct {
	matcher any
	mode    DiveMode

	*MixinMatcherGomock
}

func NewDiveMatcher(matcher any, mode DiveMode) *DiveMatcher {
	return &DiveMatcher{matcher: matcher, mode: mode}
}

// Dive applies the given matcher to each (every) element of the slice.
// Note: Dive is very close to gomega.HaveEach
func Dive(matcher any) types.BeMatcher { return NewDiveMatcher(matcher, DiveModeEvery) }

// DiveAny applies the given matcher to each element and succeeds in case if it succeeds at least at one item
func DiveAny(matcher any) types.BeMatcher { return NewDiveMatcher(matcher, DiveModeAny) }

// DiveFirst applies the given matcher to the first element of the given slice
func DiveFirst(matcher any) types.BeMatcher { return NewDiveMatcher(matcher, DiveModeFirst) }

func (dm *DiveMatcher) Match(actual interface{}) (bool, error) {
	matcher := Psi(dm.matcher)

	// todo: nice error if actual is not a slice-ish
	// as other way it panics
	slice := cast.AsSliceOfAny(actual)

	switch dm.mode {
	case DiveModeEvery:
		if len(slice) == 0 {
			return false, nil
		}

		for _, item := range slice {
			success, err := matcher.Match(item)
			if err != nil {
				return false, err
			}
			if !success {
				return false, nil
			}
		}
		return true, nil

	case DiveModeAny:
		if len(slice) == 0 {
			return true, nil
		}

		for _, item := range slice {
			success, err := matcher.Match(item)
			if err != nil {
				return false, err
			}
			if success {
				return true, nil
			}
		}

		return false, nil
	case DiveModeFirst:
		if len(slice) == 0 {
			return false, fmt.Errorf("dive[first] expects non-empty slice")
		}
		return matcher.Match(slice[0])
	}

	panic("invalid DeepMatcher mode")
}

func (dm *DiveMatcher) FailureMessage(actual any) string {
	return fmt.Sprintf("to %s on %s of given list", Psi(dm.matcher).FailureMessage(actual), dm.mode)
}
func (dm *DiveMatcher) NegatedFailureMessage(actual any) string {
	return "not " + dm.FailureMessage(actual)
}
