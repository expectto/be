package be_ctx

import (
	"github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
)

func Ctx(args ...any) types.BeMatcher {
	if len(args) == 0 {
		return psi_matchers.NewCtxMatcher()
	}

	// not sure about this:
	return psi.Psi(args...)
}

func CtxWithValue(key string, vArg ...any) types.BeMatcher {
	return psi_matchers.NewCtxValueMatcher(key, vArg...)
}

func CtxWithDeadline(deadline any) types.BeMatcher {
	return psi_matchers.NewCtxDeadlineMatcher(deadline)
}

func CtxWithError(err any) types.BeMatcher {
	return psi_matchers.NewCtxErrMatcher(err)
}
