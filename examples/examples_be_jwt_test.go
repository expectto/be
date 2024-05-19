package examples

import (
	"github.com/expectto/be/be_jwt"
	"github.com/golang-jwt/jwt/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Examples on Matching JWT", func() {
	// Here's A JWT signed with secret="my-secret"
	// with payload: {"sub":"1","name":"John Doe"}
	const tokenStr = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIiwibmFtZSI6IkpvaG4gRG9lIn0.o6m3ELBBXXiveSSfK-hdxdlbKoB3UsktDhqlt28etWk"

	// Parsed token
	var token *jwt.Token
	BeforeEach(func() {
		var err error
		token, err = jwt.Parse(tokenStr, func(_ *jwt.Token) (any, error) {
			return []byte("my-secret"), nil
		})
		Expect(err).To(Succeed())
	})

	Context("parsed token as actual value", func() {

		It("should expect parsed token to be a token", func() {
			Expect(token).To(be_jwt.Token())
		})

		It("should expect parsed token to be the token", func() {
			Expect(token).To(be_jwt.Token(tokenStr))
		})

		It("should expect parsed token to be a valid token", func() {
			Expect(token).To(be_jwt.Token(be_jwt.Valid()))
			// or simply using just Valid() directly:
			Expect(token).To(be_jwt.Valid())
		})

		It("should expect parsed token to be the token with specified details", func() {
			Expect(token).To(be_jwt.Token(
				be_jwt.Valid(),
				be_jwt.HavingClaims(
					HaveKeyWithValue("sub", "1"),
					HaveKeyWithValue("name", "John Doe"),
				),
				be_jwt.HavingMethodAlg("HS256"),
				be_jwt.SignedVia("my-secret"),
			))
		})
	})

	Context("String token as actual value", func() {

		It("should not expect a string token to be a jwt token (without transforming)", func() {
			Expect(tokenStr).NotTo(be_jwt.Token())
		})
		Context("with unsigned transform", func() {
			// Transform is required as first argument

			It("should expect a string token to be a jwt token via transform", func() {
				Expect(tokenStr).To(be_jwt.Token(be_jwt.TransformJwtFromString))
			})

			It("should not expect a string token to be the jwt token via non-signed transform", func() {
				Expect(tokenStr).NotTo(be_jwt.Token(be_jwt.TransformJwtFromString, token))
			})

			It("should expect a string token to be the jwt token via signed transform", func() {
				Expect(tokenStr).To(be_jwt.Token(be_jwt.TransformSignedJwtFromString("my-secret"), token))
			})

			It("should not expect a string token to be a valid token via unsigned transform", func() {
				Expect(tokenStr).To(be_jwt.Token(
					be_jwt.TransformJwtFromString,
					Not(be_jwt.Valid()),
				))
			})

			It("should expect string token to be a jwt token with given details matched", func() {
				Expect(tokenStr).To(be_jwt.Token(
					be_jwt.TransformJwtFromString,
					be_jwt.HavingClaims(
						HaveKeyWithValue("sub", "1"),
						HaveKeyWithValue("name", "John Doe"),
					),
					be_jwt.HavingMethodAlg("HS256"),
					be_jwt.SignedVia("my-secret"),
				))
			})
		})

		Context("signed transform", func() {
			It("should expect a string token to be valid via signed transform with proper secret", func() {
				Expect(tokenStr).To(be_jwt.Token(
					be_jwt.TransformSignedJwtFromString("my-secret"),
					be_jwt.Valid(),
				))
			})

			It("should not expect a string token to be a token via invalidly signed transform", func() {
				Expect(tokenStr).NotTo(be_jwt.Token(
					be_jwt.TransformSignedJwtFromString("invalid-secret"),
				))
			})
		})
	})

})
