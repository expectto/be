package psi_matchers_test

import (
	. "github.com/expectto/be/internal/psi_matchers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/expectto/be/types"
)

var _ = Describe("AlwaysMatcher", func() {
	var matcher types.BeMatcher

	BeforeEach(func() {
		matcher = NewAlwaysMatcher()
	})

	Context("Matching", func() {
		DescribeTable("Always match",
			func(actual any) {
				success, err := matcher.Match(actual)
				Expect(success).To(BeTrue())
				Expect(err).NotTo(HaveOccurred())

				success = matcher.Matches(actual)
				Expect(success).To(BeTrue())

				Expect(matcher.FailureMessage(actual)).To(BeEmpty())
				Expect(matcher.NegatedFailureMessage(actual)).To(BeEmpty())
				Expect(matcher.String()).To(BeEmpty())
			},
			Entry("nil", nil),
			Entry("string", "foobar"),
			Entry("zero", 0),
			Entry("positive int", 5),
			Entry("negative int", -5),
			Entry("bool true", true),
			Entry("bool false", false),
			Entry("empty slice", []any{}),
			Entry("string slice", []string{"foo", "bar"}),
			Entry("empty map", map[string]any{}),
			Entry("filled map", map[string]any{"foo": "bar"}),
		)
	})
})
