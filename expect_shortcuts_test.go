package be_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/expectto/be"
)

func TestNoError(t *testing.T) {
	// pass: nil error reports nothing
	rt := &recT{}
	if ok := be.NoError(rt, nil); !ok {
		t.Fatalf("expected success")
	}
	if len(rt.errs) != 0 || len(rt.fatals) != 0 {
		t.Fatalf("passing NoError must not report: errs=%v fatals=%v", rt.errs, rt.fatals)
	}

	// fail: non-nil error is HARD (Fatalf, not Errorf)
	rt = &recT{}
	if ok := be.NoError(rt, errors.New("boom")); ok {
		t.Fatalf("expected failure")
	}
	if len(rt.errs) != 0 {
		t.Fatalf("NoError must be hard (Fatalf); got errs=%v", rt.errs)
	}
	if len(rt.fatals) != 1 || !strings.Contains(rt.fatals[0], "boom") {
		t.Fatalf("failure should mention the error, got: %v", rt.fatals)
	}

	// message passthrough; the inline %q doubles as a vet regression guard —
	// go test runs vet, and this line fails the build if the assertion chain
	// is ever again classified as a print wrapper (see formatMsgAndArgs)
	rt = &recT{}
	be.NoError(rt, errors.New("boom"), "loading config %q", "app.yaml")
	if len(rt.fatals) != 1 || !strings.HasPrefix(rt.fatals[0], `loading config "app.yaml": `) {
		t.Fatalf("message context should be prepended, got: %v", rt.fatals)
	}
}

func TestErrorShortcut(t *testing.T) {
	// pass: non-nil error
	rt := &recT{}
	if ok := be.Error(rt, errors.New("boom")); !ok {
		t.Fatalf("expected success")
	}
	if len(rt.errs) != 0 || len(rt.fatals) != 0 {
		t.Fatalf("passing Error must not report: errs=%v fatals=%v", rt.errs, rt.fatals)
	}

	// fail: nil error is HARD
	rt = &recT{}
	if ok := be.Error(rt, nil); ok {
		t.Fatalf("expected failure")
	}
	if len(rt.fatals) != 1 || len(rt.errs) != 0 {
		t.Fatalf("Error must be hard (Fatalf); got fatals=%v errs=%v", rt.fatals, rt.errs)
	}
}

func TestErrorIs(t *testing.T) {
	sentinel := errors.New("sentinel")
	wrapped := fmt.Errorf("ctx: %w", sentinel)

	// pass: errors.Is through the wrap
	rt := &recT{}
	if ok := be.ErrorIs(rt, wrapped, sentinel); !ok {
		t.Fatalf("expected success")
	}
	if len(rt.errs) != 0 || len(rt.fatals) != 0 {
		t.Fatalf("passing ErrorIs must not report: errs=%v fatals=%v", rt.errs, rt.fatals)
	}

	// fail: unrelated error is HARD
	rt = &recT{}
	if ok := be.ErrorIs(rt, errors.New("other"), sentinel); ok {
		t.Fatalf("expected failure")
	}
	if len(rt.fatals) != 1 || len(rt.errs) != 0 {
		t.Fatalf("ErrorIs must be hard (Fatalf); got fatals=%v errs=%v", rt.fatals, rt.errs)
	}

	// fail: nil error
	rt = &recT{}
	if ok := be.ErrorIs(rt, nil, sentinel); ok {
		t.Fatalf("nil error must not match a target")
	}

	// message passthrough with an inline directive (vet regression guard)
	rt = &recT{}
	be.ErrorIs(rt, nil, sentinel, "checking case %d", 7)
	if len(rt.fatals) != 1 || !strings.HasPrefix(rt.fatals[0], "checking case 7: ") {
		t.Fatalf("message context should be prepended, got: %v", rt.fatals)
	}
}

func TestZeroNonZero(t *testing.T) {
	type cfg struct {
		Host string
		Port int
	}

	be.Expect(t, 0).To(be.Zero())
	be.Expect(t, "").To(be.Zero())
	be.Expect(t, false).To(be.Zero())
	be.Expect(t, cfg{}).To(be.Zero())
	var p *int
	be.Expect(t, p).To(be.Zero())

	be.Expect(t, 1).To(be.NonZero())
	be.Expect(t, "x").To(be.NonZero())
	be.Expect(t, cfg{Port: 8080}).To(be.NonZero())
	be.Expect(t, 0).NotTo(be.NonZero())
	be.Expect(t, cfg{}).NotTo(be.NonZero())
}

func TestRootNumericAliases(t *testing.T) {
	be.Expect(t, 10).To(be.Gt(5))
	be.Expect(t, 10).To(be.Gte(10))
	be.Expect(t, 3).To(be.Lt(5))
	be.Expect(t, 5).To(be.Lte(5))
	be.Expect(t, 10).To(be.GreaterThan(5))
	be.Expect(t, 10).To(be.GreaterThanEqual(10))
	be.Expect(t, 3).To(be.LessThan(5))
	be.Expect(t, 5).To(be.LessThanEqual(5))
	be.Expect(t, 5).To(be.InRange(1, true, 10, true))
	be.Expect(t, 5).To(be.Positive())
	be.Expect(t, -5).To(be.Negative())

	// failure message stays native/compact through the alias
	rt := &recT{}
	be.Expect(rt, 3).To(be.Gt(5))
	if len(rt.errs) != 1 || rt.errs[0] != "Expected 3 to be > 5" {
		t.Fatalf("unexpected failure output: %v", rt.errs)
	}
}
