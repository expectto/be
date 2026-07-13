package belint_test

import (
	"testing"

	"github.com/expectto/be/x/belint"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestBelint(t *testing.T) {
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), belint.Analyzer, "a")
}
