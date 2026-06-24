package testify_test

import (
	"testing"

	"github.com/expectto/be/be_math"
	betestify "github.com/expectto/be/x/testify"
	"github.com/stretchr/testify/mock"
)

// fakeService is a hand-written testify mock (mockery would generate the same shape).
type fakeService struct{ mock.Mock }

func (f *fakeService) Do(n int) string { return f.Called(n).String(0) }

func TestMockArgumentMatching(t *testing.T) {
	svc := &fakeService{}
	// The argument is matched by a be matcher rather than a literal value.
	svc.On("Do", betestify.Mock(be_math.GreaterThan(10))).Return("big")

	if got := svc.Do(42); got != "big" {
		t.Fatalf(`Do(42) = %q, want "big"`, got)
	}
	svc.AssertExpectations(t)
}
