package psi_matchers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPsiMatchers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PsiMatchers Suite")
}
