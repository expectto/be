package examples

import (
	"github.com/expectto/be"
	"github.com/expectto/be/be_string"
	. "github.com/expectto/be/options"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Showcase for MatchersString", func() {
	It("should correctly match non-empty string", func() {
		Expect("Hello").To(be_string.NonEmptyString())
		Expect("").NotTo(be_string.NonEmptyString())
	})

	It("should correctly match empty string", func() {
		Expect("").To(be_string.EmptyString())
		Expect("Hello").NotTo(be_string.EmptyString())
	})

	It("should correctly match strings with only alphabets", func() {
		Expect("abcXYZ").To(be_string.Only(Alpha))
		Expect("123abc").NotTo(be_string.Only(Alpha))
	})

	It("should correctly match strings with only integer numeric values", func() {
		Expect("123").To(be_string.Only(Numeric))
		Expect("abc123").NotTo(be_string.Only(Numeric))
		Expect("125.0").NotTo(be_string.Only(Numeric))
	})

	It("should correctly match strings with only float numeric values", func() {
		Expect("123.0").NotTo(be_string.Only(Numeric))
		Expect("123").To(be_string.Float())         // float is a numeric as well
		Expect(".5").NotTo(be_string.Only(Numeric)) // leading zero is ok to be ommited
	})

	It("should correctly match strings with alphanumeric values", func() {
		Expect("abc123").To(be_string.Only(Alpha | Numeric))
		Expect("!@#").NotTo(be_string.Only(Alpha | Numeric))
	})

	It("should correctly match strings with title case", func() {
		Expect("Hello World").To(be_string.Titled())
		Expect("hello world").NotTo(be_string.Titled())
	})

	It("should correctly match strings with lowercase only", func() {
		Expect("hello").To(be_string.LowerCaseOnly())
		Expect("Hello").NotTo(be_string.LowerCaseOnly())
	})

	It("should correctly match strings using wildcard pattern", func() {
		Expect("abc123").To(be_string.MatchWildcard("abc*"))
		Expect("xyz").NotTo(be_string.MatchWildcard("abc*"))
	})

	It("should correctly match valid email addresses", func() {
		Expect("user@example.com").To(be_string.ValidEmail())
		Expect("invalid-email").NotTo(be_string.ValidEmail())
	})

	Context("MatchTemplate", func() {
		It("should correctly match strings using a template", func() {
			Expect("Hello John! Your number is 42. Goodbye John.").To(
				be_string.MatchTemplate(
					"Hello {{Name}}! Your number is {{Number}}. Goodbye {{Name}}.",
					be_string.Var("Name", "John"),
					be_string.Var("Number", be_string.Only(Numeric)),
				),
			)
			Expect("Invalid template").NotTo(be_string.MatchTemplate("Hello {{Name}}. Goodbye {{Name}}."))
		})

		It("should perform basic string & template matching", func() {
			Expect("Hello Jack! Your email is ask@example.com. Bye Jack").To(
				be_string.MatchTemplate(
					`Hello {{User}}! Your email is {{Email}}. Bye {{User}}`,

					// Inside input message we should have either Jack or Jill
					be_string.Var("User", be.Any("Jack", "Jill")),
					// any valid email with @example.com suffix
					be_string.Var("Email", be.All(be_string.ValidEmail(), HaveSuffix("@example.com"))),
				),
			)
		})
	})
})
