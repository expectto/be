package belint_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/expectto/be/x/belint"
)

func TestBelint(t *testing.T) {
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), belint.Analyzer, "a")
}
