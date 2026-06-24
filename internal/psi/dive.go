package psi

import (
	"fmt"
	"reflect"

	"github.com/amberpixels/k1/cast"
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

func (dm *DiveMatcher) Match(actual any) (bool, error) {
	matcher := Psi(dm.matcher)

	// Collect the elements to dive over via reflection. Slices and arrays dive
	// over their elements; maps dive over their values. Anything else fails
	// gracefully instead of panicking (cast.AsSliceOfAny would panic).
	rv := reflect.ValueOf(actual)
	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	var slice []any
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		slice = make([]any, rv.Len())
		for i := range slice {
			slice[i] = rv.Index(i).Interface()
		}
	case reflect.Map:
		// Maps are unordered, so positional modes are not meaningful.
		if dm.mode == DiveModeFirst || dm.mode == DiveModeNth {
			return false, fmt.Errorf("dive[%s] is not supported on a map (maps are unordered)", dm.mode)
		}
		slice = make([]any, 0, rv.Len())
		for _, k := range rv.MapKeys() {
			slice = append(slice, rv.MapIndex(k).Interface())
		}
	default:
		return false, fmt.Errorf("dive[%s] expects a slice, array or map, got %T", dm.mode, actual)
	}

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
