package be_test

import (
	"bytes"
	"github.com/expectto/be"
	"github.com/expectto/be/be_http"
	"github.com/expectto/be/be_jwt"
	"github.com/expectto/be/be_strings"
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
			be_http.HavingHeader("Authorization",
				be.StringAsTemplate("Bearer {{jwt}}",
					be_strings.Var("jwt",
						be.JwtToken(
							be_jwt.TransformSignedJwtFromString("my-secret"),
							be_jwt.Valid(),
							be_jwt.HavingClaims("name", "John Doe"),
						),
					),
				),
			),
		))
	})
})
