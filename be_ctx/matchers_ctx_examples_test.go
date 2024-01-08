package be_ctx_test

import (
	"context"
	"github.com/expectto/be/be_ctx"
	"github.com/expectto/be/internal/psi_matchers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MatchersCtx", func() {
	ctx := context.Background()

	It("should match a ctx", func() {
		Expect(ctx).To(be_ctx.NewCtxMatcher())
	})

	It("should match a ctx with a value", func() {
		ctx := context.WithValue(ctx, "foo", "bar")
		Expect(ctx).To(be_ctx.NewCtxValueMatcher("foo"))
		Expect(ctx).To(be_ctx.NewCtxValueMatcher("foo", psi_matchers.NewEqMatcher("bar")))
	})
})
