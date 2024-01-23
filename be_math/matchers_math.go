// Package be_math provides Be matchers for mathematical operations
package be_math

import (
	"github.com/expectto/be/be_reflected"
	"github.com/expectto/be/internal/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gcustom"
	"math"
)

// GreaterThan succeeds if actual is numerically greater than the passed-in value.
func GreaterThan(arg any) types.BeMatcher {
	return Psi(gomega.BeNumerically(">", arg))
}

// GreaterThanEqual succeeds if actual is numerically greater than or equal to the passed-in value.
func GreaterThanEqual(arg any) types.BeMatcher {
	return Psi(gomega.BeNumerically(">=", arg))
}

// LessThan succeeds if actual is numerically less than the passed-in value.
func LessThan(arg any) types.BeMatcher {
	return Psi(gomega.BeNumerically("<", arg))
}

// LessThanEqual succeeds if actual is numerically less than or equal to the passed-in value.
func LessThanEqual(arg any) types.BeMatcher {
	return Psi(gomega.BeNumerically("<=", arg))
}

// Approx succeeds if actual is numerically approximately equal to the passed-in value within the specified threshold.
func Approx(compareTo, threshold any) types.BeMatcher {
	return Psi(gomega.BeNumerically("~", compareTo, threshold))
}

// InRange succeeds if actual is numerically within the specified range.
// The range is defined by the 'from' and 'until' values, and inclusivity is determined
// by the 'fromInclusive' and 'untilInclusive' flags.
func InRange(from any, fromInclusive bool, until any, untilInclusive bool) types.BeMatcher {
	group := make([]types.BeMatcher, 2)
	if fromInclusive {
		group[0] = Gte(from)
	} else {
		group[0] = Gt(from)
	}
	if untilInclusive {
		group[1] = Lte(until)
	} else {
		group[1] = Lt(until)
	}
	return psi_matchers.NewAllMatcher(cast.AsSliceOfAny(group)...)
}

// Odd succeeds if actual is an odd numeric value.
func Odd() types.BeMatcher {
	return Psi(
		be_reflected.AsNumeric(),
		WithFallibleTransform(func(actual any) any {
			return int(cast.AsFloat(actual))%2 != 0
		}, gomega.BeTrue()),
	)
}

// Even succeeds if actual is an even numeric value.
func Even() types.BeMatcher {
	return Psi(
		be_reflected.AsNumeric(),
		WithFallibleTransform(func(actual any) any {
			return int(cast.AsFloat(actual))%2 == 0
		}, gomega.BeTrue()),
	)
}

// Negative succeeds if actual is a negative numeric value.
func Negative() types.BeMatcher { return LessThan(0.0) }

// Positive succeeds if actual is a positive numeric value.
func Positive() types.BeMatcher { return GreaterThan(0.0) }

// Zero succeeds if actual is numerically approximately equal to zero.
// Any type of int/float will work for comparison.
func Zero() types.BeMatcher { return Approx(0, 0) }

// Integral succeeds if actual is an integral float, meaning it has zero decimal places.
// This matcher checks if the numeric value has no fractional component.
func Integral() types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		f := cast.AsFloat(actual)
		return f-float64(int(f)) == 0, nil
	}))
}

// DivisibleBy succeeds if actual is numerically divisible by the passed-in value.
func DivisibleBy(divisor any) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		return math.Mod(cast.AsFloat(actual), cast.AsFloat(divisor)) == 0, nil
	}))
}

// Shorter Names:

// Gt is an alias for GreaterThan, succeeding if actual is numerically greater than the passed-in value.
func Gt(arg any) types.BeMatcher { return GreaterThan(arg) }

// Gte is an alias for GreaterThanEqual, succeeding if actual is numerically greater than or equal to the passed-in value.
func Gte(arg any) types.BeMatcher { return GreaterThanEqual(arg) }

// Lt is an alias for LessThan, succeeding if actual is numerically less than the passed-in value.
func Lt(arg any) types.BeMatcher { return LessThan(arg) }

// Lte is an alias for LessThanEqual, succeeding if actual is numerically less than or equal to the passed-in value.
func Lte(arg any) types.BeMatcher { return LessThanEqual(arg) }
