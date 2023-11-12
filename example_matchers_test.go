package be_test

import (
	"bytes"
	"github.com/expectto/be"
	"github.com/expectto/be/be_http"
	"github.com/expectto/be/be_url"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("matchers_http", func() {
	It("should match an HTTP request", func() {
		req, _ := http.NewRequest("POST", "https://example.com/path?foo=bar", bytes.NewReader([]byte("hello world")))
		req.Header.Set("X-Something", "something")
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")

		Expect(req).To(be.HttpRequest(
			be_http.RequestHavingURL(be.URL(
				be_url.UrlHavingHost("example.com"),
				be_url.UrlWithHttps(),
				be_url.UrlHavingPath("/path"),
				be_url.UrlHavingSearchParam("foo", "bar"),
			)),
			be_http.RequestHavingMethod("POST"),
			be_http.RequestHavingHeader(
				"X-Something", "something",
			),
			be_http.RequestHavingHeader(
				"Authorization",
				// TODO
				//HavePrefix("Bearer "),
				//be_jwt.Token(
				//	HaveClaims("name", "John Doe"),
				//	BeValid(),
				//),
			),
		))
	})
})
