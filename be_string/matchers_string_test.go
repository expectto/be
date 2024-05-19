package be_string_test

import (
	"github.com/expectto/be/be_string"
	. "github.com/expectto/be/options"
	"github.com/expectto/be/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/rand"
)

var _ = Describe("BeStrings (simple matchers)", func() {
	DescribeTable("should positively match", func(matcher types.BeMatcher, actual any) {
		// check gomega-compatible matching:
		success, err := matcher.Match(actual)
		Expect(err).Should(Succeed())
		Expect(success).To(BeTrue())

		// check gomock-compatible matching:
		success = matcher.Matches(actual)
		Expect(success).To(BeTrue())
	},
		Entry("NonEmptyString", be_string.NonEmptyString(), "Hello"),
		Entry("NonEmptyString: one character", be_string.NonEmptyString(), "a"),
		Entry("NonEmptyString: just space", be_string.NonEmptyString(), " "),
		Entry("EmptyString", be_string.EmptyString(), ""),

		Entry("Only(Alpha) lowercase only", be_string.Only(Alpha), "abcdefg"),
		Entry("Only(Alpha) uppercase only", be_string.Only(Alpha), "ABCDEFG"),
		Entry("Only(Alpha) mixed case", be_string.Only(Alpha), "AbCdEfG"),

		Entry("Only(Alpha|Whitespace) only letters", be_string.Only(Alpha|Whitespace), "HelloWorld"),
		Entry("Only(Alpha|Whitespace) only spaces", be_string.Only(Alpha|Whitespace), "    "),
		Entry("Only(Alpha|Whitespace) mixed letters and spaces", be_string.Only(Alpha|Whitespace), "Hello World"),
		Entry("Only(Alpha|Whitespace) letters with leading/trailing spaces", be_string.Only(Alpha|Whitespace), "  Hello World  "),

		Entry("Only(Alpha|Punctuation): only letters", be_string.Only(Alpha|Punctuation), "HelloWorld"),
		Entry("Only(Alpha|Punctuation): only punctuation", be_string.Only(Alpha|Punctuation), "!().,"),
		Entry("Only(Alpha|Punctuation): mixed letters and punctuation", be_string.Only(Alpha|Punctuation), "HelloWorld!"),

		Entry("Only(Alpha|Whitespace|Punctuation) only letters", be_string.Only(Alpha|Whitespace|Punctuation), "HelloWorld"),
		Entry("Only(Alpha|Whitespace|Punctuation) only whitespace", be_string.Only(Alpha|Whitespace|Punctuation), "  "),
		Entry("Only(Alpha|Whitespace|Punctuation) only punctuation", be_string.Only(Alpha|Whitespace|Punctuation), "!().,"),
		Entry("Only(Alpha|Whitespace|Punctuation) mixed letters, whitespace, and punctuation", be_string.Only(Alpha|Whitespace|Punctuation), "Hello, World! How are you?"),

		Entry("Only(Whitespace): only whitespace", be_string.Only(Whitespace), "  "),
		Entry("Only(Whitespace): space with tab with newline", be_string.Only(Whitespace), "\n\t "),

		Entry("Only(Numeric)", be_string.Only(Numeric), "12345"),
		Entry("Only(Numeric): one digit", be_string.Only(Numeric), "1"),
		Entry("Only(Numeric): big digits", be_string.Only(Numeric), "9999999999999"),

		Entry("Only(Numeric|Whitespace): only numbers", be_string.Only(Numeric|Whitespace), "12345"),
		Entry("Only(Numeric|Whitespace): numbers with whitespace", be_string.Only(Numeric|Whitespace), "123 45"),
		Entry("Only(Numeric|Whitespace): numbers with leading and trailing whitespace", be_string.Only(Numeric|Whitespace), "  12345  "),

		Entry("Only(Alpha|Numeric): upper case", be_string.Only(Alpha|Numeric), "ABC123"),
		Entry("Only(Alpha|Numeric): lower case", be_string.Only(Alpha|Numeric), "abc123"),
		Entry("Only(Alpha|Numeric): mixed case", be_string.Only(Alpha|Numeric), "ABCxyz987"),
		Entry("Only(Alpha|Numeric): only nums", be_string.Only(Alpha|Numeric), "123456789"),
		Entry("Only(Alpha|Numeric): only alpha", be_string.Only(Alpha|Numeric), "abcdef"),

		Entry("Only(Alpha|Numeric|Whitespace): alphanumeric", be_string.Only(Alpha|Numeric|Whitespace), "abc123"),
		Entry("Only(Alpha|Numeric|Whitespace): alphanumeric with whitespace", be_string.Only(Alpha|Numeric|Whitespace), "abc123 xyz"),
		Entry("Only(Alpha|Numeric|Whitespace): alphanumeric with leading and trailing whitespace", be_string.Only(Alpha|Numeric|Whitespace), "  abc123 xyz  "),

		Entry("Only(Alpha|Numeric|Punctuation): alphanumeric with punctuation", be_string.Only(Alpha|Numeric|Punctuation), "abc123,xyz"),
		Entry("Only(Alpha|Numeric|Punctuation): alphanumeric with leading and trailing punctuation", be_string.Only(Alpha|Numeric|Punctuation), "(abc123.xyz)"),

		Entry("Only(Alpha|Numeric|Whitespace|Punctuation): alphanumeric with whitespace and punctuation", be_string.Only(Alpha|Numeric|Whitespace|Punctuation), "abc 123,xyz"),
		Entry("Only(Alpha|Numeric|Whitespace|Punctuation): alphanumeric with leading and trailing punctuation and whitespace", be_string.Only(Alpha|Numeric|Whitespace|Punctuation), "(abc 123.xyz)"),
		Entry("Only(Alpha|Numeric|Whitespace|Punctuation): alphanumeric with punctuation and whitespace", be_string.Only(Alpha|Numeric|Whitespace|Punctuation), "abc 123, xyz"),

		Entry("Only(Alpha|Numeric|Dots)", be_string.Only(Alpha|Numeric|Dots), "Abc123.5"),
		Entry("Only(Alpha|Numeric|Dots): only nums+dots", be_string.Only(Alpha|Numeric|Dots), "3.141592653589793"),
		Entry("Only(Alpha|Numeric|Dots): multiple dots", be_string.Only(Alpha|Numeric|Dots), "a.b.c.1.2.3"),
		Entry("Only(Alpha|Numeric|Dots): only dot", be_string.Only(Alpha|Numeric|Dots), "."),
		Entry("Only(Alpha|Numeric|Dots): only dots", be_string.Only(Alpha|Numeric|Dots), "..."),

		Entry("Float", be_string.Float(), "3.14"),
		Entry("Float: negative", be_string.Float(), "-3.14"),
		Entry("Float: integral", be_string.Float(), "5.00"),
		Entry("Float: integral (without dot)", be_string.Float(), "5"),

		Entry("Titled", be_string.Titled(), "This Is Titled"),
		Entry("Titled:one word", be_string.Titled(), "Yo"),

		Entry("LowerCaseOnly", be_string.LowerCaseOnly(), "this is lowercase"),
		Entry("LowerCaseOnly: one character", be_string.LowerCaseOnly(), "x"),
		Entry("LowerCaseOnly: not trimmed", be_string.LowerCaseOnly(), "   hello    "),

		Entry("UpperCaseOnly", be_string.UpperCaseOnly(), "THIS IS CAPS"),
		Entry("UpperCaseOnly: one character", be_string.UpperCaseOnly(), "X"),
		Entry("UpperCaseOnly: not trimmed", be_string.UpperCaseOnly(), "   HELLO    "),

		Entry("ContainingSubstring: contains 'abc'", be_string.ContainingSubstring("lazy"), "The quick brown fox jumps over the lazy"),
		Entry("ContainingSubstring: contains '123'", be_string.ContainingSubstring("123"), "The password is 123456"),
		Entry("ContainingSubstring: contains 'xyz'", be_string.ContainingSubstring("xyz"), "xyz is the last three characters"),

		Entry("ContainingOnlyCharacters: contains only 'abc'", be_string.ContainingOnlyCharacters("abc"), "aaaaab"),
		Entry("ContainingOnlyCharacters: contains only '123'", be_string.ContainingOnlyCharacters("123"), "123"),
		Entry("ContainingOnlyCharacters: contains only 'xyz'", be_string.ContainingOnlyCharacters("xyz"), "xyzxyzxyzxyz"),

		Entry("ContainingCharacters: contains 'abc'", be_string.ContainingCharacters("abc"), "abc"),
		Entry("ContainingCharacters: contains 'abc123'", be_string.ContainingCharacters("abc123"), "1a2b3c"),
		Entry("ContainingCharacters: contains '123'", be_string.ContainingCharacters("123"), "foo111112222233331112223bar"),
		Entry("ContainingCharacters: empty chars list", be_string.ContainingCharacters(""), "anything"),
		Entry("ContainingCharacters: empty given & empty actual", be_string.ContainingCharacters(""), ""),

		Entry("MatchWildcard: prefix one char", be_string.MatchWildcard("*ello"), "Hello"),
		Entry("MatchWildcard: prefix longer", be_string.MatchWildcard("*orld"), "Hello World"),
		Entry("MatchWildcard: suffix one char", be_string.MatchWildcard("Hell*"), "Hello"),
		Entry("MatchWildcard: suffix longer", be_string.MatchWildcard("Hello W*"), "Hello World"),
		Entry("MatchWildcard: in the middle", be_string.MatchWildcard("H*d"), "Hello World"),
		Entry("MatchWildcard: all-star", be_string.MatchWildcard("*"), "Hello World"),
		Entry("MatchWildcard: all-star for nothing", be_string.MatchWildcard("*"), ""),

		Entry("ValidEmail", be_string.ValidEmail(), "test@example.com"),
	)

	DescribeTable("should negatively match", func(matcher types.BeMatcher, actual any) {
		// Check gomega-compatible matching:
		success, err := matcher.Match(actual)
		Expect(err).Should(Succeed())
		Expect(success).To(BeFalse())

		// Check gomock-compatible matching:
		success = matcher.Matches(actual)
		Expect(success).To(BeFalse())
	},
		Entry("NonEmptyString: empty string", be_string.NonEmptyString(), ""),
		Entry("EmptyString: non-empty string", be_string.EmptyString(), "Hello"),

		Entry("Only(Alpha): alphanumeric string", be_string.Only(Alpha), "abc123"),
		Entry("Only(Alpha): string with space", be_string.Only(Alpha), "Hello World"),
		Entry("Only(Alpha): string with special characters", be_string.Only(Alpha), "Hello@World"),
		Entry("Only(Alpha): empty string", be_string.Only(Alpha), ""),

		Entry("Only(Alpha|Whitespace): contains numbers", be_string.Only(Alpha|Whitespace), "abc123"),
		Entry("Only(Alpha|Whitespace): contains punctuation", be_string.Only(Alpha|Whitespace), "abc!"),
		Entry("Only(Alpha|Whitespace): contains both numbers and punctuation", be_string.Only(Alpha|Whitespace), "abc 123!"),
		Entry("Only(Alpha|Whitespace): empty string", be_string.Only(Alpha|Whitespace), ""),

		Entry("Only(Alpha|Punctuation): contains whitespace", be_string.Only(Alpha|Punctuation), "abc def"),
		Entry("Only(Alpha|Punctuation): contains numbers", be_string.Only(Alpha|Punctuation), "abc123"),
		Entry("Only(Alpha|Punctuation): empty string", be_string.Only(Alpha|Punctuation), ""),

		Entry("Only(Alpha|Whitespace|Punctuation) contains numbers", be_string.Only(Alpha|Whitespace|Punctuation), "abc123"),
		Entry("Only(Alpha|Whitespace|Punctuation) empty string", be_string.Only(Alpha|Whitespace|Punctuation), ""),

		Entry("Only(Whitespace): contains letters", be_string.Only(Whitespace), "abc"),
		Entry("Only(Whitespace): contains numbers", be_string.Only(Whitespace), "123"),
		Entry("Only(Whitespace): empty string", be_string.Only(Whitespace), ""),

		Entry("Only(Numeric|Whitespace): contains letters", be_string.Only(Numeric|Whitespace), "abc"),
		Entry("Only(Numeric|Whitespace): contains punctuation", be_string.Only(Numeric|Whitespace), "1,2,3"),
		Entry("Only(Numeric|Whitespace): contains letters and punctuation", be_string.Only(Numeric|Whitespace), "abc 123!"),
		Entry("Only(Numeric|Whitespace): empty string", be_string.Only(Numeric|Whitespace), ""),

		Entry("Only(Alpha|Numeric|Whitespace): contains punctuation", be_string.Only(Alpha|Numeric|Whitespace), "abc123!"),
		Entry("Only(Alpha|Numeric|Whitespace): contains special characters", be_string.Only(Alpha|Numeric|Whitespace), "abc 123@"),
		Entry("Only(Alpha|Numeric|Whitespace): contains letters, numbers, and punctuation", be_string.Only(Alpha|Numeric|Whitespace), "abc 123!"),
		Entry("Only(Alpha|Numeric|Whitespace): empty string", be_string.Only(Alpha|Numeric|Whitespace), ""),

		Entry("Only(Alpha|Numeric|Punctuation): contains whitespace", be_string.Only(Alpha|Numeric|Punctuation), "abc 123"),
		Entry("Only(Alpha|Numeric|Punctuation): contains whitespace and special characters", be_string.Only(Alpha|Numeric|Punctuation), "abc 123@"),

		Entry("Only(Alpha|Numeric|Whitespace|Punctuation): contains special characters", be_string.Only(Alpha|Numeric|Whitespace|Punctuation), "abc$% 123@"),
		Entry("Only(Alpha|Numeric|Whitespace|Punctuation): empty string", be_string.Only(Alpha|Numeric|Whitespace|Punctuation), ""),

		Entry("Only(Numeric): alphanumeric string", be_string.Only(Numeric), "abc123"),
		Entry("Only(Numeric): string with space", be_string.Only(Numeric), "123 456"),
		Entry("Only(Numeric): string with special characters", be_string.Only(Numeric), "123@456"),
		Entry("Only(Numeric): empty string", be_string.Only(Numeric), ""),

		Entry("Only(Alpha|Numeric): string with space", be_string.Only(Alpha|Numeric), "abc 123"),
		Entry("Only(Alpha|Numeric): string with special characters", be_string.Only(Alpha|Numeric), "abc@123"),
		Entry("Only(Alpha|Numeric): empty string", be_string.Only(Alpha|Numeric), ""),

		Entry("Only(Alpha|Numeric|Dots): string with space", be_string.Only(Alpha|Numeric|Dots), "abc 123"),
		Entry("Only(Alpha|Numeric|Dots): string with special characters", be_string.Only(Alpha|Numeric|Dots), "abc@123"),
		Entry("Only(Alpha|Numeric|Dots): empty string", be_string.Only(Alpha|Numeric|Dots), ""),

		Entry("Float: string with non-numeric characters", be_string.Float(), "3.14abc"),
		Entry("Float: string with space", be_string.Float(), "3.14 5"),
		Entry("Float: string with special characters", be_string.Float(), "3.14@5"),
		Entry("Float: empty string", be_string.Float(), ""),

		Entry("Titled: non-titled string", be_string.Titled(), "hello world"),
		Entry("LowerCaseOnly: string with upper case", be_string.LowerCaseOnly(), "HelloWorld"),

		Entry("ContainingSubstring: does not contain substring", be_string.ContainingSubstring("xyz"), "abc123"),
		Entry("ContainingSubstring: empty string", be_string.ContainingSubstring("xyz"), ""),

		Entry("ContainingOnlyCharacters: contains other characters", be_string.ContainingOnlyCharacters("abc"), "defabc123"),
		Entry("ContainingOnlyCharacters: empty string", be_string.ContainingOnlyCharacters("abc"), ""),
		Entry("ContainingOnlyCharacters: contains whitespace", be_string.ContainingOnlyCharacters("abc"), "a b c"),

		Entry("ContainingCharacters: missing required characters", be_string.ContainingCharacters("abc"), "def123"),
		Entry("ContainingCharacters: empty string", be_string.ContainingCharacters("abc"), ""),
		Entry("ContainingCharacters: contains extra characters", be_string.ContainingCharacters("abc"), "def123"),

		Entry("MatchWildcard: no match", be_string.MatchWildcard("abc*"), "xyz123"),

		Entry("ValidEmail: invalid email", be_string.ValidEmail(), "test@example@com"),
		Entry("ValidEmail: just letters", be_string.ValidEmail(), "testexample"),
		Entry("ValidEmail: numeric string", be_string.ValidEmail(), "1000"),
	)

	// All be_string matchers expects input to be a string.
	// They will not succeed and return a short "to be type of string" failure message
	DescribeTable("non-string type tests", func(matcher types.BeMatcher) {
		notStrings := []any{
			0, false, map[string]any{}, []string{}, func() {}, nil,
		}
		actual := notStrings[rand.Intn(len(notStrings))] // not a string

		// Check gomega-compatible matching:
		success, err := matcher.Match(actual)
		Expect(err).Should(Succeed())
		Expect(success).To(BeFalse())

		// Check gomock-compatible matching:
		success = matcher.Matches(actual)
		Expect(success).To(BeFalse())

		// and the message should be the same
		Expect(matcher.FailureMessage(actual)).To(HaveSuffix("to be type of string"))
		//Expect(matcher.NegatedFailureMessage(actual)).To(HaveSuffix("not to be type of string"))
	},
		// todo: at some point Entries should be auto-generated
		Entry("NonEmptyString", be_string.NonEmptyString()),
		Entry("EmptyString", be_string.EmptyString()),
		Entry("Only(Alpha)", be_string.Only(Alpha)),
		Entry("Float", be_string.Float()),
		Entry("Titled", be_string.Titled()),
		Entry("LowerCaseOnly", be_string.LowerCaseOnly()),
		Entry("ContainingSubstring", be_string.ContainingSubstring("xyz")),
		Entry("ContainingOnlyCharacters", be_string.ContainingOnlyCharacters("abc")),
		Entry("ContainingCharacters", be_string.ContainingCharacters("abc")),
		Entry("MatchWildcard", be_string.MatchWildcard("abc*")),
		Entry("ValidEmail", be_string.ValidEmail()),
	)

})

var _ = Describe("BeStrings (template matching)", func() {
	It("Should match a template with 2 variables", func() {
		matcher := be_string.MatchTemplate(
			"Hello {{UserName}}! Given email is {{Email}}",
			be_string.V("UserName", be_string.Only(Alpha|Numeric)),
			be_string.V("Email", be_string.ValidEmail()),
		)
		input := "Hello Foo123! Given email is hello@gmail.com"
		Expect(input).To(matcher)
	})
})
