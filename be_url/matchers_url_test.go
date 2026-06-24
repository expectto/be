package be_url_test

import (
	"net/url"

	"github.com/expectto/be/be_url"
	"github.com/expectto/be/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// mustParse is a tiny helper to build *url.URL fixtures for the tables.
func mustParse(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	Expect(err).NotTo(HaveOccurred())
	return u
}

var _ = Describe("BeUrl", func() {

	DescribeTable("should positively match", func(matcher types.BeMatcher, actual any) {
		// check gomega-compatible matching:
		success, err := matcher.Match(actual)
		Expect(err).Should(Succeed())
		Expect(success).To(BeTrue())

		// check gomock-compatible matching:
		success = matcher.Matches(actual)
		Expect(success).To(BeTrue())
	},
		// URL with TransformUrlFromString (raw string -> *url.URL)
		Entry("URL transforms a raw string and matches its scheme",
			be_url.URL(be_url.TransformUrlFromString, be_url.HavingScheme("https")),
			"https://example.com/path"),
		Entry("URL transforms a raw string and matches multiple fields",
			be_url.URL(
				be_url.TransformUrlFromString,
				be_url.HavingScheme("https"),
				be_url.HavingHost("example.com"),
				be_url.HavingPath("/path"),
				be_url.HavingSearchParam("status", "active"),
			),
			"https://example.com/path?status=active"),
		Entry("URL with no args matches any valid *url.URL",
			be_url.URL(),
			mustParse("https://example.com")),

		// HavingScheme / WithHttps / WithHttp
		Entry("HavingScheme https", be_url.HavingScheme("https"), mustParse("https://example.com")),
		Entry("WithHttps", be_url.WithHttps(), mustParse("https://example.com")),
		Entry("WithHttp", be_url.WithHttp(), mustParse("http://example.com")),

		// HavingHost / HavingHostname
		Entry("HavingHost with port", be_url.HavingHost("example.com:8080"), mustParse("http://example.com:8080/x")),
		Entry("HavingHostname without port", be_url.HavingHostname("example.com"), mustParse("http://example.com:8080/x")),

		// HavingPath
		Entry("HavingPath", be_url.HavingPath("/path/to/page"), mustParse("https://example.com/path/to/page")),

		// HavingPort / NotHavingPort
		Entry("HavingPort", be_url.HavingPort("8080"), mustParse("http://example.com:8080/x")),
		Entry("NotHavingPort on a port-less url", be_url.NotHavingPort(), mustParse("https://example.com/x")),

		// HavingSearchParam
		Entry("HavingSearchParam exact value", be_url.HavingSearchParam("status", "active"),
			mustParse("https://example.com?status=active")),

		// HavingRawQuery
		Entry("HavingRawQuery", be_url.HavingRawQuery("a=1&b=2"), mustParse("https://example.com?a=1&b=2")),

		// HavingUsername / HavingPassword / HavingUserinfo
		Entry("HavingUsername", be_url.HavingUsername("user"), mustParse("https://user:pass@example.com")),
		Entry("HavingPassword", be_url.HavingPassword("pass"), mustParse("https://user:pass@example.com")),
		Entry("HavingUserinfo", be_url.HavingUserinfo("user:pass"), mustParse("https://user:pass@example.com")),
	)

	DescribeTable("should negatively match", func(matcher types.BeMatcher, actual any) {
		// check gomega-compatible matching:
		success, err := matcher.Match(actual)
		Expect(err).Should(Succeed())
		Expect(success).To(BeFalse())

		// check gomock-compatible matching:
		success = matcher.Matches(actual)
		Expect(success).To(BeFalse())
	},
		// URL with TransformUrlFromString
		Entry("URL transforms a raw string with the wrong scheme",
			be_url.URL(be_url.TransformUrlFromString, be_url.HavingScheme("https")),
			"http://example.com/path"),

		// HavingScheme / WithHttps / WithHttp
		Entry("HavingScheme wrong", be_url.HavingScheme("https"), mustParse("http://example.com")),
		Entry("WithHttps on http url", be_url.WithHttps(), mustParse("http://example.com")),
		Entry("WithHttp on https url", be_url.WithHttp(), mustParse("https://example.com")),

		// HavingHost / HavingHostname
		Entry("HavingHost wrong", be_url.HavingHost("other.com"), mustParse("http://example.com/x")),
		Entry("HavingHostname wrong", be_url.HavingHostname("other.com"), mustParse("http://example.com/x")),

		// HavingPath
		Entry("HavingPath wrong", be_url.HavingPath("/other"), mustParse("https://example.com/path")),

		// HavingPort / NotHavingPort
		Entry("HavingPort wrong", be_url.HavingPort("9090"), mustParse("http://example.com:8080/x")),
		Entry("NotHavingPort on a url that has the given port", be_url.NotHavingPort("8080"),
			mustParse("http://example.com:8080/x")),

		// HavingSearchParam
		Entry("HavingSearchParam wrong value", be_url.HavingSearchParam("status", "inactive"),
			mustParse("https://example.com?status=active")),

		// HavingRawQuery
		Entry("HavingRawQuery wrong", be_url.HavingRawQuery("a=2"), mustParse("https://example.com?a=1")),

		// HavingUsername / HavingPassword / HavingUserinfo
		Entry("HavingUsername wrong", be_url.HavingUsername("other"), mustParse("https://user:pass@example.com")),
		Entry("HavingPassword wrong", be_url.HavingPassword("other"), mustParse("https://user:pass@example.com")),
		Entry("HavingUserinfo wrong", be_url.HavingUserinfo("user:other"), mustParse("https://user:pass@example.com")),
	)

	DescribeTable("should return a valid failure message", func(matcher types.BeMatcher, actual any, message string) {
		// FailureMessage is considered to be called after matching:
		_, _ = matcher.Match(actual)

		failureMessage := matcher.FailureMessage(actual)
		Expect(failureMessage).To(Equal(message))
	},
		Entry("HavingScheme mismatch",
			be_url.HavingScheme("https"), mustParse("http://example.com"),
			"Expected\n    <string>: http\nto equal\n    <string>: https"),
		Entry("HavingPath mismatch",
			be_url.HavingPath("/expected"), mustParse("https://example.com/actual"),
			"Expected\n    <string>: /actual\nto equal\n    <string>: /expected"),
	)

	It("Values matches url.Values directly", func() {
		v := url.Values{"page": {"2"}, "sort": {"name"}}
		Expect(v).To(be_url.Values(
			be_url.HavingSearchParam("page", "2"),
			be_url.HavingSearchParam("sort", "name"),
		))
		Expect(v).NotTo(be_url.Values(be_url.HavingSearchParam("page", "9")))
	})

	It("NotHavingSearchParam distinguishes absent from present", func() {
		u := mustParse("https://example.com/?status=active")
		Expect(u).To(be_url.NotHavingSearchParam("missing"))
		Expect(u).NotTo(be_url.NotHavingSearchParam("status"))
	})
})
