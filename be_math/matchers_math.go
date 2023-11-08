package be_math

import (
	"github.com/expectto/be/be_reflected"
	"github.com/expectto/be/internal/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
)

// Todo: Integral (means its an integral float, so zero decimals)
// todo: DivisibleBy
// todo: Negative & Positive (syntax sugar)

func GreaterThan(arg any) types.BeMatcher {
	return Psi(gomega.BeNumerically(">", arg))
}

func GreaterThanEqual(arg any) types.BeMatcher {
	return Psi(gomega.BeNumerically(">=", arg))
}

func LessThan(arg any) types.BeMatcher {
	return Psi(gomega.BeNumerically("<", arg))
}

func LessThanEqual(arg any) types.BeMatcher {
	return Psi(gomega.BeNumerically("<=", arg))
}

func ApproxEqual(compareTo, threshold any) types.BeMatcher {
	return Psi(gomega.BeNumerically("~", compareTo, threshold))
}

func InRange(from, fromInclusive bool, until any, untilInclusive bool) types.BeMatcher {
	group := make([]types.BeMatcher, 2, 2)
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

func Odd() types.BeMatcher {
	return Psi(
		be_reflected.AsNumeric(),
		WithFallibleTransform(func(actual any) any {
			// todo: not accurate!
			return int(cast.AsFloat(actual))%2 != 0
		}, gomega.BeTrue()),
	)
}

func Even() types.BeMatcher {
	return Psi(
		be_reflected.AsNumeric(),
		WithFallibleTransform(func(actual any) any {
			return int(cast.AsFloat(actual))%2 == 0
		}, gomega.BeTrue()),
	)
}

// Shorter Names:

func Gt(arg any) types.BeMatcher      { return GreaterThan(arg) }
func Gte(arg any) types.BeMatcher     { return GreaterThanEqual(arg) }
func Lt(arg any) types.BeMatcher      { return LessThan(arg) }
func Lte(arg any) types.BeMatcher     { return LessThanEqual(arg) }
func Approx(c, t any) types.BeMatcher { return ApproxEqual(c, t) }
