// Package psi contains helpers that extends gomega library
// Name psi stands for previous letter from Omega
// (as we want to have a name that is close to gomega, but not to be a gomega)
// todo: rewrite better description
package psi

import (
	"github.com/expectto/be/types"
)

// Psi is a main converter function that converts given input into a PsiMatcher
func Psi(args ...any) types.BeMatcher {
	if len(args) == 0 {
		return Always()
	}
	if len(args) == 1 {
		return AsMatcher(args[0])
	}

	matchers := make([]types.BeMatcher, 0)

	// Cast each arg as:
	// 1. transform func: will be wrapped via WithFallibleTransform then
	// 2. Matcher (Gomega/Gomock/Psi)
	// 3. any raw value will be converted to EqualMatcher
	for i, arg := range args {
		if IsTransformFunc(arg) { // 1
			transformMatcher := WithFallibleTransform(arg, Psi(args[i+1:]...))
			matchers = append(matchers, Psi(transformMatcher))
			return All(matchers...)
		}

		matchers = append(matchers, Psi(arg)) // 2 or 3
	}

	return All(matchers...)
}
