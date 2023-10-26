package types

import (
	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

// BeMatcher currently stands for matcher that fits both Gomega and Gomock libraries
// todo: that's a draft yet, and a subject to be changed:
//
//	we want to be more flexible: not to have gomock as a MUST probably
//	but let to have testify be suported as well
//	Though `gomega` will probably remain as a MUST dependency
type BeMatcher interface {
	gomega.OmegaMatcher
	gomock.Matcher
}
