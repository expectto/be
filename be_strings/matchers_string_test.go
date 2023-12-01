package be_strings_test

import (
	"github.com/expectto/be"
	"github.com/expectto/be/be_strings"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MatchersString", func() {
	It("should perform basic string & template matching", func() {
		Expect("Hello Jack! Your email is ask@example.com. Bye Jack").To(
			be_strings.MatchTemplate(
				`Hello {{User}}! Your email is {{Email}}. Bye {{User}}`,

				// Inside input message we should have either Jack or Jill
				be_strings.Var("User", be.Any("Jack", "Jill")),
				// any valid email with @example.com suffix
				be_strings.Var("Email", be.All(be_strings.ValidEmail(), HaveSuffix("@example.com"))),
			),
		)
	})

	It("should perform matchin on a sparse template", func() {
		Expect("Hello Jack! Your email is ask@example.com. Bye Jack").To(
			be_strings.MatchTemplate(
				// {{...}} means it can match ANYTHING (but not empty string)
				// {{..?}} means it can match ANYTHING (even empty string)
				// so here we actually simply check that email is somewhere in the middle of the string
				`{{...}}{{Email}}{{..?}}`,

				// any valid email with @example.com suffix
				be_strings.Var("Email", be.All(be_strings.ValidEmail(), HaveSuffix("@example.com"))),
			),
		)
	})
})
