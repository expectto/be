package belint

// plugin.go registers belint with golangci-lint's module plugin system — the
// lint entry point most repos (and their CI, and their coding agents) actually
// run. See README.md for the .custom-gcl.yml / .golangci.yml wiring.

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("belint", newPlugin)
}

type plugin struct{}

func newPlugin(any) (register.LinterPlugin, error) { return &plugin{}, nil }

func (*plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{Analyzer}, nil
}

func (*plugin) GetLoadMode() string { return register.LoadModeTypesInfo }
