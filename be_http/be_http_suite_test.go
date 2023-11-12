package be_http_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBeHttp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BeHttp Suite")
}
