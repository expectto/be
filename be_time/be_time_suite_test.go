package be_time_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBeTime(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BeTime Suite")
}
