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
	})
	//
	//Context("FailureMessage", func() {
	//	It("should return the failure message for incorrect length", func() {
	//		actual := []int{1, 2}
	//		failureMsg := matcher.FailureMessage(actual)
	//		expectedMsg := fmt.Sprintf("Expected\n    []int{1, 2}\nto have length = 3")
	//		Expect(failureMsg).To(Equal(expectedMsg))
	//	})
	//
	//	It("should return the failure message for matcher comparison", func() {
	//		matcher = psi_matchers.NewHaveLengthMatcher(BeNumerically(">", 5))
	//		actual := []int{1, 2, 3}
	//		failureMsg := matcher.FailureMessage(actual)
	//		expectedMsg := fmt.Sprintf("Expected\n    []int{1, 2, 3}\nlength to be > 5")
	//		Expect(failureMsg).To(Equal(expectedMsg))
	//	})
	//})
	//
	//Context("NegatedFailureMessage", func() {
	//	It("should return the negated failure message for incorrect length", func() {
	//		actual := []int{1, 2}
	//		negatedFailureMsg := matcher.NegatedFailureMessage(actual)
	//		expectedMsg := fmt.Sprintf("Expected\n    []int{1, 2}\nnot to have length = 3")
	//		Expect(negatedFailureMsg).To(Equal(expectedMsg))
	//	})
	//
	//	It("should return the negated failure message for matcher comparison", func() {
	//		matcher = psi_matchers.NewHaveLengthMatcher(BeNumerically(">", 5))
	//		actual := []int{1, 2, 3}
	//		negatedFailureMsg := matcher.NegatedFailureMessage(actual)
	//		expectedMsg := fmt.Sprintf("Expected\n    []int{1, 2, 3}\nlength not to be > 5")
	//		Expect(negatedFailureMsg).To(Equal(expectedMsg))
	//	})
	//})
})
