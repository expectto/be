package psi_test

import (
	"testing"

	"github.com/expectto/be/internal/psi"
	"github.com/onsi/gomega"
)

// gt0 is a real matcher (actual > 0). A raw func(any)(bool,error) cannot be used
// here: be treats a one-value+error func as a transform, not a predicate.
func gt0() gomega.OmegaMatcher { return gomega.BeNumerically(">", 0) }

func TestDiveModes(t *testing.T) {
	g := gomega.NewWithT(t)

	every := psi.NewDiveMatcher(gt0(), psi.DiveModeEvery)
	g.Expect(every.Match([]int{1, 2, 3})).To(gomega.BeTrue())
	g.Expect(every.Match([]int{1, -2, 3})).To(gomega.BeFalse())

	any := psi.NewDiveMatcher(gt0(), psi.DiveModeAny)
	g.Expect(any.Match([]int{-1, -2, 3})).To(gomega.BeTrue())
	g.Expect(any.Match([]int{-1, -2, -3})).To(gomega.BeFalse())

	first := psi.NewDiveMatcher(gt0(), psi.DiveModeFirst)
	g.Expect(first.Match([]int{5, -1})).To(gomega.BeTrue())

	nth := psi.NewDiveMatcher(gt0(), psi.DiveModeNth, 1)
	g.Expect(nth.Match([]int{-1, 5})).To(gomega.BeTrue())
}

// TestDiveSupportsArrays guards the reflect-based conversion that lets Dive accept
// fixed-size arrays, not just slices.
func TestDiveSupportsArrays(t *testing.T) {
	g := gomega.NewWithT(t)
	m := psi.NewDiveMatcher(gt0(), psi.DiveModeEvery)
	g.Expect(m.Match([2]int{1, 2})).To(gomega.BeTrue())
}

// TestDiveFailsGracefullyOnNonSlice ensures a non-slice actual yields an error
// rather than panicking and crashing the whole test run.
func TestDiveFailsGracefullyOnNonSlice(t *testing.T) {
	g := gomega.NewWithT(t)
	for _, actual := range []any{42, "hello", nil, struct{}{}} {
		m := psi.NewDiveMatcher(gt0(), psi.DiveModeEvery)
		g.Expect(func() {
			ok, err := m.Match(actual)
			g.Expect(err).To(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeFalse())
		}).NotTo(gomega.Panic())
	}
}
