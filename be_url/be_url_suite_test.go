package be_url_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBeUrl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BeUrl Suite")
}
