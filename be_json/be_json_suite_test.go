package be_json_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBeJson(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BeJson Suite")
}
