module github.com/expectto/be/x/testify

go 1.26

// This driver is a separate module so the core `github.com/expectto/be` module
// never lists testify (or any test framework) as a dependency. Users who want
// the testify adapter opt in by importing this module explicitly.

require (
	github.com/expectto/be v1.0.0-rc.3
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/IGLOU-EU/go-wildcard v1.0.3 // indirect
	github.com/amberpixels/k1 v0.1.6 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/onsi/gomega v1.42.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Local development resolves the core from the in-repo source. External consumers
// ignore this replace and use the version pinned in the require block above.
replace github.com/expectto/be => ../..
