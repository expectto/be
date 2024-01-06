package be_ctx

import (
	"context"
	"fmt"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega/format"
	"go.uber.org/mock/gomock"
)

var (
	FailCtxNotAContext        = fmt.Errorf("ctx is not a context")
	FailCtxValueExpected      = fmt.Errorf("ctx value expected but not found")
	FailCtxValueNotMatched    = fmt.Errorf("ctx value failed to match")
	FailCtxErrNotMatched      = fmt.Errorf("ctx err() failed to match")
	FailCtxDeadlineExpected   = fmt.Errorf("ctx deadlined expected but not found")
	FailCtxDeadlineNotMatched = fmt.Errorf("ctx deadlined failed to match")
)

// CtxMatcher is a matcher for ctx// Each instance of CtxMatcher can match across only one thing:// (1) ctx value or (2) error or (3) deadline or (4) done signal
// Do not fill multiple things together here. Use separate instances instead
type CtxMatcher struct {
	failReason error

	// 1.  Matching inner ctx value:
	expectedValueKey string
	valueMatcher     types.BeMatcher

	// 2. Error matching
	errMatcher types.BeMatcher

	// 3. Deadline matching
	deadlineMatcher types.BeMatcher

	// 4. Done matching
	doneMatcher types.BeMatcher
}

var _ gomock.Matcher = &CtxMatcher{}
var _ types.BeMatcher = &CtxMatcher{}

func (cm *CtxMatcher) Match(v any) (success bool, err error) {
	return cm.match(v)
}

func (cm *CtxMatcher) FailureMessage(v any) (message string) {
	return format.Message(v, "failed")
}

func (cm *CtxMatcher) NegatedFailureMessage(v any) (message string) {
	return format.Message(v, "not failed")
}

func (cm *CtxMatcher) Matches(v any) bool {
	success, _ := cm.match(v)
	return success
}

func (cm *CtxMatcher) String() string {
	return "failed"
}

func (cm *CtxMatcher) match(v any) (bool, error) {
	ctx, ok := v.(context.Context)
	if !ok {
		cm.failReason = FailCtxNotAContext
		return false, nil
	}
	// (1) matching context value
	if cm.expectedValueKey != "" {
		foundValue := ctx.Value(cm.expectedValueKey)

		if cm.valueMatcher == nil {
			// simply match existance of a value
			if foundValue == nil {
				cm.failReason = FailCtxValueExpected
				return false, nil
			}

			return true, nil
		}

		succeed, err := cm.valueMatcher.Match(foundValue)
		if err != nil {
			return false, err
		}
		if !succeed {
			cm.failReason = fmt.Errorf("%w: %s", FailCtxValueNotMatched, cm.valueMatcher.FailureMessage(foundValue))
		}
		return succeed, nil
	}
	// (2) matching context err
	if cm.errMatcher != nil {
		succeed, err := cm.errMatcher.Match(ctx.Err())
		if err != nil {
			return false, err
		}
		if !succeed {
			cm.failReason = FailCtxErrNotMatched
		}
		return succeed, nil
	}
	// (3) matching context deadline
	if cm.deadlineMatcher != nil {
		// first simple check if deadline exists
		deadline, ok := ctx.Deadline()
		if !ok {
			cm.failReason = FailCtxDeadlineExpected
			return false, nil
		}

		succeed, err := cm.deadlineMatcher.Match(deadline)
		if err != nil {
			return false, err
		}
		if !succeed {
			cm.failReason = fmt.Errorf("%w: %s", FailCtxDeadlineNotMatched, cm.deadlineMatcher.FailureMessage(deadline))
		}
		return succeed, nil
	}
	// (4) matching context Done signal TODO
	//ctx.Done()

	return true, nil
}

func NewCtxMatcher() *CtxMatcher {
	return &CtxMatcher{}
}

func NewCtxValueMatcher(key string, m types.BeMatcher) *CtxMatcher {
	return &CtxMatcher{expectedValueKey: key, valueMatcher: m}
}

func NewCtxDeadlineMatcher(m types.BeMatcher) *CtxMatcher {
	return &CtxMatcher{deadlineMatcher: m}
}

func NewCtxErrMatcher(m types.BeMatcher) *CtxMatcher {
	return &CtxMatcher{errMatcher: m}
}
