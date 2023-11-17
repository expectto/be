package be_http_test

import (
	"bytes"
	"github.com/expectto/be/be_http"
	"github.com/expectto/be/be_jwt"
	"github.com/expectto/be/be_url"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("MatchersHttp", func() {
	It("should check a request", func() {
		req, _ := http.NewRequest("POST", "https://example.com/path?foo=bar", bytes.NewReader([]byte("hello world")))
		req.Header.Set("X-Something", "something")
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")

		gomega.Expect(req).To(be_http.Request(
			be_http.HavingURL(be_url.URL(
				be_url.WithHttps(),
				be_url.HavingHost("example.com"),
				be_url.HavingPath("/path"),
				be_url.HavingSearchParam("foo", "bar"),
			)),
			be_http.HavingMethod("POST"),
			be_http.HavingHeader(
				"X-Something", "something",
			),
			be_http.HavingHeader(
				"Authorization",
				gomega.HavePrefix("Bearer "),
				be_jwt.Token(
					be_jwt.HavingClaims("name", "John Doe"),
					be_jwt.BeingValid(),
				),
			),
		))
	})
})
