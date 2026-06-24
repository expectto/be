package be_json_test

import (
	"io"
	"strings"

	"github.com/expectto/be/be_json"
	"github.com/expectto/be/be_math"
	"github.com/expectto/be/be_reflected"
	"github.com/expectto/be/be_string"
	"github.com/expectto/be/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// sampleJSON is a JSON document exercised across several test entries.
const sampleJSON = `{
	"name": "gopher",
	"n": 42,
	"email": "user@tests.com",
	"tags": ["a", "b", "c"],
	"nested": {"inner": "value", "count": 3},
	"items": [{"key": "foo"}, {"key": "bar"}]
}`

// NOTE: JSON numbers always decode to float64. So a value like 42 is matched by
// be_reflected.AsFloat() (NOT be_reflected.AsInteger(), which inspects reflect.Kind
// and only succeeds for Go integer kinds). See the report for details.

var _ = Describe("BeJson", func() {

	DescribeTable("should positively match (string-like input)", func(matcher types.BeMatcher, actual any) {
		// check gomega-compatible matching:
		success, err := matcher.Match(actual)
		Expect(err).Should(Succeed())
		Expect(success).To(BeTrue())

		// check gomock-compatible matching:
		success = matcher.Matches(actual)
		Expect(success).To(BeTrue())
	},
		// JsonAsString mode -- simple string value
		Entry("string value via HaveKeyValue",
			be_json.Matcher(be_json.JsonAsString, be_json.HaveKeyValue("name", "gopher")),
			sampleJSON),

		// numeric value (JSON numbers decode to float64)
		Entry("numeric value via HaveKeyValue",
			be_json.Matcher(be_json.JsonAsString, be_json.HaveKeyValue("n", float64(42))),
			sampleJSON),

		// key presence only (no value assertion)
		Entry("key presence only",
			be_json.Matcher(be_json.JsonAsString, be_json.HaveKeyValue("email")),
			sampleJSON),

		// numeric value matched by other be matchers as the value
		Entry("numeric value matched by AsFloat + GreaterThan",
			be_json.Matcher(be_json.JsonAsString,
				be_json.HaveKeyValue("n", be_reflected.AsFloat(), be_math.GreaterThan(10))),
			sampleJSON),

		// string value matched by be_string matchers
		Entry("string value matched by ValidEmail",
			be_json.Matcher(be_json.JsonAsString,
				be_json.HaveKeyValue("email", be_string.ValidEmail())),
			sampleJSON),

		Entry("string value matched by ValidEmail + ContainingSubstring",
			be_json.Matcher(be_json.JsonAsString,
				be_json.HaveKeyValue("email", be_string.ValidEmail(), be_string.ContainingSubstring("@tests.com"))),
			sampleJSON),

		Entry("string value matched by NonEmptyString",
			be_json.Matcher(be_json.JsonAsString,
				be_json.HaveKeyValue("name", be_string.NonEmptyString())),
			sampleJSON),

		// nested object value
		Entry("nested object value via HaveKeyValue",
			be_json.Matcher(be_json.JsonAsString,
				be_json.HaveKeyValue("nested", be_json.HaveKeyValue("inner", "value"))),
			sampleJSON),

		Entry("nested object with reflected/math matcher on inner key",
			be_json.Matcher(be_json.JsonAsString,
				be_json.HaveKeyValue("nested", And(
					be_reflected.AsObject(),
					be_json.HaveKeyValue("count", be_reflected.AsFloat(), be_math.GreaterThan(2)),
				))),
			sampleJSON),

		// array value
		Entry("array value via HaveKeyValue (equality)",
			be_json.Matcher(be_json.JsonAsString,
				be_json.HaveKeyValue("tags", []any{"a", "b", "c"})),
			sampleJSON),

		Entry("array value reflected as slice with length",
			be_json.Matcher(be_json.JsonAsString,
				be_json.HaveKeyValue("tags", And(be_reflected.AsSlice(), HaveLen(3)))),
			sampleJSON),

		Entry("array of objects via ContainElement",
			be_json.Matcher(be_json.JsonAsString,
				be_json.HaveKeyValue("items", ContainElement(be_json.HaveKeyValue("key", "foo")))),
			sampleJSON),

		// Not(HaveKeyValue) for an absent key
		Entry("Not for absent key",
			be_json.Matcher(be_json.JsonAsString, Not(be_json.HaveKeyValue("missing_field"))),
			sampleJSON),

		// whole-document equality
		Entry("whole document equality via map",
			be_json.Matcher(be_json.JsonAsString, map[string]any{"foo": "bar"}),
			`{"foo":"bar"}`),

		// validity-only (no value args)
		Entry("string is valid json (validity-only)",
			be_json.Matcher(be_json.JsonAsString),
			`{"any":"json"}`),
	)

	DescribeTable("should negatively match (string-like input)", func(matcher types.BeMatcher, actual any) {
		// check gomega-compatible matching:
		success, err := matcher.Match(actual)
		Expect(err).Should(Succeed())
		Expect(success).To(BeFalse())

		// check gomock-compatible matching:
		success = matcher.Matches(actual)
		Expect(success).To(BeFalse())
	},
		Entry("wrong string value",
			be_json.Matcher(be_json.JsonAsString, be_json.HaveKeyValue("name", "not-gopher")),
			sampleJSON),

		Entry("wrong numeric value",
			be_json.Matcher(be_json.JsonAsString, be_json.HaveKeyValue("n", float64(7))),
			sampleJSON),

		Entry("numeric value not GreaterThan 100",
			be_json.Matcher(be_json.JsonAsString,
				be_json.HaveKeyValue("n", be_reflected.AsFloat(), be_math.GreaterThan(100))),
			sampleJSON),

		Entry("absent key requested with HaveKeyValue",
			be_json.Matcher(be_json.JsonAsString, be_json.HaveKeyValue("deleted_field")),
			sampleJSON),

		Entry("Not(HaveKeyValue) on present key fails",
			be_json.Matcher(be_json.JsonAsString, Not(be_json.HaveKeyValue("name"))),
			sampleJSON),

		Entry("nested object wrong inner value",
			be_json.Matcher(be_json.JsonAsString,
				be_json.HaveKeyValue("nested", be_json.HaveKeyValue("inner", "wrong"))),
			sampleJSON),

		Entry("array of objects missing element",
			be_json.Matcher(be_json.JsonAsString,
				be_json.HaveKeyValue("items", ContainElement(be_json.HaveKeyValue("key", "absent")))),
			sampleJSON),
	)

	// Reader-mode is tested with a fresh io.Reader factory per call: an io.Reader is
	// single-use, so it cannot be matched twice (Match then Matches). We only call
	// Match once here, recreating the reader for each entry via the factory.
	DescribeTable("should match (io.Reader input)", func(matcher types.BeMatcher, newReader func() io.Reader, expected bool) {
		success, err := matcher.Match(newReader())
		Expect(err).Should(Succeed())
		Expect(success).To(Equal(expected))
	},
		Entry("reader mode -- simple string value matches",
			be_json.Matcher(be_json.JsonAsReader, be_json.HaveKeyValue("name", "gopher")),
			func() io.Reader { return strings.NewReader(sampleJSON) }, true),

		Entry("reader mode -- numeric value matched by AsFloat",
			be_json.Matcher(be_json.JsonAsReader, be_json.HaveKeyValue("n", be_reflected.AsFloat())),
			func() io.Reader { return strings.NewReader(sampleJSON) }, true),

		Entry("reader mode -- nested object matches",
			be_json.Matcher(be_json.JsonAsReader, be_json.HaveKeyValue("nested", be_json.HaveKeyValue("inner", "value"))),
			func() io.Reader { return strings.NewReader(sampleJSON) }, true),

		Entry("reader mode -- wrong value does not match",
			be_json.Matcher(be_json.JsonAsReader, be_json.HaveKeyValue("name", "wrong")),
			func() io.Reader { return strings.NewReader(sampleJSON) }, false),
	)

	// On input that can't be parsed as JSON the matcher surfaces the parse error
	// (v1 contract: un-evaluatable input -> error, not a silent non-match) and must
	// never panic.
	DescribeTable("should error (no panic) on invalid json input", func(matcher types.BeMatcher, actual any) {
		Expect(func() {
			success, err := matcher.Match(actual)
			Expect(err).To(HaveOccurred())
			Expect(success).To(BeFalse())
		}).NotTo(Panic())
	},
		Entry("malformed json string",
			be_json.Matcher(be_json.JsonAsString, be_json.HaveKeyValue("name", "gopher")),
			`{"name": `),
		Entry("non-json string",
			be_json.Matcher(be_json.JsonAsString, be_json.HaveKeyValue("name", "gopher")),
			`this is not json`),
	)

	DescribeTable("should return a valid failure message", func(matcher types.BeMatcher, actual any, substr string) {
		// FailureMessage is considered to be called after matching:
		_, _ = matcher.Match(actual)

		failureMessage := matcher.FailureMessage(actual)
		Expect(failureMessage).To(ContainSubstring(substr))
	},
		Entry("wrong string value mentions expected key",
			be_json.Matcher(be_json.JsonAsString, be_json.HaveKeyValue("name", "not-gopher")),
			sampleJSON, "name"),

		Entry("absent key failure mentions the key",
			be_json.Matcher(be_json.JsonAsString, be_json.HaveKeyValue("deleted_field")),
			sampleJSON, "deleted_field"),
	)
})
