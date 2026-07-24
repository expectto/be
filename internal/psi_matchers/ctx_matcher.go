package psi_matchers

import (
	"context"
	"errors"
	"fmt"

	"github.com/onsi/gomega/format"

	. "github.com/expectto/be/internal/psi" //nolint:staticcheck // should be moved to lintignore
	"github.com/expectto/be/types"
)

var (
	ErrCtxNotAContext     = errors.New("be a ctx")
	ErrCtxValueExpected   = errors.New("have the ctx.value")
	ErrCtxValueNotMatched = errors.New(
		"have the ctx.value",
	) // same text as won't be used directly (but still can be distinguished via errors.Is()
	ErrCtxErrorNotMatched    = errors.New("have the expected ctx error")
	ErrCtxDeadlineExpected   = errors.New("have a ctx deadline")
	ErrCtxDeadlineNotMatched = errors.New("have the expected ctx deadline")
)

// CtxMatcher is a matcher for ctx// Each instance of CtxMatcher can match across only one thing:// (1) ctx value or (2) error or (3) deadline or (4) done signal
// Do not fill multiple things together here. Use separate instances instead
type CtxMatcher struct {
	failReason error

	// 1.  Matching inner ctx value:
	key   any
	value any

	// 2. Error matching
	// matchErr records that error matching was requested (via NewCtxErrMatcher),
	// so a nil errFn means "expect no error" rather than "no error matcher set".
	matchErr bool
	errFn    any

	// 3. Deadline matching
	deadline any

	// 4. Done matching
	// doneMatcher types.BeMatcher
}

// types.BeMatcher embeds types.GomockMatcher, so this single assertion also
// guarantees gomock compatibility without importing go.uber.org/mock.
var _ types.BeMatcher = &CtxMatcher{}

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
	return &CtxMatcher{matchErr: true, errFn: errFn}
}

func (cm *CtxMatcher) Match(v any) (bool, error) {
	return cm.match(v)
}

func (cm *CtxMatcher) FailureMessage(v any) string {
	return format.Message(v, fmt.Sprintf("to %s", cm.failReason))
}

func (cm *CtxMatcher) NegatedFailureMessage(v any) string {
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
			cm.failReason = fmt.Errorf(
				"%w key=`%s` that failed on match:\n%s",
				ErrCtxValueNotMatched,
				cm.key,
				valueMatcher.FailureMessage(foundValue),
			)
		}
		return succeed, nil
	}
	// (2) matching context err
	if cm.matchErr {
		// CtxWithError(nil) asserts the context carries no error.
		if cm.errFn == nil {
			if ctx.Err() != nil {
				cm.failReason = fmt.Errorf("%w: expected no error, got %w", ErrCtxErrorNotMatched, ctx.Err())
				return false, nil //nolint:nilerr // ctx.Err() is the value under test; a non-nil ctx error means "no match", not a matcher failure
			}
			return true, nil
		}

		errMatcher := Psi(cm.errFn)
		succeed, err := errMatcher.Match(ctx.Err())
		if err != nil {
			return false, err
		}
		if !succeed {
			cm.failReason = fmt.Errorf("%w: %s", ErrCtxErrorNotMatched, errMatcher.FailureMessage(ctx.Err()))
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
	// ctx.Done()

	return true, nil
}
