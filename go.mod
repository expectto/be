module github.com/expectto/be

go 1.23.6

require (
	github.com/IGLOU-EU/go-wildcard v1.0.3 // latest
	github.com/golang-jwt/jwt/v5 v5.2.1 // latest
	github.com/onsi/ginkgo/v2 v2.22.2 // latest
	github.com/onsi/gomega v1.36.2 // latest
	github.com/stretchr/testify v1.9.0
	go.uber.org/mock v0.5.0 // latest
	golang.org/x/text v0.22.0 // latest
)

// TODO: testify / ginkgo should be moved into inner "drivers" package
//       So the `be` itself is pure matchers, and inside user can select testify or ginkgo/gomega

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/pprof v0.0.0-20241210010833-40e02aabc2ad // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/tools v0.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
