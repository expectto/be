# be_http
--
    import "github.com/expectto/be/be_http"

Package be_http provides matchers for url.Request todo: more detailed
documentation here is required

## Usage

#### func  DELETE

```go
func DELETE() types.BeMatcher
```
DELETE returns a matcher that succeeds if the actual *http.Request has a method
"DELETE".

#### func  GET

```go
func GET() types.BeMatcher
```
GET returns a matcher that succeeds if the actual *http.Request has a method
"GET".

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
HavingHeader matches requests that have header with a given key. (1) If no args
are given, it simply matches a request with existed header by key. (2) If
len(args) == 1 && args[0] is a stringish, it matches a request with header `Key:
Args[0]` (3) if len(args) == 1 && args[0] is not stringish, it is considered to
be matcher for header's value Examples: - HavingHeader("X-Header") matches
request with non-empty X-Header header - HavingHeader("X-Header", "X-Value")
matches request with X-Header: X-Value - HavingHeader("X-Header",
HavePrefix("Bearer ")) matchers request with header(X-Header)'s value matching
given HavePrefix matcher - todo: support multiple header values todo: fixme I'm
ugly for now

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

#### func  OPTIONS

```go
func OPTIONS() types.BeMatcher
```
OPTIONS returns a matcher that succeeds if the actual *http.Request has a method
"OPTIONS".

#### func  PATCH

```go
func PATCH() types.BeMatcher
```
PATCH returns a matcher that succeeds if the actual *http.Request has a method
"PATCH".

#### func  POST

```go
func POST() types.BeMatcher
```
POST returns a matcher that succeeds if the actual *http.Request has a method
"POST".

#### func  PUT

```go
func PUT() types.BeMatcher
```
PUT returns a matcher that succeeds if the actual *http.Request has a method
"PUT".

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
