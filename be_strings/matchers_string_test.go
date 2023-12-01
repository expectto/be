package be_strings_test

import (
	"github.com/expectto/be"
	"github.com/expectto/be/be_strings"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MatchersString", func() {
	It("should perform basic string & template matching", func() {
		Expect("Hello Jack! Your email is ask!example.com. Bye Jack").To(
			be_strings.MatchTemplate(
				`Hello {{User}}! Your email is {{Email}}. Bye {{User}}`,

				// Inside input message we should have either Jack or Jill
				be_strings.Var("User", be.Any("Jack", "Jill")),
				// any valid email with @example.com suffix
				be_strings.Var("Email", be.All(be_strings.ValidEmail(), HaveSuffix("@example.com"))),
			),
		)
	})
})
