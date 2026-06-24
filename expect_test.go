package be_test

import (
	"fmt"
	"testing"

	"github.com/expectto/be"
	"github.com/expectto/be/be_math"
)

// recT is a fake be.TestingT that records calls instead of failing a real test,
// so we can assert how Expect/Require react to passing and failing matches.
type recT struct {
	helpers int
	errs    []string
	fatals  []string
}

func (r *recT) Helper()                   { r.helpers++ }
func (r *recT) Errorf(f string, a ...any) { r.errs = append(r.errs, fmt.Sprintf(f, a...)) }
func (r *recT) Fatalf(f string, a ...any) { r.fatals = append(r.fatals, fmt.Sprintf(f, a...)) }

// The whole point of the native driver: the stdlib *testing.T satisfies it.
var _ be.TestingT = (*testing.T)(nil)

func TestExpectWorksWithRealTestingT(t *testing.T) {
	be.Expect(t, 10).To(be_math.GreaterThan(5))
	be.Require(t, 10).ToNot(be_math.LessThan(0))
}

func TestExpectSoftPass(t *testing.T) {
	rt := &recT{}
	ok := be.Expect(rt, 10).To(be_math.GreaterThan(5))
	if !ok {
		t.Fatalf("expected success")
	}
	if len(rt.errs) != 0 || len(rt.fatals) != 0 {
		t.Fatalf("a passing assertion must not report failures: errs=%v fatals=%v", rt.errs, rt.fatals)
	}
}

func TestExpectSoftFailReportsCompactMessage(t *testing.T) {
	rt := &recT{}
	ok := be.Expect(rt, 3).To(be_math.GreaterThan(5))
	if ok {
		t.Fatalf("expected failure")
	}
	if len(rt.fatals) != 0 {
		t.Fatalf("Expect must be soft (Errorf), not fatal; got fatals=%v", rt.fatals)
	}
	if len(rt.errs) != 1 {
		t.Fatalf("expected exactly one Errorf, got %v", rt.errs)
	}
	if rt.errs[0] != "Expected 3 to be > 5" {
		t.Fatalf("message not compact/native, got: %q", rt.errs[0])
	}
}

func TestRequireFailIsFatal(t *testing.T) {
	rt := &recT{}
	be.Require(rt, 3).To(be_math.GreaterThan(5))
	if len(rt.fatals) != 1 {
		t.Fatalf("Require must be fatal (Fatalf); got fatals=%v errs=%v", rt.fatals, rt.errs)
	}
}

func TestNotTo(t *testing.T) {
	rt := &recT{}
	if ok := be.Expect(rt, 3).NotTo(be_math.GreaterThan(5)); !ok {
		t.Fatalf("3 is not > 5, so NotTo should pass")
	}
	if ok := be.Expect(rt, 10).ToNot(be_math.GreaterThan(5)); ok {
		t.Fatalf("10 IS > 5, so ToNot should fail")
	}
}

func TestExpectMarksHelper(t *testing.T) {
	rt := &recT{}
	be.Expect(rt, 1).To(be_math.GreaterThan(0))
	if rt.helpers == 0 {
		t.Fatalf("Expect should call t.Helper so failures point at the caller line")
	}
}
