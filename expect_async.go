package be

// expect_async.go provides native async assertions: Eventually and Consistently.
// They are a be-native poll loop against TestingT — NOT a wrapper around
// gomega's Eventually — so failures go through the same path as Expectation.To
// and keep the compact, framework-native message format (see internal/beformat).

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/expectto/be/internal/beformat"
	"github.com/expectto/be/internal/psi"
)

const (
	defaultEventuallyTimeout    = time.Second
	defaultConsistentlyDuration = 100 * time.Millisecond
	defaultPollingInterval      = 10 * time.Millisecond
)

// asyncConfig carries the polling knobs for Eventually / Consistently.
type asyncConfig struct {
	timeout time.Duration
	polling time.Duration
	ctx     context.Context
}

// EventuallyOption configures Eventually and Consistently.
type EventuallyOption func(*asyncConfig)

// WithTimeout sets how long Eventually keeps polling before failing
// (default 1s), or how long Consistently keeps verifying before succeeding
// (default 100ms).
func WithTimeout(d time.Duration) EventuallyOption {
	return func(c *asyncConfig) { c.timeout = d }
}

// WithPolling sets the interval between polls (default 10ms).
func WithPolling(d time.Duration) EventuallyOption {
	return func(c *asyncConfig) { c.polling = d }
}

// WithContext bounds the poll loop by a context: when the context is done the
// assertion fails immediately instead of waiting for the timeout.
func WithContext(ctx context.Context) EventuallyOption {
	return func(c *asyncConfig) { c.ctx = ctx }
}

// Eventually polls actual until it satisfies the matcher, failing the test
// (softly, via Errorf) if it never does within the timeout. actual may be:
//   - a plain value (matched repeatedly — useful for stateful matchers),
//   - func() T — polled each interval,
//   - func() (T, error) — a returned error means "not ready yet"; polling continues.
//
// Example:
//
//	be.Eventually(t, queue.Len, be.Gte(3))
//	be.Eventually(t, fetchStatus, be.Eq("ready"), be.WithTimeout(5*time.Second))
//
// The failure message reports the last mismatch in the same compact format as
// be.Expect. Returns true on success.
func Eventually(t TestingT, actual any, matcher any, opts ...EventuallyOption) bool {
	t.Helper()
	cfg := newAsyncConfig(defaultEventuallyTimeout, opts)
	e := &Expectation{t: t, actual: actual}

	poll, err := asPoller(actual)
	if err != nil {
		return e.fail(err.Error())
	}

	m := psi.Psi(matcher)
	deadline := time.After(cfg.timeout)
	var ctxDone <-chan struct{}
	if cfg.ctx != nil {
		ctxDone = cfg.ctx.Done()
	}

	var lastMsg string
	for {
		v, pollErr := poll()
		if pollErr != nil {
			lastMsg = fmt.Sprintf("polled function returned error: %v", pollErr)
		} else if ok, merr := m.Match(v); merr != nil {
			lastMsg = merr.Error()
		} else if ok {
			return true
		} else {
			lastMsg = beformat.Compact(m.FailureMessage(v))
		}

		select {
		case <-ctxDone:
			return e.fail(fmt.Sprintf("context done while waiting: %s", lastMsg))
		case <-deadline:
			return e.fail(fmt.Sprintf("timed out after %s: %s", cfg.timeout, lastMsg))
		case <-time.After(cfg.polling):
		}
	}
}

// Consistently polls actual and requires it to satisfy the matcher on EVERY
// poll for the whole duration (default 100ms, set via WithTimeout). The first
// mismatch fails the test (softly, via Errorf) immediately. actual takes the
// same forms as in Eventually. Returns true if the matcher held throughout.
//
//	be.Consistently(t, queue.Len, be.Zero())
func Consistently(t TestingT, actual any, matcher any, opts ...EventuallyOption) bool {
	t.Helper()
	cfg := newAsyncConfig(defaultConsistentlyDuration, opts)
	e := &Expectation{t: t, actual: actual}

	poll, err := asPoller(actual)
	if err != nil {
		return e.fail(err.Error())
	}

	m := psi.Psi(matcher)
	deadline := time.After(cfg.timeout)
	var ctxDone <-chan struct{}
	if cfg.ctx != nil {
		ctxDone = cfg.ctx.Done()
	}

	for {
		v, pollErr := poll()
		if pollErr != nil {
			return e.fail(fmt.Sprintf("polled function returned error: %v", pollErr))
		}
		ok, merr := m.Match(v)
		if merr != nil {
			return e.fail(merr.Error())
		}
		if !ok {
			return e.fail(beformat.Compact(m.FailureMessage(v)))
		}

		select {
		case <-ctxDone:
			return e.fail("context done before the consistency window elapsed")
		case <-deadline:
			return true
		case <-time.After(cfg.polling):
		}
	}
}

func newAsyncConfig(defaultTimeout time.Duration, opts []EventuallyOption) *asyncConfig {
	cfg := &asyncConfig{timeout: defaultTimeout, polling: defaultPollingInterval}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

var errType = reflect.TypeFor[error]()

// asPoller normalizes actual into a poll function. Niladic funcs returning one
// value (or a value and an error) are invoked on each poll; anything else —
// including func() with no results, so be.Panic still works — is treated as a
// plain value and matched as-is.
func asPoller(actual any) (func() (any, error), error) {
	v := reflect.ValueOf(actual)
	if !v.IsValid() || v.Kind() != reflect.Func {
		return func() (any, error) { return actual, nil }, nil
	}
	ft := v.Type()
	if ft.NumIn() != 0 || ft.NumOut() == 0 {
		return func() (any, error) { return actual, nil }, nil
	}
	if ft.NumOut() > 2 {
		return nil, fmt.Errorf("Eventually/Consistently poll function must return (T) or (T, error), got %s", ft)
	}
	if ft.NumOut() == 2 && !ft.Out(1).Implements(errType) {
		return nil, fmt.Errorf("Eventually/Consistently poll function's second return value must be an error, got %s", ft)
	}

	return func() (any, error) {
		out := v.Call(nil)
		if len(out) == 2 && !out[1].IsNil() {
			return nil, out[1].Interface().(error)
		}
		return out[0].Interface(), nil
	}, nil
}
