package be

import (
	"github.com/expectto/be/internal/psi"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
)

// todo: implement gomega.Not as psi.Not
func BeNonEmptyString() types.BeMatcher {
	return All(
		String(),
		psi.Psi(gomega.Not(gomega.BeEmpty())),
	)
}

// todo implement... need to select third-party solution for wildcard
func MatchWildcard(wildcard string) types.BeMatcher {
	return nil
}
