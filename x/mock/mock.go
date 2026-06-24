// Package bemock adapts expectto/be matchers for use as testify/mockery mock
// argument matchers.
//
// It lives in its own module (github.com/expectto/be/x/mock) so the core `be`
// module never depends on testify. Assertions don't live here — use the core
// be.AssertThat / be.RequireThat (or be.Expect/Require) instead. Gomock needs no
// adapter either (every be matcher already satisfies gomock.Matcher); this
// module's sole job is the testify/mockery adapter below.
package bemock

import (
	"github.com/expectto/be/types"
	"github.com/stretchr/testify/mock"
)

// MatchedBy adapts a be matcher into a testify/mock argument matcher, the matcher
// equivalent of testify's own mock.MatchedBy(func). It works with hand-written
// testify mocks and mockery-generated mocks alike:
//
//	svc.On("Do", bemock.MatchedBy(be_math.GreaterThan(10))).Return("ok")
func MatchedBy(matcher types.BeMatcher) any {
	return mock.MatchedBy(func(actual any) bool {
		ok, _ := matcher.Match(actual)
		return ok
	})
}
