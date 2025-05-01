package psi

import (
	"fmt"

	"github.com/amberpixels/abu/cast"
)

type DiveMode string

const (
	DiveModeEvery DiveMode = "every"
	DiveModeAny   DiveMode = "any"
	DiveModeFirst DiveMode = "first"
	DiveModeNth   DiveMode = "nth"
)

type DiveMatcher struct {
	matcher any
	mode    DiveMode

	// when mode is DiveModeNth, then we keep nth element
	n int

	*MixinMatcherGomock
}

func NewDiveMatcher(matcher any, mode DiveMode, args ...any) *DiveMatcher {
	dm := &DiveMatcher{matcher: matcher, mode: mode}

	if mode == DiveModeNth {
		if len(args) == 0 {
			panic("DiveNth expects value of `n` as an argument")
		}
		if !cast.IsInt(args[0]) {
			panic("DiveNth expects value of `n` to be an integer")
		}

		dm.n = cast.AsInt(args[0])
	}

	return dm
}

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
	case DiveModeNth:
		if len(slice) == 0 {
			return false, fmt.Errorf("dive[nth] expects non-empty slice")
		}
		if dm.n >= len(slice) {
			return false, fmt.Errorf("dive[nth] expects `n` to be less than length of slice")
		}
		return matcher.Match(slice[dm.n])
	}

	panic("invalid DeepMatcher mode")
}

func (dm *DiveMatcher) FailureMessage(actual any) string {
	return fmt.Sprintf("to %s on %s of given list", Psi(dm.matcher).FailureMessage(actual), dm.mode)
}
func (dm *DiveMatcher) NegatedFailureMessage(actual any) string {
	return "not " + dm.FailureMessage(actual)
}
