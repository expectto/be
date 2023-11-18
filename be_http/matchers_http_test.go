package be_http_test

import (
	"bytes"
	"github.com/expectto/be/be_http"
	"github.com/expectto/be/be_json"
	"github.com/expectto/be/be_jwt"
	"github.com/expectto/be/be_reflected"
	"github.com/expectto/be/be_strings"
	"github.com/expectto/be/be_url"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("MatchersHttp", func() {
	It("should check a request", func() {
		// 1. Let's say we test a function that returns a *http.Request
		// req, err := SomeFunc()
		req, _ := http.NewRequest(http.MethodPost,
			"https://example.com/path?status=active&v=1&q=Hello+World",
			bytes.NewReader([]byte(`{
				"hello": "world",
				"n": 3,
				"details": [{"key":"foo"},{"key":"bar"}],
				"ids":["id1", "id2", "id3"]
			}`)),
		)
		req.Header.Set("X-Custom", "Hey-There")
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")

		// 2. Let's match everything about the request
		gomega.Expect(req).To(be_http.Request(
			be_http.HavingURL(be_url.URL(
				be_url.WithHttps(),
				be_url.HavingHost("example.com"),
				be_url.HavingPath("/path"),
				be_url.HavingSearchParam("status", "active"),
				be_url.HavingSearchParam("v", be_reflected.AsNumericString()), // any number
				be_url.HavingSearchParam("q", "Hello World"),
			)),
			be_http.HavingMethod("POST"),
			be_http.HavingBody(
				be_json.Matcher(
					be_json.JsonAsReader,
					be_json.HaveKeyValue("hello", "world"),
					be_json.HaveKeyValue("n", be_reflected.AsIntish()), // any int number
					// match ids to be sliceof strings
					// and details ot
					/*
						"hello": "world",
						"n": 3,
						"details": [{"key":"foo"},{"key":"bar"}],
						"ids":["id1", "id2", "id3"]
					*/
				),
			),

			be_http.HavingHeader("X-Custom", "Hey-There"),
			be_http.HavingHeader("Authorization"),
			be_strings.Template("Bearer {{jwt}}",
				be_strings.With("jwt",
					be_jwt.Token(
						be_jwt.BeingValid(),
						be_jwt.HavingClaims("name", "John Doe"),
					),
				),
			),
		)),
	})
})
