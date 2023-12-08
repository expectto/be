package be_strings

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MatchersString", func() {
	Context("EmptyString matcher", func() {
		It("should match empty strings", func() {
			success, err := EmptyString().Match("")
			Expect(err).Should(Succeed())
			Expect(success).To(BeTrue())
		})

		It("should not match empty strings", func() {
			success, err := EmptyString().Match("non an empty string")
			Expect(err).Should(Succeed())
			Expect(success).To(BeFalse())
		})

		It("should not match non strings", func() {
			success, err := EmptyString().Match(999)
			Expect(err).Should(Succeed())
			Expect(success).To(BeFalse())
		})
	})

})
