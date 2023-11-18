package be_strings

import (
	"fmt"
	"github.com/IGLOU-EU/go-wildcard"
	"github.com/expectto/be/be_reflected"
	"github.com/expectto/be/internal/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gcustom"
)

// Todo: add aliases for some string-related reflected matchers

func NonEmptyString() types.BeMatcher {
	return psi_matchers.NewAllMatcher(
		be_reflected.AsString(),
		psi_matchers.NewNotMatcher(gomega.BeEmpty()),
	)
}

func EmptyString() types.BeMatcher {
	return psi_matchers.NewAllMatcher(
		be_reflected.AsString(),
		Psi(gomega.BeEmpty()),
	)
}

func Wildcard(pattern string) types.BeMatcher {
	return Psi(gcustom.MakeMatcher("Wildcard", func(actual interface{}) (bool, error) {
		if !cast.IsString(actual, cast.AllowCustomTypes(), cast.AllowPointers()) {
			return false, fmt.Errorf("string expected, got %T", actual)
		}

		return wildcard.Match(pattern, cast.AsString(actual)), nil
	}))
}

func Template(pattern string, args ...any) types.BeMatcher {
	return psi_matchers.NewAlwaysMatcher()
}

func With(key string, v any) types.BeMatcher {
	return psi_matchers.NewAlwaysMatcher()
}
