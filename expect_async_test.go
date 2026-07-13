package be_test

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/expectto/be"
)

// Tests count polls instead of measuring wall-clock time: correctness is
// "flips after N polls", not "took N milliseconds".

func TestEventuallyFlipsAfterNPolls(t *testing.T) {
	var polls atomic.Int64
	counter := func() int64 { return polls.Add(1) }

	rt := &recT{}
	ok := be.Eventually(rt, counter, be.Gte(int64(3)), be.WithPolling(time.Millisecond))
	if !ok {
		t.Fatalf("expected Eventually to succeed once the counter reaches 3")
	}
	if got := polls.Load(); got < 3 {
		t.Fatalf("expected at least 3 polls, got %d", got)
	}
	if len(rt.errs) != 0 || len(rt.fatals) != 0 {
		t.Fatalf("passing Eventually must not report: errs=%v fatals=%v", rt.errs, rt.fatals)
	}
}

func TestEventuallyPlainValue(t *testing.T) {
	be.Eventually(t, 5, be.Eq(5))
}

func TestEventuallyTimeoutFailure(t *testing.T) {
	rt := &recT{}
	ok := be.Eventually(rt, func() int { return 3 }, be.Gt(5),
		be.WithTimeout(20*time.Millisecond), be.WithPolling(2*time.Millisecond))
	if ok {
		t.Fatalf("expected failure")
	}
	if len(rt.fatals) != 0 {
		t.Fatalf("Eventually must fail softly (Errorf); got fatals=%v", rt.fatals)
	}
	if len(rt.errs) != 1 {
		t.Fatalf("expected exactly one Errorf, got %v", rt.errs)
	}
	// timeout context + the last compact mismatch, same format as be.Expect
	if !strings.Contains(rt.errs[0], "timed out after 20ms") ||
		!strings.Contains(rt.errs[0], "Expected 3 to be > 5") {
		t.Fatalf("unexpected failure message: %q", rt.errs[0])
	}
}

func TestEventuallyFuncReturningError(t *testing.T) {
	var polls atomic.Int64
	fetch := func() (int, error) {
		if polls.Add(1) < 3 {
			return 0, errors.New("not ready yet")
		}
		return 42, nil
	}

	rt := &recT{}
	if ok := be.Eventually(rt, fetch, be.Eq(42), be.WithPolling(time.Millisecond)); !ok {
		t.Fatalf("errors from the poll function should mean 'retry', not 'fail'")
	}

	// an error on every poll surfaces in the timeout message
	rt = &recT{}
	alwaysErr := func() (int, error) { return 0, errors.New("kaput") }
	ok := be.Eventually(rt, alwaysErr, be.Eq(42),
		be.WithTimeout(15*time.Millisecond), be.WithPolling(2*time.Millisecond))
	if ok {
		t.Fatalf("expected failure")
	}
	if len(rt.errs) != 1 || !strings.Contains(rt.errs[0], "polled function returned error: kaput") {
		t.Fatalf("unexpected failure message: %v", rt.errs)
	}
}

func TestEventuallyWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already-done context fails on the first non-matching poll

	rt := &recT{}
	ok := be.Eventually(rt, func() int { return 3 }, be.Gt(5),
		be.WithContext(ctx), be.WithPolling(time.Millisecond))
	if ok {
		t.Fatalf("expected failure")
	}
	if len(rt.errs) != 1 || !strings.Contains(rt.errs[0], "context done") {
		t.Fatalf("unexpected failure message: %v", rt.errs)
	}
}

func TestConsistentlyHolds(t *testing.T) {
	var polls atomic.Int64
	counter := func() int64 { polls.Add(1); return 1 }

	rt := &recT{}
	ok := be.Consistently(rt, counter, be.Eq(int64(1)),
		be.WithTimeout(15*time.Millisecond), be.WithPolling(time.Millisecond))
	if !ok {
		t.Fatalf("expected Consistently to succeed")
	}
	if got := polls.Load(); got < 2 {
		t.Fatalf("Consistently must poll more than once, got %d polls", got)
	}
	if len(rt.errs) != 0 || len(rt.fatals) != 0 {
		t.Fatalf("passing Consistently must not report: errs=%v fatals=%v", rt.errs, rt.fatals)
	}
}

func TestConsistentlyFailsOnFirstMismatch(t *testing.T) {
	var polls atomic.Int64
	flipsBad := func() int64 { return polls.Add(1) } // 1, 2, 3... breaks be.Lte(2) on poll 3

	rt := &recT{}
	ok := be.Consistently(rt, flipsBad, be.Lte(int64(2)), be.WithPolling(time.Millisecond))
	if ok {
		t.Fatalf("expected failure on the third poll")
	}
	if got := polls.Load(); got != 3 {
		t.Fatalf("Consistently must stop at the first mismatch (poll 3), got %d polls", got)
	}
	if len(rt.fatals) != 0 {
		t.Fatalf("Consistently must fail softly (Errorf); got fatals=%v", rt.fatals)
	}
	if len(rt.errs) != 1 || !strings.Contains(rt.errs[0], "Expected 3 to be <= 2") {
		t.Fatalf("unexpected failure message: %v", rt.errs)
	}
}

func TestEventuallyRejectsBadPollFunc(t *testing.T) {
	rt := &recT{}
	ok := be.Eventually(rt, func() (int, int) { return 1, 2 }, be.Eq(1),
		be.WithTimeout(10*time.Millisecond))
	if ok {
		t.Fatalf("expected failure for a func returning (int, int)")
	}
	if len(rt.errs) != 1 || !strings.Contains(rt.errs[0], "second return value must be an error") {
		t.Fatalf("unexpected failure message: %v", rt.errs)
	}
}
