package be_jwt_test

import (
	"github.com/expectto/be/be_jwt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Examples on Matching JWT", func() {
	It("should match JWT without secret", func() {
		// Here's A JWT signed with secret="foobar"
		// with payload: {"sub":"1","name":"John Doe"}
		const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIiwibmFtZSI6IkpvaG4gRG9lIn0.JrZ_ZcJ5GcxLsKl6c0VHfiH7RKFI-kvXCMAQ1zXJ45Q"

		// Basic usage: simply checks if given string is a jwt
		Expect(token).To(be_jwt.Token())
	})
})
