module github.com/expectto/be

go 1.25

require (
	github.com/IGLOU-EU/go-wildcard v1.0.3 // latest
	github.com/amberpixels/k1 v0.1.4 // latest
	github.com/golang-jwt/jwt/v5 v5.3.0 // latest
	github.com/onsi/ginkgo/v2 v2.27.3 // latest
	github.com/onsi/gomega v1.38.3 // latest
	github.com/stretchr/testify v1.11.1
	go.uber.org/mock v0.6.0 // latest
	golang.org/x/text v0.32.0 // latest
)

// TODO: testify / ginkgo should be moved into inner "drivers" package
//       So the `be` itself is pure matchers, and inside user can select testify or ginkgo/gomega

require (
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20250403155104-27863c87afa6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/mod v0.30.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/tools v0.39.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
