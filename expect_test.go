package be_test

import (
	"fmt"
	"strings"
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

func TestExpectMessageContext(t *testing.T) {
	rt := &recT{}
	be.Expect(rt, 3).To(be_math.GreaterThan(5), "checking case 7")
	if len(rt.errs) != 1 {
		t.Fatalf("expected one failure, got %v", rt.errs)
	}
	if !strings.HasPrefix(rt.errs[0], "checking case 7: ") {
		t.Fatalf("message context should be prepended, got: %q", rt.errs[0])
	}

	// format-string form (built with fmt to avoid go vet's printf heuristic on To)
	rt2 := &recT{}
	be.Expect(rt2, 3).To(be_math.GreaterThan(5), fmt.Sprintf("case %d", 7))
	if !strings.HasPrefix(rt2.errs[0], "case 7: ") {
		t.Fatalf("formatted context should be prepended, got: %q", rt2.errs[0])
	}
}

func TestExpectMarksHelper(t *testing.T) {
	rt := &recT{}
	be.Expect(rt, 1).To(be_math.GreaterThan(0))
	if rt.helpers == 0 {
		t.Fatalf("Expect should call t.Helper so failures point at the caller line")
	}
}

func TestAssertThatSoft(t *testing.T) {
	// Passing match reports nothing.
	rt := &recT{}
	if ok := be.AssertThat(rt, 10, be_math.GreaterThan(5)); !ok {
		t.Fatalf("expected success")
	}
	if len(rt.errs) != 0 || len(rt.fatals) != 0 {
		t.Fatalf("a passing assertion must not report: errs=%v fatals=%v", rt.errs, rt.fatals)
	}

	// Failing match is soft (Errorf, not Fatalf) and carries the compact message.
	rt = &recT{}
	if ok := be.AssertThat(rt, 3, be_math.GreaterThan(5)); ok {
		t.Fatalf("expected failure")
	}
	if len(rt.fatals) != 0 {
		t.Fatalf("AssertThat must be soft; got fatals=%v", rt.fatals)
	}
	if len(rt.errs) != 1 || rt.errs[0] != "Expected 3 to be > 5" {
		t.Fatalf("unexpected soft failure output: %v", rt.errs)
	}
}

func TestRequireThatIsFatal(t *testing.T) {
	rt := &recT{}
	be.RequireThat(rt, 3, be_math.GreaterThan(5))
	if len(rt.fatals) != 1 {
		t.Fatalf("RequireThat must be fatal (Fatalf); got fatals=%v errs=%v", rt.fatals, rt.errs)
	}
}

func TestAssertThatMessageContext(t *testing.T) {
	rt := &recT{}
	be.AssertThat(rt, 3, be_math.GreaterThan(5), "checking case 7")
	if len(rt.errs) != 1 || !strings.HasPrefix(rt.errs[0], "checking case 7: ") {
		t.Fatalf("message context should be prepended, got: %v", rt.errs)
	}
}
