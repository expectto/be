package be_jwt

import (
	"fmt"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
)

func Token(args ...any) types.BeMatcher {
	return psi_matchers.NewNeverMatcher(fmt.Errorf("todo: not implemented"))
}

func HavingClaims(args ...any) types.BeMatcher {
	return psi_matchers.NewNeverMatcher(fmt.Errorf("todo: not implemented"))
}

func BeingValid() types.BeMatcher {
	return psi_matchers.NewNeverMatcher(fmt.Errorf("todo: not implemented"))
}
