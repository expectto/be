package be_test

import (
	"errors"
	"fmt"
	"strings"
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

func TestEmptyAndNe(t *testing.T) {
	be.Expect(t, "").To(be.Empty())
	be.Expect(t, []int{}).To(be.Empty())
	be.Expect(t, map[string]int{}).To(be.Empty())
	be.Expect(t, []int{1}).To(be.NotEmpty())

	be.Expect(t, 5).To(be.Ne(6))
	be.Expect(t, 5).NotTo(be.Ne(5))
}

func TestContainSubstring(t *testing.T) {
	be.Expect(t, "hello world").To(be.ContainSubstring("o w"))
	be.Expect(t, "hello").NotTo(be.ContainSubstring("xyz"))
}

func TestMatchErrorAs(t *testing.T) {
	sentinel := &pathError{path: "/etc/passwd"}
	wrapped := fmt.Errorf("opening: %w", sentinel)

	be.Expect(t, wrapped).To(be.MatchErrorAs[*pathError]())                // wrapped typed error matches
	be.Expect(t, sentinel).To(be.MatchErrorAs[*pathError]())               // direct match
	be.Expect(t, errors.New("plain")).NotTo(be.MatchErrorAs[*pathError]()) // unrelated error fails

	var nilErr error
	be.Expect(t, nilErr).NotTo(be.MatchErrorAs[*pathError]()) // nil fails

	// failure message names the type
	rt := &recT{}
	be.Expect(rt, errors.New("plain")).To(be.MatchErrorAs[*pathError]())
	if len(rt.errs) != 1 || !strings.Contains(rt.errs[0], "*be_test.pathError") {
		t.Fatalf("failure should name the target type, got: %v", rt.errs)
	}

	// non-error actual is an error, not a mismatch
	rt = &recT{}
	be.Expect(rt, 42).To(be.MatchErrorAs[*pathError]())
	if len(rt.errs) != 1 || !strings.Contains(rt.errs[0], "expects an error") {
		t.Fatalf("non-error actual should produce a clear error, got: %v", rt.errs)
	}
}

type pathError struct{ path string }

func (e *pathError) Error() string { return "path error: " + e.path }

func TestHaveField(t *testing.T) {
	type address struct{ City string }
	type user struct {
		Name    string
		Age     int
		Address address
	}
	u := user{Name: "Alice", Age: 30, Address: address{City: "Lisbon"}}

	be.Expect(t, u).To(be.HaveField("Name", "Alice"))
	be.Expect(t, u).To(be.HaveField("Age", be.Gte(18)))        // matcher value
	be.Expect(t, u).To(be.HaveField("Address.City", "Lisbon")) // nesting
	be.Expect(t, &u).To(be.HaveField("Name", "Alice"))         // pointer to struct
	be.Expect(t, u).NotTo(be.HaveField("Name", "Bob"))

	// method form
	be.Expect(t, &pathError{path: "/x"}).To(be.HaveField("Error()", "path error: /x"))

	// field mismatch names the field
	rt := &recT{}
	be.Expect(rt, u).To(be.HaveField("Name", "Bob"))
	if len(rt.errs) != 1 || !strings.Contains(rt.errs[0], "Name") {
		t.Fatalf("failure should name the field, got: %v", rt.errs)
	}
}

func TestHaveFields(t *testing.T) {
	type user struct {
		Name string
		Age  int
	}
	u := user{Name: "Alice", Age: 30}

	be.Expect(t, u).To(be.HaveFields(map[string]any{
		"Name": "Alice",
		"Age":  be.Gt(18),
	}))
	be.Expect(t, u).NotTo(be.HaveFields(map[string]any{
		"Name": "Alice",
		"Age":  be.Gt(50),
	}))

	// deterministic failure output: with two mismatching fields, the first
	// (sorted) key is always the one reported
	for range 10 {
		rt := &recT{}
		be.Expect(rt, u).To(be.HaveFields(map[string]any{
			"Name": "Bob",
			"Age":  99,
		}))
		if len(rt.errs) != 1 || !strings.Contains(rt.errs[0], "Age") {
			t.Fatalf("sorted-key order should make Age fail first, got: %v", rt.errs)
		}
	}
}

func TestHaveLengthComposable(t *testing.T) {
	// HaveLength accepts a count...
	be.Expect(t, []int{1, 2, 3}).To(be.HaveLength(3))
	// ...or a matcher for the length
	be.Expect(t, []int{1, 2, 3}).To(be.HaveLength(be.Gte(2)))
	be.Expect(t, "hello").To(be.HaveLength(be.InRange(1, true, 10, true)))
	be.Expect(t, []int{1}).NotTo(be.HaveLength(be.Gte(2)))
}

func TestIdenticalAndVia(t *testing.T) {
	type box struct{ n int }
	p := &box{n: 1}
	be.Expect(t, p).To(be.Identical(p))
	be.Expect(t, p).To(be.NotIdentical(&box{n: 1})) // equal value, different pointer

	// Via projects through an accessor before matching.
	get := func(b box) int { return b.n }
	be.Expect(t, box{n: 42}).To(be.Via(get, be.Eq(42)))
	be.Expect(t, box{n: 42}).NotTo(be.Via(get, be.Eq(0)))
}
