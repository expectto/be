package be_ctx_test

import (
	"context"
	"github.com/expectto/be/be_ctx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MatchersCtx", func() {
	ctx := context.Background()

	It("should match a ctx", func() {
		Expect(ctx).To(be_ctx.Ctx())
	})

	It("should match a ctx with a value", func() {
		ctx := context.WithValue(ctx, "foo", "bar")
		// just by key
		Expect(ctx).To(be_ctx.CtxWithValue("foo"))
		// key + value directly
		Expect(ctx).To(be_ctx.CtxWithValue("foo", "bar"))
		// key + value via matcher
		Expect(ctx).To(be_ctx.CtxWithValue("foo", HavePrefix("ba")))
	})

	It("should not match when a string given instead of ctx", func() {
		Expect("not a ctx but a string").NotTo(be_ctx.Ctx())
	})
})
