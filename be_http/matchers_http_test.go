package be_http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/expectto/be"
	"github.com/expectto/be/be_ctx"
	"github.com/expectto/be/be_http"
	"github.com/expectto/be/be_json"
	"github.com/expectto/be/be_url"
	"github.com/expectto/be/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// newRequest builds a *http.Request via http.NewRequest, failing the spec on error.
func newRequest(method, url string, body string) *http.Request {
	var r *http.Request
	var err error
	if body == "" {
		r, err = http.NewRequest(method, url, nil)
	} else {
		r, err = http.NewRequest(method, url, strings.NewReader(body))
	}
	Expect(err).ShouldNot(HaveOccurred())
	return r
}

// jsonRequest builds a request with a JSON body and a proper content-type header.
func jsonRequest(method, url, body string) *http.Request {
	r := newRequest(method, url, body)
	r.Header.Set("Content-Type", "application/json")
	return r
}

var _ = Describe("MatchersHttp", func() {

	DescribeTable("should positively match", func(matcher types.BeMatcher, actual any) {
		// check gomega-compatible matching:
		success, err := matcher.Match(actual)
		Expect(err).Should(Succeed())
		Expect(success).To(BeTrue())

		// check gomock-compatible matching:
		success = matcher.Matches(actual)
		Expect(success).To(BeTrue())
	},
		// Request: any valid *http.Request
		Entry("Request() matches any *http.Request",
			be_http.Request(), newRequest(http.MethodGet, "https://example.com/path", "")),
		Entry("Request(string) matches the whole URL",
			be_http.Request("https://example.com/path?status=active"),
			newRequest(http.MethodGet, "https://example.com/path?status=active", "")),

		// Method matching
		Entry("GET() matches a GET request",
			be_http.GET(), newRequest(http.MethodGet, "https://example.com", "")),
		Entry("POST() matches a POST request",
			be_http.POST(), newRequest(http.MethodPost, "https://example.com", "")),
		Entry("PUT() matches a PUT request",
			be_http.PUT(), newRequest(http.MethodPut, "https://example.com", "")),
		Entry("PATCH() matches a PATCH request",
			be_http.PATCH(), newRequest(http.MethodPatch, "https://example.com", "")),
		Entry("DELETE() matches a DELETE request",
			be_http.DELETE(), newRequest(http.MethodDelete, "https://example.com", "")),
		Entry("HavingMethod(POST) matches a POST request",
			be_http.HavingMethod(http.MethodPost), newRequest(http.MethodPost, "https://example.com", "")),
		Entry("HavingMethod(HavePrefix) matches by matcher",
			be_http.HavingMethod(HavePrefix("PO")), newRequest(http.MethodPost, "https://example.com", "")),

		// Host matching
		Entry("HavingHost matches host",
			be_http.HavingHost("example.com"), newRequest(http.MethodGet, "https://example.com/path", "")),
		Entry("HavingHost(HaveSuffix) matches host by matcher",
			be_http.HavingHost(HaveSuffix(".com")), newRequest(http.MethodGet, "https://example.com/path", "")),

		// Proto matching
		Entry("HavingProto matches HTTP/1.1 (httptest default)",
			be_http.HavingProto("HTTP/1.1"), httptest.NewRequest(http.MethodGet, "https://example.com", nil)),

		// URL matching composed with be_url matchers
		Entry("HavingURL+be_url.HavingPath matches path",
			be_http.HavingURL(be_url.HavingPath("/path")),
			newRequest(http.MethodGet, "https://example.com/path", "")),
		Entry("HavingURL+be_url.URL composes multiple be_url matchers",
			be_http.HavingURL(be_url.URL(
				be_url.WithHttps(),
				be_url.HavingHost("example.com"),
				be_url.HavingPath("/path"),
				be_url.HavingSearchParam("status", "active"),
			)),
			newRequest(http.MethodGet, "https://example.com/path?status=active", "")),

		// Header matching
		Entry("HavingHeader by key only matches presence",
			be_http.HavingHeader("Content-Type"),
			jsonRequest(http.MethodPost, "https://example.com", `{}`)),
		Entry("HavingHeader by key+value matches",
			be_http.HavingHeader("Content-Type", "application/json"),
			jsonRequest(http.MethodPost, "https://example.com", `{}`)),
		Entry("HavingHeader by key+matcher matches",
			be_http.HavingHeader("Content-Type", HavePrefix("application/")),
			jsonRequest(http.MethodPost, "https://example.com", `{}`)),

		// Body matching composed with be.JSON / be_json.HaveKeyValue
		Entry("HavingBody+be.JSON+HaveKeyValue matches body field",
			be_http.HavingBody(be.JSON(
				be_json.JsonAsReader,
				be_json.HaveKeyValue("hello", "world"),
			)),
			jsonRequest(http.MethodPost, "https://example.com", `{"hello":"world"}`)),
		Entry("HavingBody+be.JSON multiple key/values",
			be_http.HavingBody(be.JSON(
				be_json.JsonAsReader,
				be_json.HaveKeyValue("hello", "world"),
				be_json.HaveKeyValue("active", true),
			)),
			jsonRequest(http.MethodPost, "https://example.com", `{"hello":"world","active":true}`)),

		// Combined request matcher
		Entry("Request(...) composes method+url+header",
			be_http.Request(
				be_http.POST(),
				be_http.HavingURL(be_url.HavingPath("/path")),
				be_http.HavingHeader("Content-Type", "application/json"),
			),
			jsonRequest(http.MethodPost, "https://example.com/path", `{}`)),
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
		Entry("Request(string) does not match a different URL",
			be_http.Request("https://example.com/other"),
			newRequest(http.MethodGet, "https://example.com/path", "")),

		Entry("GET() does not match a POST request",
			be_http.GET(), newRequest(http.MethodPost, "https://example.com", "")),
		Entry("POST() does not match a GET request",
			be_http.POST(), newRequest(http.MethodGet, "https://example.com", "")),
		Entry("HavingMethod(POST) does not match a GET request",
			be_http.HavingMethod(http.MethodPost), newRequest(http.MethodGet, "https://example.com", "")),

		Entry("HavingHost does not match a different host",
			be_http.HavingHost("other.com"), newRequest(http.MethodGet, "https://example.com/path", "")),

		Entry("HavingURL+be_url.HavingPath does not match a different path",
			be_http.HavingURL(be_url.HavingPath("/other")),
			newRequest(http.MethodGet, "https://example.com/path", "")),

		Entry("HavingHeader does not match a missing header",
			be_http.HavingHeader("X-Absent"),
			newRequest(http.MethodGet, "https://example.com", "")),
		Entry("HavingHeader does not match a wrong header value",
			be_http.HavingHeader("Content-Type", "text/plain"),
			jsonRequest(http.MethodPost, "https://example.com", `{}`)),

		Entry("HavingBody+be.JSON does not match a wrong field value",
			be_http.HavingBody(be.JSON(
				be_json.JsonAsReader,
				be_json.HaveKeyValue("hello", "mars"),
			)),
			jsonRequest(http.MethodPost, "https://example.com", `{"hello":"world"}`)),
		Entry("HavingBody+be.JSON does not match a missing field",
			be_http.HavingBody(be.JSON(
				be_json.JsonAsReader,
				be_json.HaveKeyValue("missing"),
			)),
			jsonRequest(http.MethodPost, "https://example.com", `{"hello":"world"}`)),
	)

	DescribeTable("should fail (no match) on non-*http.Request input", func(matcher types.BeMatcher, actual any) {
		Expect(func() {
			success, _ := matcher.Match(actual)
			Expect(success).To(BeFalse())
		}).NotTo(Panic())
	},
		Entry("Request() on a string", be_http.Request(), "not-a-request"),
		Entry("GET() on a number", be_http.GET(), 42),
		Entry("HavingHost() on nil", be_http.HavingHost("x"), nil),
	)

	// Regression: HavingBody must not panic when the request carries no body
	// (req.Body == nil, e.g. a GET built with a nil body).
	It("HavingBody does not panic on a nil-body request", func() {
		req := newRequest(http.MethodGet, "https://example.com", "")
		Expect(req.Body).To(BeNil())
		Expect(func() {
			success, _ := be_http.HavingBody(
				be.JSON(be_json.JsonAsReader, be_json.HaveKeyValue("x", "y")),
			).Match(req)
			Expect(success).To(BeFalse())
		}).NotTo(Panic())
	})

	// HavingCtx matches against the request's context (req.Context()).
	It("HavingCtx matches the request context via be_ctx matchers", func() {
		type ctxKey string
		req := newRequest(http.MethodGet, "https://example.com", "")
		req = req.WithContext(context.WithValue(req.Context(), ctxKey("requestID"), "abc-123"))

		Expect(req).To(be_http.HavingCtx(be_ctx.CtxWithValue(ctxKey("requestID"), "abc-123")))
		Expect(req).NotTo(be_http.HavingCtx(be_ctx.CtxWithValue(ctxKey("requestID"), "nope")))
	})

	DescribeTable("should return a valid failure message", func(matcher types.BeMatcher, actual any, substr string) {
		// FailureMessage is considered to be called after matching:
		_, _ = matcher.Match(actual)

		failureMessage := matcher.FailureMessage(actual)
		Expect(failureMessage).To(ContainSubstring(substr))
	},
		Entry("GET() vs POST request reports the actual method",
			be_http.GET(), newRequest(http.MethodPost, "https://example.com", ""), http.MethodPost),
		Entry("HavingHost mismatch reports the expected host",
			be_http.HavingHost("other.com"), newRequest(http.MethodGet, "https://example.com", ""), "other.com"),
		Entry("HavingURL mismatch reports the expected path",
			be_http.HavingURL(be_url.HavingPath("/other")),
			newRequest(http.MethodGet, "https://example.com/path", ""), "/other"),
	)
})
