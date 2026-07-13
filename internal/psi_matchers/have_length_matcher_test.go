package psi_matchers_test

import (
	. "github.com/expectto/be/internal/psi_matchers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HaveLengthMatcher", func() {
	Context("Match", func() {
		It("should match the length correctly", func() {
			actual := []int{1, 2, 3}
			matcher := NewHaveLengthMatcher(3)

			Expect(matcher.Match(actual)).To(BeTrue())
		})

		It("should fail for incorrect length", func() {
			actual := []int{1, 2}
			matcher := NewHaveLengthMatcher(3)

			Expect(matcher.Match(actual)).To(BeFalse())
		})

		It("should handle invalid input type", func() {
			actual := 123 // Invalid type
			matcher := NewHaveLengthMatcher(3)
			_, err := matcher.Match(actual)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("HaveLen matcher expects a string/array/map/channel/slice"))
		})

		It("should match length against a matcher instead of a count", func() {
			matcher := NewHaveLengthMatcher(BeNumerically(">", 5))

			Expect(matcher.Match([]int{1, 2, 3, 4, 5, 6})).To(BeTrue())
			Expect(matcher.Match([]int{1, 2, 3})).To(BeFalse())
			Expect(matcher.Match("long enough")).To(BeTrue())
		})
	})

	Context("FailureMessage", func() {
		It("should return the failure message for incorrect length", func() {
			matcher := NewHaveLengthMatcher(3)
			actual := []int{1, 2}
			Expect(matcher.FailureMessage(actual)).To(SatisfyAll(
				HavePrefix("Expected"),
				HaveSuffix("to have length = 3"),
			))
		})

		It("should return the failure message for matcher comparison", func() {
			matcher := NewHaveLengthMatcher(BeNumerically(">", 5))
			actual := []int{1, 2, 3}
			Expect(matcher.FailureMessage(actual)).To(SatisfyAll(
				HavePrefix("Expected"),
				ContainSubstring("length to be >"),
				ContainSubstring("5"),
			))
		})
	})

	Context("NegatedFailureMessage", func() {
		It("should return the negated failure message for incorrect length", func() {
			matcher := NewHaveLengthMatcher(3)
			actual := []int{1, 2}
			Expect(matcher.NegatedFailureMessage(actual)).To(SatisfyAll(
				HavePrefix("Expected"),
				HaveSuffix("not to have length = 3"),
			))
		})

		It("should return the negated failure message for matcher comparison", func() {
			matcher := NewHaveLengthMatcher(BeNumerically(">", 5))
			actual := []int{1, 2, 3}
			Expect(matcher.NegatedFailureMessage(actual)).To(SatisfyAll(
				HavePrefix("Expected"),
				ContainSubstring("length not to be >"),
				ContainSubstring("5"),
			))
		})
	})
})
