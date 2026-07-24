// Command belint runs the belint analyzer standalone (with -fix support):
//
//	go run github.com/expectto/be/x/belint/cmd/belint@latest -fix ./...
//
// For day-to-day linting prefer the golangci-lint plugin registration — see
// the package README.
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/expectto/be/x/belint"
)

func main() {
	singlechecker.Main(belint.Analyzer)
}
