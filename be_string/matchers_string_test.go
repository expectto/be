package be_string

import (
	"github.com/expectto/be/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BeStrings", func() {
	DescribeTable("should positively match", func(matcher types.BeMatcher, actual any) {
		// check gomega-compatible matching:
		success, err := matcher.Match(actual)
		Expect(err).Should(Succeed())
		Expect(success).To(BeTrue())

		// check gomock-compatible matching:
		success = matcher.Matches(actual)
		Expect(success).To(BeTrue())
	},
		Entry("NonEmptyString", NonEmptyString(), "Hello"),
		Entry("NonEmptyString: one character", NonEmptyString(), "a"),
		Entry("NonEmptyString: just space", NonEmptyString(), " "),
		Entry("EmptyString", EmptyString(), ""),

		Entry("Alpha: lowercase only", Alpha(), "abcdefg"),
		Entry("Alpha: uppercase only", Alpha(), "ABCDEFG"),
		Entry("Alpha: mixed case", Alpha(), "AbCdEfG"),

		Entry("AlphaWhitespace: only letters", AlphaWhitespace(), "HelloWorld"),
		Entry("AlphaWhitespace: only spaces", AlphaWhitespace(), "    "),
		Entry("AlphaWhitespace: mixed letters and spaces", AlphaWhitespace(), "Hello World"),
		Entry("AlphaWhitespace: letters with leading/trailing spaces", AlphaWhitespace(), "  Hello World  "),

		Entry("AlphaWithPunctuation: only letters", AlphaWithPunctuation(), "HelloWorld"),
		Entry("AlphaWithPunctuation: only punctuation", AlphaWithPunctuation(), "!().,"),
		Entry("AlphaWithPunctuation: mixed letters and punctuation", AlphaWithPunctuation(), "HelloWorld!"),

		Entry("AlphaWhitespaceWithPunctuation: only letters", AlphaWhitespaceWithPunctuation(), "HelloWorld"),
		Entry("AlphaWhitespaceWithPunctuation: only whitespace", AlphaWhitespaceWithPunctuation(), "  "),
		Entry("AlphaWhitespaceWithPunctuation: only punctuation", AlphaWhitespaceWithPunctuation(), "!().,"),
		Entry("AlphaWhitespaceWithPunctuation: mixed letters, whitespace, and punctuation", AlphaWhitespaceWithPunctuation(), "Hello, World! How are you?"),

		Entry("Whitespace: only whitespace", Whitespace(), "  "),
		Entry("Whitespace: space with tab with newline", Whitespace(), "\n\t "),

		Entry("Numeric", Numeric(), "12345"),
		Entry("Numeric: one digit", Numeric(), "1"),
		Entry("Numeric: big digits", Numeric(), "9999999999999"),

		Entry("NumericWhitespace: only numbers", NumericWhitespace(), "12345"),
		Entry("NumericWhitespace: numbers with whitespace", NumericWhitespace(), "123 45"),
		Entry("NumericWhitespace: numbers with leading and trailing whitespace", NumericWhitespace(), "  12345  "),

		Entry("AlphaNumeric: upper case", AlphaNumeric(), "ABC123"),
		Entry("AlphaNumeric: lower case", AlphaNumeric(), "abc123"),
		Entry("AlphaNumeric: mixed case", AlphaNumeric(), "ABCxyz987"),
		Entry("AlphaNumeric: only nums", AlphaNumeric(), "123456789"),
		Entry("AlphaNumeric: only alpha", AlphaNumeric(), "abcdef"),

		Entry("AlphaNumericWhitespace: alphanumeric", AlphaNumericWhitespace(), "abc123"),
		Entry("AlphaNumericWhitespace: alphanumeric with whitespace", AlphaNumericWhitespace(), "abc123 xyz"),
		Entry("AlphaNumericWhitespace: alphanumeric with leading and trailing whitespace", AlphaNumericWhitespace(), "  abc123 xyz  "),

		Entry("AlphaNumericWithPunctuation: alphanumeric with punctuation", AlphaNumericWithPunctuation(), "abc123,xyz"),
		Entry("AlphaNumericWithPunctuation: alphanumeric with leading and trailing punctuation", AlphaNumericWithPunctuation(), "(abc123.xyz)"),

		Entry("AlphaNumericWhitespaceWithPunctuation: alphanumeric with whitespace and punctuation", AlphaNumericWhitespaceWithPunctuation(), "abc 123,xyz"),
		Entry("AlphaNumericWhitespaceWithPunctuation: alphanumeric with leading and trailing punctuation and whitespace", AlphaNumericWhitespaceWithPunctuation(), "(abc 123.xyz)"),
		Entry("AlphaNumericWhitespaceWithPunctuation: alphanumeric with punctuation and whitespace", AlphaNumericWhitespaceWithPunctuation(), "abc 123, xyz"),

		Entry("AlphaNumericWithDots", AlphaNumericWithDots(), "Abc123.5"),
		Entry("AlphaNumericWithDots: only nums+dots", AlphaNumericWithDots(), "3.141592653589793"),
		Entry("AlphaNumericWithDots: multiple dots", AlphaNumericWithDots(), "a.b.c.1.2.3"),
		Entry("AlphaNumericWithDots: only dot", AlphaNumericWithDots(), "."),
		Entry("AlphaNumericWithDots: only dots", AlphaNumericWithDots(), "..."),

		Entry("Float", Float(), "3.14"),
		Entry("Float: negative", Float(), "-3.14"),
		Entry("Float: integral", Float(), "5.00"),
		Entry("Float: integral (without dot)", Float(), "5"),

		Entry("Titled", Titled(), "This Is Titled"),
		Entry("Titled:one word", Titled(), "Yo"),

		Entry("LowerCaseOnly", LowerCaseOnly(), "this is lowercase"),
		Entry("LowerCaseOnly: one character", LowerCaseOnly(), "x"),
		Entry("LowerCaseOnly: not trimmed", LowerCaseOnly(), "   hello    "),

		Entry("UpperCaseOnly", UpperCaseOnly(), "THIS IS CAPS"),
		Entry("UpperCaseOnly: one character", UpperCaseOnly(), "X"),
		Entry("UpperCaseOnly: not trimmed", UpperCaseOnly(), "   HELLO    "),

		Entry("ContainingSubstring: contains 'abc'", ContainingSubstring("lazy"), "The quick brown fox jumps over the lazy"),
		Entry("ContainingSubstring: contains '123'", ContainingSubstring("123"), "The password is 123456"),
		Entry("ContainingSubstring: contains 'xyz'", ContainingSubstring("xyz"), "xyz is the last three characters"),

		Entry("ContainingOnlyCharacters: contains only 'abc'", ContainingOnlyCharacters("abc"), "aaaaab"),
		Entry("ContainingOnlyCharacters: contains only '123'", ContainingOnlyCharacters("123"), "123"),
		Entry("ContainingOnlyCharacters: contains only 'xyz'", ContainingOnlyCharacters("xyz"), "xyzxyzxyzxyz"),

		Entry("ContainingCharacters: contains 'abc'", ContainingCharacters("abc"), "abc"),
		Entry("ContainingCharacters: contains 'abc123'", ContainingCharacters("abc123"), "1a2b3c"),
		Entry("ContainingCharacters: contains '123'", ContainingCharacters("123"), "foo111112222233331112223bar"),
		Entry("ContainingCharacters: empty chars list", ContainingCharacters(""), "anything"),
		Entry("ContainingCharacters: empty given & empty actual", ContainingCharacters(""), ""),

		Entry("MatchWildcard: prefix one char", MatchWildcard("*ello"), "Hello"),
		Entry("MatchWildcard: prefix longer", MatchWildcard("*orld"), "Hello World"),
		Entry("MatchWildcard: suffix one char", MatchWildcard("Hell*"), "Hello"),
		Entry("MatchWildcard: suffix longer", MatchWildcard("Hello W*"), "Hello World"),
		Entry("MatchWildcard: in the middle", MatchWildcard("H*d"), "Hello World"),
		Entry("MatchWildcard: all-star", MatchWildcard("*"), "Hello World"),
		Entry("MatchWildcard: all-star for nothing", MatchWildcard("*"), ""),

		Entry("ValidEmail", ValidEmail(), "test@example.com"),
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
		Entry("NonEmptyString: empty string", NonEmptyString(), ""),
		Entry("EmptyString: non-empty string", EmptyString(), "Hello"),

		Entry("Alpha: alphanumeric string", Alpha(), "abc123"),
		Entry("Alpha: string with space", Alpha(), "Hello World"),
		Entry("Alpha: string with special characters", Alpha(), "Hello@World"),
		Entry("Alpha: empty string", Alpha(), ""),

		Entry("AlphaWhitespace: contains numbers", AlphaWhitespace(), "abc123"),
		Entry("AlphaWhitespace: contains punctuation", AlphaWhitespace(), "abc!"),
		Entry("AlphaWhitespace: contains both numbers and punctuation", AlphaWhitespace(), "abc 123!"),
		Entry("AlphaWhitespace: empty string", AlphaWhitespace(), ""),

		Entry("AlphaWithPunctuation: contains whitespace", AlphaWithPunctuation(), "abc def"),
		Entry("AlphaWithPunctuation: contains numbers", AlphaWithPunctuation(), "abc123"),
		Entry("AlphaWithPunctuation: empty string", AlphaWithPunctuation(), ""),

		Entry("AlphaWhitespaceWithPunctuation: contains numbers", AlphaWhitespaceWithPunctuation(), "abc123"),
		Entry("AlphaWhitespaceWithPunctuation: empty string", AlphaWhitespaceWithPunctuation(), ""),

		Entry("Whitespace: contains letters", Whitespace(), "abc"),
		Entry("Whitespace: contains numbers", Whitespace(), "123"),
		Entry("Whitespace: empty string", Whitespace(), ""),

		Entry("NumericWhitespace: contains letters", NumericWhitespace(), "abc"),
		Entry("NumericWhitespace: contains punctuation", NumericWhitespace(), "1,2,3"),
		Entry("NumericWhitespace: contains letters and punctuation", NumericWhitespace(), "abc 123!"),
		Entry("NumericWhitespace: empty string", NumericWhitespace(), ""),

		Entry("AlphaNumericWhitespace: contains punctuation", AlphaNumericWhitespace(), "abc123!"),
		Entry("AlphaNumericWhitespace: contains special characters", AlphaNumericWhitespace(), "abc 123@"),
		Entry("AlphaNumericWhitespace: contains letters, numbers, and punctuation", AlphaNumericWhitespace(), "abc 123!"),
		Entry("AlphaNumericWhitespace: empty string", AlphaNumericWhitespace(), ""),

		Entry("AlphaNumericWithPunctuation: contains whitespace", AlphaNumericWithPunctuation(), "abc 123"),
		Entry("AlphaNumericWithPunctuation: contains whitespace and special characters", AlphaNumericWithPunctuation(), "abc 123@"),

		Entry("AlphaNumericWhitespaceWithPunctuation: contains special characters", AlphaNumericWhitespaceWithPunctuation(), "abc$% 123@"),
		Entry("AlphaNumericWhitespaceWithPunctuation: empty string", AlphaNumericWhitespaceWithPunctuation(), ""),

		Entry("Numeric: alphanumeric string", Numeric(), "abc123"),
		Entry("Numeric: string with space", Numeric(), "123 456"),
		Entry("Numeric: string with special characters", Numeric(), "123@456"),
		Entry("Numeric: empty string", Numeric(), ""),

		Entry("AlphaNumeric: string with space", AlphaNumeric(), "abc 123"),
		Entry("AlphaNumeric: string with special characters", AlphaNumeric(), "abc@123"),
		Entry("AlphaNumeric: empty string", AlphaNumeric(), ""),

		Entry("AlphaNumericWithDots: string with space", AlphaNumericWithDots(), "abc 123"),
		Entry("AlphaNumericWithDots: string with special characters", AlphaNumericWithDots(), "abc@123"),
		Entry("AlphaNumericWithDots: empty string", AlphaNumericWithDots(), ""),

		Entry("Float: string with non-numeric characters", Float(), "3.14abc"),
		Entry("Float: string with space", Float(), "3.14 5"),
		Entry("Float: string with special characters", Float(), "3.14@5"),
		Entry("Float: empty string", Float(), ""),

		Entry("Titled: non-titled string", Titled(), "hello world"),
		Entry("LowerCaseOnly: string with upper case", LowerCaseOnly(), "HelloWorld"),

		Entry("ContainingSubstring: does not contain substring", ContainingSubstring("xyz"), "abc123"),
		Entry("ContainingSubstring: empty string", ContainingSubstring("xyz"), ""),

		Entry("ContainingOnlyCharacters: contains other characters", ContainingOnlyCharacters("abc"), "defabc123"),
		Entry("ContainingOnlyCharacters: empty string", ContainingOnlyCharacters("abc"), ""),
		Entry("ContainingOnlyCharacters: contains whitespace", ContainingOnlyCharacters("abc"), "a b c"),

		Entry("ContainingCharacters: missing required characters", ContainingCharacters("abc"), "def123"),
		Entry("ContainingCharacters: empty string", ContainingCharacters("abc"), ""),
		Entry("ContainingCharacters: contains extra characters", ContainingCharacters("abc"), "def123"),

		Entry("MatchWildcard: no match", MatchWildcard("abc*"), "xyz123"),

		Entry("ValidEmail: invalid email", ValidEmail(), "test@example@com"),
		Entry("ValidEmail: just letters", ValidEmail(), "testexample"),
		Entry("ValidEmail: numeric string", ValidEmail(), "1000"),
	)

})
