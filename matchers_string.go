package be

import (
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
)

func BeNonEmptyString() types.BeMatcher {
	return All(
		String(),
		Not(gomega.BeEmpty()),
	)
}

// todo implement... need to select third-party solution for wildcard
func MatchWildcard(wildcard string) types.BeMatcher {
	return nil
}
