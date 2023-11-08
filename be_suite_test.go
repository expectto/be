package be_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBe(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Be Suite")
}
