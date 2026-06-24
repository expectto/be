package be_ctx_test

import (
	"context"
	"time"

	"github.com/expectto/be/be_ctx"
	"github.com/expectto/be/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// helpers building context.Context actuals for the tables below.

// plainCtx returns a fresh background context.
func plainCtx() context.Context { return context.Background() }

// valueCtx returns a context carrying the given key/value pair.
func valueCtx(key, value any) context.Context {
	return context.WithValue(context.Background(), key, value)
}

// fixedDeadline is a stable deadline shared between actual/expected entries.
var fixedDeadline = time.Date(2030, time.January, 1, 0, 0, 0, 0, time.UTC)

// deadlineCtx returns a context with the fixed deadline (cancel intentionally
// dropped: the context lives only for the duration of a single matcher call).
func deadlineCtx() context.Context {
	ctx, cancel := context.WithDeadline(context.Background(), fixedDeadline)
	_ = cancel
	return ctx
}

// timeoutCtx returns a context built via WithTimeout (i.e. it does have a deadline).
func timeoutCtx() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	_ = cancel
	return ctx
}

// canceledCtx returns a context that has already been canceled.
func canceledCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

// deadlineExceededCtx returns a context whose deadline is already in the past,
// so ctx.Err() is context.DeadlineExceeded.
func deadlineExceededCtx() context.Context {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Hour))
	_ = cancel
	return ctx
}

var _ = Describe("BeCtx", func() {

	DescribeTable("should positively match", func(matcher types.BeMatcher, actual any) {
		// check gomega-compatible matching:
		success, err := matcher.Match(actual)
		Expect(err).Should(Succeed())
		Expect(success).To(BeTrue())

		// check gomock-compatible matching:
		success = matcher.Matches(actual)
		Expect(success).To(BeTrue())
	},
		// Ctx() with no args matches any context.Context:
		Entry("background ctx is a ctx", be_ctx.Ctx(), plainCtx()),
		Entry("TODO ctx is a ctx", be_ctx.Ctx(), context.TODO()),
		Entry("value ctx is a ctx", be_ctx.Ctx(), valueCtx("k", "v")),
		Entry("canceled ctx is a ctx", be_ctx.Ctx(), canceledCtx()),

		// CtxWithValue: existence-only matching (no value arg):
		Entry("ctx has value for key `foobar`", be_ctx.CtxWithValue("foobar"), valueCtx("foobar", 100)),
		Entry("ctx has value for key `name`", be_ctx.CtxWithValue("name"), valueCtx("name", "alice")),

		// CtxWithValue: exact value matching:
		Entry("ctx value foobar==100", be_ctx.CtxWithValue("foobar", 100), valueCtx("foobar", 100)),
		Entry("ctx value name==alice", be_ctx.CtxWithValue("name", "alice"), valueCtx("name", "alice")),

		// CtxWithDeadline: exact deadline matching against the actual time.Time:
		Entry("ctx deadline matches fixed deadline", be_ctx.CtxWithDeadline(fixedDeadline), deadlineCtx()),

		// CtxWithError: cancellation / deadline-exceeded errors:
		Entry("canceled ctx err is context.Canceled", be_ctx.CtxWithError(context.Canceled), canceledCtx()),
		Entry("expired ctx err is context.DeadlineExceeded", be_ctx.CtxWithError(context.DeadlineExceeded), deadlineExceededCtx()),
		Entry("live ctx err is nil", be_ctx.CtxWithError(nil), plainCtx()),
	)

	DescribeTable("should negatively match", func(matcher types.BeMatcher, actual any) {
		// check gomega-compatible matching:
		success, err := matcher.Match(actual)
		Expect(err).Should(Succeed())
		Expect(success).To(BeFalse())

		// check gomock-compatible matching:
		success = matcher.Matches(actual)
		Expect(success).To(BeFalse())
	},
		// not a context.Context at all:
		Entry("string is not a ctx", be_ctx.CtxWithValue("k"), "not-a-ctx"),
		Entry("int is not a ctx", be_ctx.CtxWithError(nil), 42),

		// CtxWithValue: missing key:
		Entry("ctx misses key `absent`", be_ctx.CtxWithValue("absent"), valueCtx("foobar", 100)),
		Entry("background ctx has no values", be_ctx.CtxWithValue("foobar"), plainCtx()),

		// CtxWithValue: wrong value:
		Entry("ctx value foobar!=200", be_ctx.CtxWithValue("foobar", 200), valueCtx("foobar", 100)),
		Entry("ctx value name!=bob", be_ctx.CtxWithValue("name", "bob"), valueCtx("name", "alice")),

		// CtxWithDeadline: no deadline present:
		Entry("background ctx has no deadline", be_ctx.CtxWithDeadline(fixedDeadline), plainCtx()),
		// CtxWithDeadline: wrong deadline:
		Entry("ctx deadline differs from expected", be_ctx.CtxWithDeadline(fixedDeadline.Add(time.Hour)), deadlineCtx()),

		// CtxWithError: wrong / unexpected error state:
		Entry("live ctx err is not context.Canceled", be_ctx.CtxWithError(context.Canceled), plainCtx()),
		Entry("canceled ctx err is not DeadlineExceeded", be_ctx.CtxWithError(context.DeadlineExceeded), canceledCtx()),

		// timeoutCtx does have a deadline but is still live (err == nil), so it
		// must NOT match an expectation of context.Canceled:
		Entry("timeout ctx err is not context.Canceled", be_ctx.CtxWithError(context.Canceled), timeoutCtx()),
	)

	DescribeTable("should return a valid failure message", func(matcher types.BeMatcher, actual any, message string) {
		// FailureMessage is considered to be called after matching:
		_, _ = matcher.Match(actual)

		Expect(matcher.FailureMessage(actual)).To(Equal(message))
	},
		Entry("not a context.Context",
			be_ctx.CtxWithError(nil), "not-a-ctx",
			"Expected\n    <string>: not-a-ctx\nto be a ctx",
		),
		Entry("missing ctx value key",
			be_ctx.CtxWithValue("absent"), plainCtx(),
			"Expected\n    <context.backgroundCtx>: {\n        emptyCtx: <suppressed context>,\n    }\nto have the ctx.value key=`absent`",
		),
	)

	// CtxWithError(nil) asserts the context carries NO error. It must reject a
	// context that has errored (previously it matched any context).
	DescribeTable("CtxWithError(nil) asserts no error", func(actual context.Context, wantMatch bool) {
		success, err := be_ctx.CtxWithError(nil).Match(actual)
		Expect(err).Should(Succeed())
		Expect(success).To(Equal(wantMatch))
	},
		Entry("live context has no error -> match", plainCtx(), true),
		Entry("canceled context has an error -> no match", canceledCtx(), false),
		Entry("deadline-exceeded context has an error -> no match", deadlineExceededCtx(), false),
	)
})
