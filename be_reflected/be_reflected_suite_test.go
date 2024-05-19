package be_reflected_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBeReflected(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BeReflected Suite")
}
