package be_math_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBeMath(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BeMath Suite")
}
