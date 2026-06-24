module github.com/expectto/be/x/testify

go 1.26

// This driver is a separate module so the core `github.com/expectto/be` module
// never lists testify (or any test framework) as a dependency. Users who want
// the testify adapter opt in by importing this module explicitly.

require (
	github.com/expectto/be v0.2.4
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
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Local development against the in-repo core. When the core is tagged, bump the
// require above; external consumers ignore this replace and use the tagged core.
replace github.com/expectto/be => ../..
