package psi

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/expectto/be/types"
)

func FromGomega(omega types.GomegaMatcher, messagePrefixArg ...string) types.BeMatcher {
	return &upgradedOmegaMatcher{
		GomegaMatcher:      omega,
		MixinMatcherGomock: NewMixinMatcherGomock(omega, messagePrefixArg...),
	}
}

// upgradedOmegaMatcher wraps GomegaMatcher and GomockMatcher
// Upgrade "Gomega => Psi" is done via attaching MixinMatcherGomock
type upgradedOmegaMatcher struct {
	types.GomegaMatcher
	*MixinMatcherGomock
}

var (
	// ExpectedStr2LinesRegex matches first 2 lines of standard gomega failure message
	ExpectedStr2LinesRegex = regexp.MustCompile(`Expected\n.*\n`)
)

// MixinMatcherGomock should be used for embedding to create a matcher that fits Gomock interface
type MixinMatcherGomock struct {
	cachedActual  any
	messagePrefix *string

	omega types.GomegaMatcher
}

func (igm *MixinMatcherGomock) Matches(v any) bool {
	// todo: we might cache the err if needed
	success, _ := igm.omega.Match(v)
	return success
}

func (igm *MixinMatcherGomock) String() string {
	gomegaFailureMessage := igm.omega.FailureMessage(igm.cachedActual)

	// considering that Failure message is a standard message made by FormatMessage
	// If the message is:
	// > Expected
	// >  <something>
	// > <something more>
	// it will remove first 2 lines
	if ExpectedStr2LinesRegex.MatchString(gomegaFailureMessage) {
		gomegaFailureMessage = ExpectedStr2LinesRegex.ReplaceAllString(gomegaFailureMessage, "")
	}

	// build a prefix (either was given, or by type of underlying matcher)
	var messagePrefix string
	if igm.messagePrefix != nil {
		messagePrefix = *igm.messagePrefix
	} else {
		messagePrefix = fmt.Sprintf("%T", igm.omega)
	}

	// Ensure prefix and message is separate by a single space
	messagePrefix = strings.TrimSuffix(messagePrefix, " ")
	if messagePrefix != "" {
		messagePrefix += " "
	}
	gomegaFailureMessage = strings.TrimPrefix(gomegaFailureMessage, " ")
	return messagePrefix + gomegaFailureMessage
}

func NewMixinMatcherGomock(Ω types.GomegaMatcher, messagePrefixArg ...string) *MixinMatcherGomock {
	igm := &MixinMatcherGomock{omega: Ω}
	if len(messagePrefixArg) > 0 {
		igm.messagePrefix = new(string)
		*igm.messagePrefix = messagePrefixArg[0]
	}

	return igm
}
