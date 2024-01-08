package be_ctx_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBeCtx(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BeCtx Suite")
}
