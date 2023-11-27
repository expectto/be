package be_jwt_test

import (
	"fmt"
	"github.com/expectto/be/be_jwt"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("MatchersJwt", func() {
	It("...", func() {
		fmt.Println(be_jwt.BeingValid())
	})
})
