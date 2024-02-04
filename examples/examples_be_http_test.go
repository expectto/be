package examples

import (
	"bytes"
	"github.com/expectto/be"
	"github.com/expectto/be/be_http"
	"github.com/expectto/be/be_json"
	"github.com/expectto/be/be_reflected"
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
			be_http.HavingURL(be.URL(
				be_url.HavingHost("example.com"),
				be_url.WithHttps(),
				be_url.HavingPath("/path"),
				be_url.HavingSearchParam("foo", "bar"),
			)),
			be_http.POST(),
			be_http.HavingHeader(
				"X-Something", "something",
			),
			be_http.HavingHeader("Authorization"), //be.StringAsTemplate("Bearer {{jwt}}",
			//	be_strings.Var("jwt",
			//		be.JwtToken(
			//			be_jwt.TransformSignedJwtFromString("my-secret"),
			//			be_jwt.Valid(),
			//			be_jwt.HavingClaims("name", "John Doe"),
			//		),
			//	),
			//),

		))
	})

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
		Expect(req).To(be_http.Request(
			// 2.1. Match the URL
			be_http.HavingURL(be_url.URL(
				be_url.WithHttps(),
				be_url.HavingHost("example.com"),
				be_url.HavingPath("/path"),
				be_url.HavingSearchParam("status", "active"),
				be_url.HavingSearchParam("v", be_reflected.AsNumericString()), // any number
				be_url.HavingSearchParam("q", "Hello World"),
			)),

			be_http.HavingMethod("POST"),

			// 2.2. Match the body
			be_http.HavingBody(
				be_json.Matcher(
					be_json.JsonAsReader,
					be_json.HaveKeyValue("hello", "world"),
					//be_json.HaveKeyValue("n", be_reflected.AsIntish()), // any int number
					//be_json.HaveKeyValue("ids", be_reflected.AsSliceOf[string]),
					be_json.HaveKeyValue("details", And(
						be_reflected.AsObjects(),
						be.HaveLength(2),
						ContainElements(
							be_json.HaveKeyValue("key", "foo"),
							be_json.HaveKeyValue("key", "bar"),
						),
					)),
				),
			),

			// 2.3. Matching the headers

			//be_http.HavingHeader("X-Custom", "Hey-There"),
			be_http.HavingHeader("Authorization", HavePrefix("Bearer ")),
			//be_strings.MatchTemplate("Bearer {{jwt}}",
			//	be_strings.Var("jwt",
			//		be_jwt.Token(
			//			be_jwt.TransformJwtFromString,
			//			be_jwt.Valid(),
			//			be_jwt.HavingClaims("name", "John Doe"),
			//		),
			//	),
			//),

			// todo: add example with Time in header, so we can test be_time
		))
	})
})
