package testify

import (
	"github.com/expectto/be/types"
	"github.com/stretchr/testify/mock"
)

// Mock adapts a be matcher into a testify/mock argument matcher, so be matchers
// can match mock call arguments. It works with hand-written testify mocks and
// mockery-generated mocks alike:
//
//	svc.On("Do", betestify.Mock(be_math.GreaterThan(10))).Return("ok")
func Mock(matcher types.BeMatcher) any {
	return mock.MatchedBy(func(actual any) bool {
		ok, _ := matcher.Match(actual)
		return ok
	})
}
