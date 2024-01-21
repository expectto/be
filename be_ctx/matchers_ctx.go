package be_ctx

import (
	"github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
)

// Ctx succeeds if the actual value is a context.Context.
// If no arguments are provided, it matches any context.Context.
// Otherwise, it uses the Psi matcher to match the provided arguments against the actual context's values.
func Ctx(args ...any) types.BeMatcher {
	if len(args) == 0 {
		return psi_matchers.NewCtxMatcher()
	}

	// todo: weak solution, fixme
	return psi.Psi(args...)
}

// CtxWithValue succeeds if the actual value is a context.Context and contains a key-value pair
// where the key matches the provided key and the value matches the provided arguments using any other matchers.
func CtxWithValue(key any, vs ...any) types.BeMatcher {
	return psi_matchers.NewCtxValueMatcher(key, vs...)
}

// CtxWithDeadline succeeds if the actual value is a context.Context and its deadline matches the provided deadline.
func CtxWithDeadline(deadline any) types.BeMatcher {
	return psi_matchers.NewCtxDeadlineMatcher(deadline)
}

// CtxWithError succeeds if the actual value is a context.Context and its error matches the provided error value.
func CtxWithError(err any) types.BeMatcher {
	return psi_matchers.NewCtxErrMatcher(err)
}
