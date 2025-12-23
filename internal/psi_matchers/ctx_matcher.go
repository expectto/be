package psi_matchers

import (
	"context"
	"fmt"

	. "github.com/expectto/be/internal/psi" //nolint:staticcheck // should be moved to lintignore
	"github.com/expectto/be/types"
	"github.com/onsi/gomega/format"
	"go.uber.org/mock/gomock"
)

var (
	ErrCtxNotAContext        = fmt.Errorf("be a ctx")
	ErrCtxValueExpected      = fmt.Errorf("have the ctx.value")
	ErrCtxValueNotMatched    = fmt.Errorf("have the ctx.value") // same text as won't be used directly (but still can be distinguished via errors.Is()
	ErrCtxErrNotMatched      = fmt.Errorf("todo: better error handling here")
	ErrCtxDeadlineExpected   = fmt.Errorf("todo: better error handling here")
	ErrCtxDeadlineNotMatched = fmt.Errorf("todo: better error handling here")
)

// CtxMatcher is a matcher for ctx// Each instance of CtxMatcher can match across only one thing:// (1) ctx value or (2) error or (3) deadline or (4) done signal
// Do not fill multiple things together here. Use separate instances instead
type CtxMatcher struct {
	failReason error

	// 1.  Matching inner ctx value:
	key   any
	value any

	// 2. Error matching
	errFn any

	// 3. Deadline matching
	deadline any

	// 4. Done matching
	//doneMatcher types.BeMatcher
}

var _ gomock.Matcher = &CtxMatcher{}
var _ types.BeMatcher = &CtxMatcher{}

func (cm *CtxMatcher) Match(v any) (success bool, err error) {
	return cm.match(v)
}

func (cm *CtxMatcher) FailureMessage(v any) (message string) {
	return format.Message(v, fmt.Sprintf("to %s", cm.failReason))
}

func (cm *CtxMatcher) NegatedFailureMessage(v any) (message string) {
	return format.Message(v, fmt.Sprintf("not to %s", cm.failReason))
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
		cm.failReason = ErrCtxNotAContext
		return false, nil
	}

	// (1) matching context value
	if cm.key != nil {
		foundValue := ctx.Value(cm.key)

		if cm.value == nil {
			// simply match existence of a value
			if foundValue == nil {
				cm.failReason = fmt.Errorf("%w key=`%s`", ErrCtxValueExpected, cm.key)
				return false, nil
			}

			return true, nil
		}

		valueMatcher := Psi(cm.value)
		succeed, err := valueMatcher.Match(foundValue)
		if err != nil {
			return false, err
		}
		if !succeed {
			cm.failReason = fmt.Errorf("%w key=`%s` that failed on match:\n%s", ErrCtxValueNotMatched, cm.key, valueMatcher.FailureMessage(foundValue))
		}
		return succeed, nil
	}
	// (2) matching context err
	if cm.errFn != nil {
		errMatcher := Psi(cm.errFn)
		succeed, err := errMatcher.Match(ctx.Err())
		if err != nil {
			return false, err
		}
		if !succeed {
			cm.failReason = fmt.Errorf("%w: %s", ErrCtxErrNotMatched, errMatcher.FailureMessage(ctx.Err()))
		}
		return succeed, nil
	}
	// (3) matching context deadline
	if cm.deadline != nil {
		// first simple check if deadline exists
		deadline, ok := ctx.Deadline()
		if !ok {
			cm.failReason = ErrCtxDeadlineExpected
			return false, nil
		}

		deadlineMatcher := Psi(cm.deadline)
		succeed, err := deadlineMatcher.Match(deadline)
		if err != nil {
			return false, err
		}
		if !succeed {
			cm.failReason = fmt.Errorf("%w: %s", ErrCtxDeadlineNotMatched, deadlineMatcher.FailureMessage(deadline))
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

func NewCtxValueMatcher(key any, valueArg ...any) *CtxMatcher {
	matcher := &CtxMatcher{key: key}
	switch len(valueArg) {
	case 0:
		return matcher
	case 1:
		matcher.value = valueArg[0]
		return matcher
	default:
		panic("NewCtxValueMatcher expects either 0 or 1 value matcher")
	}
}

func NewCtxDeadlineMatcher(deadline any) *CtxMatcher {
	return &CtxMatcher{deadline: deadline}
}

func NewCtxErrMatcher(errFn any) *CtxMatcher {
	return &CtxMatcher{errFn: errFn}
}
