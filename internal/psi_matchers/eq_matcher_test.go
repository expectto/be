package psi_matchers_test

import (
	. "github.com/expectto/be/internal/psi_matchers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EqMatcher", func() {

	DescribeTable("Match",
		func(expected, actual any, expectedResult, shouldFail bool) {
			matcher := NewEqMatcher(expected)
			result, err := matcher.Match(actual)
			Expect(result).To(Equal(expectedResult))
			Expect(err != nil).To(Equal(shouldFail))

			// Matches() is a shortcut to Match
			Expect(matcher.Matches(actual)).To(Equal(result))
		},
		Entry("Nil values (are note comparable)", nil, nil, false, true),
		Entry("Equal strings", "hello", "hello", true, false),
		Entry("Equal byte slices", []byte{1, 2, 3}, []byte{1, 2, 3}, true, false),
		Entry("Not equal strings", "hello", "world", false, false),
		Entry("Not equal byte slices", []byte{1, 2, 3}, []byte{3, 2, 1}, false, false),
	)

	Describe("Failure", func() {
		It("should return the failure message for non-equal strings", func() {
			matcher := NewEqMatcher("world")
			failureMsg := matcher.FailureMessage("hello")

			Expect(failureMsg).To(Equal("Expected\n    <string>: hello\nto equal\n    <string>: world"))
		})

		It("should return the failure String() for non-equal strings", func() {
			matcher := NewEqMatcher("world")
			_ = matcher.Matches("hello") // actual value is passed through Matches9)
			failureStr := matcher.String()

			Expect(failureStr).To(Equal("Expected\n    <string>: hello\nto equal\n    <string>: world"))
		})
	})
})
