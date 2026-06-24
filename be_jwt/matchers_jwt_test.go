package be_jwt_test

import (
	"github.com/expectto/be/be_jwt"
	"github.com/expectto/be/types"
	"github.com/golang-jwt/jwt/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	secret      = "s3cr3t-signing-key"
	wrongSecret = "totally-different-key"
)

// mustSign builds and signs a JWT with the given method, key and claims,
// returning the signed compact string. It panics on error (test setup only).
func mustSign(method jwt.SigningMethod, key any, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(method, claims)
	signed, err := token.SignedString(key)
	if err != nil {
		panic(err)
	}
	return signed
}

var (
	// A valid HS256 token signed with `secret`.
	validHS256 = mustSign(jwt.SigningMethodHS256, []byte(secret), jwt.MapClaims{
		"name":  "John Doe",
		"admin": true,
		"sub":   "1234567890",
	})

	// A valid HS384 token signed with `secret` (for alg checks).
	validHS384 = mustSign(jwt.SigningMethodHS384, []byte(secret), jwt.MapClaims{
		"name": "Jane Roe",
	})

	// A token signed with a different secret than the one we verify with.
	signedWithWrongSecret = mustSign(jwt.SigningMethodHS256, []byte(wrongSecret), jwt.MapClaims{
		"name": "John Doe",
	})

	// A structurally-valid token whose signature byte has been tampered with.
	tamperedHS256 = tamper(validHS256)
)

// tamper flips the last character of the signature segment, producing a
// structurally-parseable but signature-invalid token.
func tamper(token string) string {
	b := []byte(token)
	last := b[len(b)-1]
	if last == 'A' {
		b[len(b)-1] = 'B'
	} else {
		b[len(b)-1] = 'A'
	}
	return string(b)
}

var _ = Describe("BeJwt", func() {

	DescribeTable("should positively match", func(matcher types.BeMatcher, actual any) {
		// gomega-compatible matching:
		success, err := matcher.Match(actual)
		Expect(err).Should(Succeed())
		Expect(success).To(BeTrue())

		// gomock-compatible matching:
		success = matcher.Matches(actual)
		Expect(success).To(BeTrue())
	},
		// Valid: token parsed+verified with the correct secret is valid.
		Entry("valid signed token is Valid()",
			be_jwt.Token(be_jwt.TransformSignedJwtFromString(secret), be_jwt.Valid()), validHS256),

		// HavingClaim with the correct value.
		Entry("token has claim name=John Doe",
			be_jwt.Token(be_jwt.TransformSignedJwtFromString(secret), be_jwt.HavingClaim("name", "John Doe")), validHS256),
		Entry("token has claim admin=true",
			be_jwt.Token(be_jwt.TransformSignedJwtFromString(secret), be_jwt.HavingClaim("admin", true)), validHS256),
		Entry("token has claim sub=1234567890",
			be_jwt.Token(be_jwt.TransformSignedJwtFromString(secret), be_jwt.HavingClaim("sub", "1234567890")), validHS256),

		// HavingClaims: match the whole claims map.
		Entry("token claims contain name key",
			be_jwt.Token(be_jwt.TransformSignedJwtFromString(secret), be_jwt.HavingClaims(HaveKeyWithValue("name", "John Doe"))), validHS256),

		// HavingMethodAlg.
		Entry("HS256 token has method alg HS256",
			be_jwt.Token(be_jwt.TransformSignedJwtFromString(secret), be_jwt.HavingMethodAlg("HS256")), validHS256),
		Entry("HS384 token has method alg HS384",
			be_jwt.Token(be_jwt.TransformSignedJwtFromString(secret), be_jwt.HavingMethodAlg("HS384")), validHS384),

		// SignedVia on a secret-less parsed token.
		Entry("token is SignedVia(secret) when parsed without verification",
			be_jwt.Token(be_jwt.TransformJwtFromString, be_jwt.SignedVia(secret)), validHS256),

		// Combining several matchers on the same token.
		Entry("valid token with correct alg and claim",
			be_jwt.Token(be_jwt.TransformSignedJwtFromString(secret),
				be_jwt.Valid(),
				be_jwt.HavingMethodAlg("HS256"),
				be_jwt.HavingClaim("name", "John Doe"),
			), validHS256),
	)

	DescribeTable("should negatively match", func(matcher types.BeMatcher, actual any) {
		// gomega-compatible matching:
		success, _ := matcher.Match(actual)
		Expect(success).To(BeFalse())

		// gomock-compatible matching:
		success = matcher.Matches(actual)
		Expect(success).To(BeFalse())
	},
		// Wrong claim value.
		Entry("claim name is not Jane Doe",
			be_jwt.Token(be_jwt.TransformSignedJwtFromString(secret), be_jwt.HavingClaim("name", "Jane Doe")), validHS256),
		Entry("claim admin is not false",
			be_jwt.Token(be_jwt.TransformSignedJwtFromString(secret), be_jwt.HavingClaim("admin", false)), validHS256),

		// Wrong algorithm.
		Entry("HS256 token does not have method alg HS384",
			be_jwt.Token(be_jwt.TransformSignedJwtFromString(secret), be_jwt.HavingMethodAlg("HS384")), validHS256),

		// SignedVia with the wrong secret on a parsed (unverified) token.
		Entry("token is not SignedVia(wrongSecret)",
			be_jwt.Token(be_jwt.TransformJwtFromString, be_jwt.SignedVia(wrongSecret)), validHS256),

		// Tampered token parsed without verification, then signature-checked.
		Entry("tampered token is not SignedVia(secret)",
			be_jwt.Token(be_jwt.TransformJwtFromString, be_jwt.SignedVia(secret)), tamperedHS256),

		// When the signed-transform fails (wrong/another secret, tampered), the
		// matcher reports a clean non-match (not a Gomega error) thanks to the
		// fallible-transform wrapper.
		Entry("verifying with the wrong secret does not match",
			be_jwt.Token(be_jwt.TransformSignedJwtFromString(wrongSecret), be_jwt.Valid()), validHS256),
		Entry("token signed with another secret does not match",
			be_jwt.Token(be_jwt.TransformSignedJwtFromString(secret), be_jwt.Valid()), signedWithWrongSecret),
		Entry("tampered token does not match",
			be_jwt.Token(be_jwt.TransformSignedJwtFromString(secret), be_jwt.Valid()), tamperedHS256),
	)

	DescribeTable("should fail gracefully (no panic) on non-token input", func(matcher types.BeMatcher, actual any) {
		Expect(func() {
			success, err := matcher.Match(actual)
			// non-*jwt.Token input is reported as a clean non-match
			Expect(err).ShouldNot(HaveOccurred())
			Expect(success).To(BeFalse())
		}).NotTo(Panic())
	},
		Entry("Valid() on a plain string", be_jwt.Valid(), "not-a-token"),
		Entry("HavingClaim on an int", be_jwt.HavingClaim("name", "John Doe"), 42),
	)

	DescribeTable("should return a valid failure message", func(matcher types.BeMatcher, actual any, message string) {
		// FailureMessage is considered to be called after matching:
		_, _ = matcher.Match(actual)
		Expect(matcher.FailureMessage(actual)).To(Equal(message))
	},
		Entry("wrong claim value failure message",
			be_jwt.HavingClaim("name", "Jane Doe"),
			parse(validHS256),
			"Expected\n    <string>: John Doe\nto equal\n    <string>: Jane Doe",
		),
		Entry("wrong method alg failure message",
			be_jwt.HavingMethodAlg("HS384"),
			parse(validHS256),
			"Expected\n    <string>: HS256\nto equal\n    <string>: HS384",
		),
	)
})

// parse is a test helper that parses a signed string into a *jwt.Token (without
// verification) so matchers can be exercised directly against *jwt.Token.
func parse(signed string) *jwt.Token {
	v := be_jwt.TransformJwtFromString(signed)
	t, ok := v.(*jwt.Token)
	if !ok {
		panic("failed to parse token in test helper")
	}
	return t
}
