package be_string_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBeStrings(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BeString Suite")
}
