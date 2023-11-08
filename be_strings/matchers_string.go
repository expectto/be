package be_strings

import (
	"github.com/expectto/be/be_reflected"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
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

// todo implement... need to select third-party solution for wildcard
func Wildcard(wildcard string) types.BeMatcher {
	panic("not implemented")
}
