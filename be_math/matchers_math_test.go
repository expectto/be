package be_math_test

import (
	"strings"

	"github.com/expectto/be/be_math"
	"github.com/expectto/be/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BeMath", func() {

	DescribeTable("should positively match", func(matcher types.BeMatcher, actual any) {
		// check gomega-compatible matching:
		success, err := matcher.Match(actual)
		Expect(err).Should(Succeed())
		Expect(success).To(BeTrue())

		// check gomock-compatible matching:
		success = matcher.Matches(actual)
		Expect(success).To(BeTrue())
	},
		Entry("10 GreaterThan 5", be_math.GreaterThan(5), 10),
		Entry("10 GreaterThan 5 (alias)", be_math.Gt(5), 10),
		Entry("10 GreaterThanEqual 5", be_math.GreaterThanEqual(5), 10),
		Entry("10 GreaterThanEqual 10", be_math.GreaterThanEqual(10), 10),
		Entry("10 GreaterThanEqual 5 (alias)", be_math.Gte(5), 10),
		Entry("10 GreaterThanEqual 10 (alias)", be_math.Gte(10), 10),

		Entry("10 LessThan 20", be_math.LessThan(20), 10),
		Entry("10 LessThan 20 (alias)", be_math.Lt(20), 10),
		Entry("10 LessThanEqual 20", be_math.LessThanEqual(20), 10),
		Entry("10 LessThanEqual 10", be_math.LessThanEqual(10), 10),
		Entry("10 LessThanEqual 20 (alias)", be_math.Lte(20), 10),
		Entry("10 LessThanEqual 10 (alias)", be_math.Lte(10), 10),

		Entry("3 is within range [1, 5]", be_math.InRange(1, true, 5, true), 3),
		Entry("3 is within range (1, 5]", be_math.InRange(1, false, 5, true), 3),
		Entry("3 is within range [1, 5)", be_math.InRange(1, true, 5, false), 3),
		Entry("3 is within range (1, 5)", be_math.InRange(1, false, 5, false), 3),

		Entry("1 is within range [1, 5]", be_math.InRange(1, true, 5, true), 1),
		Entry("5 is within range [1, 5]", be_math.InRange(1, true, 5, true), 5),
		Entry("1 is within range [1, 5)", be_math.InRange(1, true, 5, false), 1),

		Entry("0.999 is ~ 1.0 within threshold 0.01", be_math.Approx(0.999, 0.01), 1.0),
		Entry("3.05 is ~ 3.04 within threshold 0.01", be_math.Approx(3.05, 0.01), 3.04),
		Entry("3.5 is ~ 3.4 within threshold 0.2", be_math.Approx(3.5, 0.2), 3.4),
		Entry("3.5 is ~ 3.0 within threshold 0.5", be_math.Approx(3.5, 0.5), 3.0),
		Entry("3.9 is ~ 3.0 within threshold 1.0", be_math.Approx(3.9, 1.0), 3.0),
		Entry("4.0 is ~ 3.0 within threshold 1.0", be_math.Approx(4.0, 1.0), 3.0),
		Entry("3.5 is ~ 3.4 within threshold 10.0", be_math.Approx(3.5, 10.0), 3.4),

		Entry("3 is an odd number", be_math.Odd(), 3),
		Entry("7 is an odd number", be_math.Odd(), 7),
		Entry("-3 is an odd number", be_math.Odd(), -3),
		Entry("-7 is an odd number", be_math.Odd(), -7),

		Entry("2 is an even number", be_math.Even(), 2),
		Entry("4 is an even number", be_math.Even(), 4),
		Entry("-2 is an even number", be_math.Even(), -2),
		Entry("-4 is an even number", be_math.Even(), -4),

		Entry("-5 is a negative number", be_math.Negative(), -5),
		Entry("-8.5 is a negative number", be_math.Negative(), -8.5),
		Entry("5 is a positive number", be_math.Positive(), 5),
		Entry("8.5 is a positive number", be_math.Positive(), 8.5),

		Entry("0 is zero", be_math.Zero(), 0.0),
		Entry("0.00000000009 is approx zero", be_math.ApproxZero(), 0.00000000009),

		Entry("5 is an integral number", be_math.Integral(), 5.0),
		Entry("-10 is an integral number", be_math.Integral(), -10.0),

		Entry("10 is divisible by 5", be_math.DivisibleBy(5), 10),
		Entry("18 is divisible by -3", be_math.DivisibleBy(-3), 18),
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
		Entry("5 is not GreaterThan 10", be_math.GreaterThan(10), 5),
		Entry("5 is not GreaterThan 5 (alias)", be_math.Gt(5), 5),
		Entry("5 is not GreaterThanEqual 10", be_math.GreaterThanEqual(10), 5),
		Entry("5 is not GreaterThanEqual 20", be_math.GreaterThanEqual(20), 5),
		Entry("5 is not GreaterThanEqual 10 (alias)", be_math.Gte(10), 5),
		Entry("5 is not GreaterThanEqual 20 (alias)", be_math.Gte(20), 5),

		Entry("20 is not LessThan 10", be_math.LessThan(10), 20),
		Entry("20 is not LessThan 10 (alias)", be_math.Lt(10), 20),
		Entry("20 is not LessThanEqual 10", be_math.LessThanEqual(10), 20),
		Entry("20 is not LessThanEqual 5", be_math.LessThanEqual(5), 20),
		Entry("20 is not LessThanEqual 10 (alias)", be_math.Lte(10), 20),
		Entry("20 is not LessThanEqual 5 (alias)", be_math.Lte(5), 20),

		Entry("6 is not within range [1, 5]", be_math.InRange(1, true, 5, true), 6),
		Entry("6 is not within range (1, 5]", be_math.InRange(1, false, 5, true), 6),
		Entry("6 is not within range [1, 5)", be_math.InRange(1, true, 5, false), 6),
		Entry("6 is not within range (1, 5)", be_math.InRange(1, false, 5, false), 6),

		Entry("0.999 is not ~ 1.0 within threshold 0.001", be_math.Approx(1.0, 0.001), 0.999),
		Entry("3.05 is not ~ 3.04 within threshold 0.001", be_math.Approx(3.04, 0.001), 3.05),
		Entry("3.5 is not ~ 3.0 within threshold 0.2", be_math.Approx(3.0, 0.2), 3.5),
		Entry("3.9 is not ~ 3.0 within threshold 0.5", be_math.Approx(3.0, 0.5), 3.9),
		Entry("4.0 is not ~ 3.0 within threshold 0.5", be_math.Approx(3.0, 0.5), 4.0),

		Entry("2 is not an odd number", be_math.Odd(), 2),
		Entry("4 is not an odd number", be_math.Odd(), 4),
		Entry("-2 is not an odd number", be_math.Odd(), -2),
		Entry("-4 is not an odd number", be_math.Odd(), -4),
		Entry("floats can't be matched as odd numbers", be_math.Odd(), 1.5),

		Entry("3 is not an even number", be_math.Even(), 3),
		Entry("7 is not an even number", be_math.Even(), 7),
		Entry("-3 is not an even number", be_math.Even(), -3),
		Entry("-7 is not an even number", be_math.Even(), -7),
		Entry("floats can't be matched as even numbers", be_math.Even(), 1.5),

		Entry("5 is not a negative number", be_math.Negative(), 5),
		Entry("8.5 is not a negative number", be_math.Negative(), 8.5),
		Entry("0 is not a negative number", be_math.Negative(), 0),
		Entry("-5 is not a positive number", be_math.Positive(), -5),
		Entry("-8.5 is not a positive number", be_math.Positive(), -8.5),

		Entry("0.1 is not zero", be_math.Zero(), 0.1),
		Entry("0.1 is not approx zero", be_math.ApproxZero(), 0.1),

		Entry("5.5 is not an integral number", be_math.Integral(), 5.5),
		Entry("-10.5 is not an integral number", be_math.Integral(), -10.5),

		Entry("10 is not divisible by 3", be_math.DivisibleBy(3), 10),
		Entry("18 is not divisible by -4", be_math.DivisibleBy(-4), 18),
	)

	DescribeTable("should return a valid failure message", func(matcher types.BeMatcher, actual any, message string) {
		// FailureMessage is considered to be called after matching:
		_, _ = matcher.Match(actual)

		failureMessage := matcher.FailureMessage(actual)
		Expect(failureMessage).To(Equal(message))

		// in all our matchers negated failure messages are simply `to be` => `not to be`
		Expect(matcher.NegatedFailureMessage(actual)).To(Equal(
			strings.Replace(failureMessage, "\nto be ", "\nnot to be ", 1),
		))
	},
		// Example of entry where FailureMessage is simply inherited from gomega's underlying matching
		Entry("5 is not GreaterThan 10", be_math.GreaterThan(10), 5, "Expected\n    <int>: 5\nto be >\n    <int>: 10"),

		// Examples of entry with custom message (gcustom.MakeMatcher matching)
		Entry("10 is not divisible by 3", be_math.DivisibleBy(3), 10, "Expected:\n    <int>: 10\nto be divisible by 3"),
		Entry("0.1 is not zero", be_math.Zero(), 0.1, "Expected:\n    <float64>: 0.1\nto be zero"),

		// Examples of entry on complex Psi matchers (chaining + transform)
		Entry("float is not odd", be_math.Odd(), 12.5, "Expected:\n    <float64>: 12.5\nto be an odd number"),
		Entry("8 is not odd", be_math.Odd(), 8, "Expected:\n    <int>: 8\nto be an odd number"),
		Entry("8 (uint) is not odd", be_math.Odd(), uint(8), "Expected:\n    <uint>: 8\nto be an odd number"),
	)
})
