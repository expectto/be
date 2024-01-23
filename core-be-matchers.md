# be
--
    import "github.com/expectto/be"


## Usage

```go
var Ctx = be_ctx.Ctx
```
Ctx is an alias for be_ctx.Ctx

```go
var HttpRequest = be_http.Request
```
HttpRequest is an alias for be_http.Request matcher

```go
var JwtToken = be_jwt.Token
```
JwtToken is an alias for be_jwt.Token matcher

```go
var StringAsTemplate = be_strings.MatchTemplate
```
StringAsTemplate is an alias for be_strings.MatchTemplate matcher

```go
var URL = be_url.URL
```
URL is an alias for be_url.URL matcher

#### func  All

```go
func All(ms ...any) types.BeMatcher
```
All is like gomega.And()

#### func  Always

```go
func Always() types.BeMatcher
```
Always does always match

#### func  Any

```go
func Any(ms ...any) types.BeMatcher
```
Any is like gomega.Or()

#### func  Eq

```go
func Eq(expected any) types.BeMatcher
```
Eq is like gomega.Equal()

#### func  HaveLength

```go
func HaveLength(args ...any) types.BeMatcher
```
HaveLength is like gomega.HaveLen() HaveLength succeeds if the actual value has
a length that matches the provided conditions. It accepts either a count value
or one or more Gomega matchers to specify the desired length conditions.

#### func  Never

```go
func Never(err error) types.BeMatcher
```
Never does never succeed (does always fail)

#### func  Not

```go
func Not(expected any) types.BeMatcher
```
Not is like gomega.Not()
