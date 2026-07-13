package a

import (
	"errors"
	"slices"
	"strings"

	"github.com/expectto/be"
)

var errSentinel = errors.New("sentinel")

type myErr struct{}

func (*myErr) Error() string { return "my" }

func rawComparisons(t be.TestingT, x, y, n int, s string) {
	be.AssertThat(t, x == y, be.True())  // want `prefer be\.Eq\(y\) over wrapping the raw expression in be\.True\(\)`
	be.AssertThat(t, x != y, be.True())  // want `prefer be\.Ne\(y\)`
	be.AssertThat(t, x == y, be.False()) // want `prefer be\.Ne\(y\)`
	be.AssertThat(t, x >= n, be.True())  // want `prefer be\.Gte\(n\)`
	be.AssertThat(t, x > 0, be.True())   // want `prefer be\.Gt\(0\)`
	be.AssertThat(t, 0 < x, be.True())   // want `prefer be\.Gt\(0\)`
	be.AssertThat(t, x < n, be.False())  // want `prefer be\.Gte\(n\)`
	be.AssertThat(t, !(x == y), be.True()) // want `prefer be\.Ne\(y\)`

	// ordered comparison on strings has no numeric matcher — not flagged
	be.AssertThat(t, s > "a", be.True())
	// plain bool subject is what be.True is FOR — not flagged
	ok := x > 0
	be.AssertThat(t, ok, be.True())
}

func lengths(t be.TestingT, xs []int, n int) {
	be.AssertThat(t, len(xs) == 0, be.True()) // want `prefer be\.Empty\(\)`
	be.AssertThat(t, len(xs) != 0, be.True()) // want `prefer be\.NotEmpty\(\)`
	be.AssertThat(t, len(xs) > 0, be.True())  // want `prefer be\.NotEmpty\(\)`
	be.AssertThat(t, len(xs) == n, be.True()) // want `prefer be\.HaveLength\(n\)`
	be.AssertThat(t, len(xs) >= n, be.True()) // want `prefer be\.HaveLength\(be\.Gte\(n\)\)`
	be.AssertThat(t, len(xs) == 0, be.False()) // want `prefer be\.NotEmpty\(\)`
}

func nils(t be.TestingT, err error) {
	be.AssertThat(t, err == nil, be.True()) // want `prefer be\.Nil\(\)`
	be.AssertThat(t, err != nil, be.True()) // want `prefer be\.NotNil\(\)`
	be.AssertThat(t, nil != err, be.True()) // want `prefer be\.NotNil\(\)`
}

func stdlibIdioms(t be.TestingT, xs []int, x int, s string, err error) {
	be.AssertThat(t, slices.Contains(xs, x), be.True())      // want `prefer be\.ContainElement\(x\)`
	be.AssertThat(t, slices.Contains(xs, x), be.False())     // want `prefer be\.Not\(be\.ContainElement\(x\)\)`
	be.AssertThat(t, strings.Contains(s, "ell"), be.True())  // want `prefer be\.ContainSubstring\("ell"\)`
	be.AssertThat(t, errors.Is(err, errSentinel), be.True()) // want `prefer be\.MatchError\(errSentinel\)`
	be.AssertThat(t, errors.Is(err, errSentinel), be.False()) // want `prefer be\.Not\(be\.MatchError\(errSentinel\)\)`

	// report-only (fix would need a new import / a type name):
	be.AssertThat(t, strings.HasPrefix(s, "he"), be.True()) // want `prefer be_string\.HavingPrefix\("he"\)`
	be.AssertThat(t, strings.HasSuffix(s, "lo"), be.True()) // want `prefer be_string\.HavingSuffix\("lo"\)`

	// errors.As with the target unused afterward: flagged
	var target *myErr
	_ = target
	be.AssertThat(t, errors.As(err, &target), be.True()) // want `prefer be\.MatchErrorAs\[\*myErr\]\(\)`
}

func errorsAsTargetUsedAfter(t be.TestingT, err error) {
	// MatchErrorAs does not bind the target, so when it is referenced after
	// the assertion the conversion is impossible — no diagnostic.
	var target *myErr
	be.AssertThat(t, errors.As(err, &target), be.True())
	if target != nil {
		_ = target.Error()
	}
}

func fluentForms(t be.TestingT, x, n int) {
	be.Expect(t, x > n).To(be.True())      // want `prefer be\.Gt\(n\)`
	be.Expect(t, x > n).NotTo(be.True())   // want `prefer be\.Lte\(n\)`
	be.Require(t, x == n).To(be.True())    // want `prefer be\.Eq\(n\)`
	be.RequireThat(t, x != n, be.True())   // want `prefer be\.Ne\(n\)`
}

func composites(t be.TestingT, xs []int, x int) {
	be.AssertThat(t, xs, be.Not(be.Nil()))         // want `prefer be\.NotNil\(\) over be\.Not\(be\.Nil\(\)\)`
	be.AssertThat(t, xs, be.HaveLength(0))         // want `prefer be\.Empty\(\) over be\.HaveLength\(0\)`
	be.AssertThat(t, xs, be.Not(be.HaveLength(0))) // want `prefer be\.NotEmpty\(\) over be\.Not\(be\.HaveLength\(0\)\)`
	be.AssertThat(t, xs, be.Not(be.Empty()))       // want `prefer be\.NotEmpty\(\) over be\.Not\(be\.Empty\(\)\)`
	be.AssertThat(t, x, be.Not(be.Eq(0)))          // want `prefer be\.NonZero\(\) over be\.Not\(be\.Eq\(0\)\)`
	be.Expect(t, xs).To(be.Not(be.Nil()))          // want `prefer be\.NotNil\(\)`

	// not zero-valued / not composable — not flagged
	be.AssertThat(t, xs, be.Not(be.Eq(x)))
	be.AssertThat(t, xs, be.HaveLength(3))
}
