package be_jwt_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBeJwt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BeJwt Suite")
}
