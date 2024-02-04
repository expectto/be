# be_http
--
    import "github.com/expectto/be/be_http"

Package be_http provides matchers for url.Request TODO: more detailed
documentation here is required

## Usage

```go
var (
	GET     = func() types.BeMatcher { return HavingMethod(http.MethodGet) }
	HEAD    = func() types.BeMatcher { return HavingMethod(http.MethodHead) }
	POST    = func() types.BeMatcher { return HavingMethod(http.MethodPost) }
	PUT     = func() types.BeMatcher { return HavingMethod(http.MethodPut) }
	PATCH   = func() types.BeMatcher { return HavingMethod(http.MethodPatch) }
	DELETE  = func() types.BeMatcher { return HavingMethod(http.MethodDelete) }
	OPTIONS = func() types.BeMatcher { return HavingMethod(http.MethodOptions) }
	CONNECT = func() types.BeMatcher { return HavingMethod(http.MethodConnect) }
	TRACE   = func() types.BeMatcher { return HavingMethod(http.MethodTrace) }
)
```
HavingMethod: Syntactic sugar

#### func  HavingBody

```go
func HavingBody(args ...any) types.BeMatcher
```
HavingBody succeeds if the actual value is a *http.Request and its body matches
the provided arguments. Note: The body is not re-streamed, so it's not available
after matching.

#### func  HavingHeader

```go
func HavingHeader(key string, args ...any) types.BeMatcher
```
HavingHeader matches requests that have header with a given key. Key is a string
key for a header, args can be nil or len(args)==1. Note: Golang's http.Header is
`map[string][]string`, and matching is done on the FIRST value of the header in
case if you have multiple-valued header that needs to be matched, use
HavingHeaders() instead

These are scenarios that can be handled here: (1) If no args are given, it
simply matches a request with existed header by key. (2) If len(args) == 1 &&
args[0] is a stringish, it matches a request with header `Key: Args[0]` (3) if
len(args) == 1 && args[0] is not stringish, it is considered to be matcher for
header's value Examples: - HavingHeader("X-Header") matches request with
non-empty X-Header header - HavingHeader("X-Header", "X-Value") matches request
with X-Header: X-Value - HavingHeader("X-Header", HavePrefix("Bearer "))
matchers request with header(X-Header)'s value matching given HavePrefix matcher

#### func  HavingHeaders

```go
func HavingHeaders(key string, args ...any) types.BeMatcher
```
HavingHeaders matches requests that have header with a given key. Key is a
string key for a header, args can be nil or len(args)==1. Note: Matching is done
on the list of header values. In case if you have single-valued header that
needs to be matched, use HavingHeader() instead

These are scenarios that can be handled here: (1) If no args are given, it
simply matches a request with existed header by key. (2) If len(args) == 1 &&
args[0] is a stringish, it matches a request with header `Key: Args[0]` (3) if
len(args) == 1 && args[0] is not stringish, it is considered to be matcher for
header's value Examples: - HavingHeader("X-Header") matches request with
non-empty X-Header header - HavingHeader("X-Header", "X-Value") matches request
with X-Header: X-Value - HavingHeader("X-Header", Dive(HavePrefix("Foo ")))
matchers request with multiple X-Header values, each of them having Foo prefix

#### func  HavingHost

```go
func HavingHost(args ...any) types.BeMatcher
```
HavingHost succeeds if the actual value is a *http.Request and its Host matches
the provided arguments.

#### func  HavingMethod

```go
func HavingMethod(args ...any) types.BeMatcher
```
HavingMethod succeeds if the actual value is a *http.Request and its HTTP method
matches the provided arguments.

#### func  HavingProto

```go
func HavingProto(args ...any) types.BeMatcher
```
HavingProto succeeds if the actual value is a *http.Request and its Proto
matches the provided arguments.

#### func  HavingURL

```go
func HavingURL(args ...any) types.BeMatcher
```
HavingURL succeeds if the actual value is a *http.Request and its URL matches
the provided arguments.

#### func  Request

```go
func Request(args ...any) types.BeMatcher
```
Request matches an actual value to be a valid *http.Request corresponding to
given inputs. Possible inputs: 1. Nil args -> so actual value MUST be any valid
*http.Request. 2. Single arg <string>. Actual value MUST be a *http.Request,
whose .URL.String() is compared against args[0]. 3. List of Omega/Gomock/Psi
matchers, that are applied to *http.Request object.

    - Supports matching http.Request properties like method, URL, body, host, proto, and headers.
    - Additional arguments can be used for matching specific headers, e.g., WithHeader("Content-Type", "application/json").
