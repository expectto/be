package be_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/expectto/be"
)

func TestNilMatchers(t *testing.T) {
	var p *int // typed nil
	be.Expect(t, p).To(be.Nil())
	be.Expect(t, nil).To(be.Nil())
	be.Expect(t, 5).To(be.NotNil())
	be.Expect(t, p).NotTo(be.NotNil())
}

func TestBoolMatchers(t *testing.T) {
	be.Expect(t, true).To(be.True())
	be.Expect(t, false).To(be.False())
	be.Expect(t, true).NotTo(be.False())
	be.Expect(t, false).NotTo(be.True())
}

func TestErrorMatchers(t *testing.T) {
	var noErr error
	be.Expect(t, noErr).To(be.Succeed())

	sentinel := errors.New("boom")
	wrapped := fmt.Errorf("ctx: %w", sentinel)
	be.Expect(t, wrapped).To(be.HaveOccurred())
	be.Expect(t, wrapped).To(be.MatchError(sentinel))    // errors.Is through the wrap
	be.Expect(t, wrapped).To(be.MatchError("ctx: boom")) // message comparison
	be.Expect(t, noErr).NotTo(be.HaveOccurred())
}

func TestPanicMatchers(t *testing.T) {
	be.Expect(t, func() { panic("kaboom") }).To(be.Panic())
	be.Expect(t, func() {}).To(be.NotPanic())
}

func TestCollectionMatchers(t *testing.T) {
	be.Expect(t, []int{1, 2, 3}).To(be.ContainElement(2))
	be.Expect(t, []int{1, 2, 3}).To(be.ContainElements(3, 1))
	be.Expect(t, []int{1, 2, 3}).NotTo(be.ContainElement(9))

	m := map[string]int{"a": 1, "b": 2}
	be.Expect(t, m).To(be.HaveKey("a"))
	be.Expect(t, m).To(be.HaveKeyWithValue("b", 2))
	be.Expect(t, m).NotTo(be.HaveKey("z"))
}
