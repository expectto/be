package examples

import (
	"github.com/expectto/be"
	"github.com/expectto/be/be_strings"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Showcase for MatchersString", func() {
	It("should correctly match non-empty string", func() {
		Expect("Hello").To(be_strings.NonEmptyString())
		Expect("").NotTo(be_strings.NonEmptyString())
	})

	It("should correctly match empty string", func() {
		Expect("").To(be_strings.EmptyString())
		Expect("Hello").NotTo(be_strings.EmptyString())
	})

	It("should correctly match strings with only alphabets", func() {
		Expect("abcXYZ").To(be_strings.Alpha())
		Expect("123abc").NotTo(be_strings.Alpha())
	})

	It("should correctly match strings with only integer numeric values", func() {
		Expect("123").To(be_strings.Numeric())
		Expect("abc123").NotTo(be_strings.Numeric())
		Expect("125.0").NotTo(be_strings.Numeric())
	})

	It("should correctly match strings with only float numeric values", func() {
		Expect("123.0").NotTo(be_strings.Numeric())
		Expect("123").To(be_strings.Float())     // float is a numeric as well
		Expect(".5").NotTo(be_strings.Numeric()) // leading zero is ok to be ommited
	})

	It("should correctly match strings with alphanumeric values", func() {
		Expect("abc123").To(be_strings.AlphaNumeric())
		Expect("!@#").NotTo(be_strings.AlphaNumeric())
	})

	It("should correctly match strings with title case", func() {
		Expect("Hello World").To(be_strings.Titled())
		Expect("hello world").NotTo(be_strings.Titled())
	})

	It("should correctly match strings with lowercase only", func() {
		Expect("hello").To(be_strings.LowerCaseOnly())
		Expect("Hello").NotTo(be_strings.LowerCaseOnly())
	})

	It("should correctly match strings using wildcard pattern", func() {
		Expect("abc123").To(be_strings.MatchWildcard("abc*"))
		Expect("xyz").NotTo(be_strings.MatchWildcard("abc*"))
	})

	It("should correctly match valid email addresses", func() {
		Expect("user@example.com").To(be_strings.ValidEmail())
		Expect("invalid-email").NotTo(be_strings.ValidEmail())
	})

	Context("MatchTemplate", func() {
		It("should correctly match strings using a template", func() {
			Expect("Hello John! Your number is 42. Goodbye John.").To(
				be_strings.MatchTemplate(
					"Hello {{Name}}! Your number is {{Number}}. Goodbye {{Name}}.",
					be_strings.Var("Name", "John"),
					be_strings.Var("Number", be_strings.Numeric()),
				),
			)
			Expect("Invalid template").NotTo(be_strings.MatchTemplate("Hello {{Name}}. Goodbye {{Name}}."))
		})

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
	})
})
