# be_url
--
    import "github.com/expectto/be/be_url"

Package be_url provides Be matchers on url.URL

## Usage

```go
var TransformSchemelessUrlFromString = func(rawURL string) (*url.URL, error) {
	result, err := url.Parse(rawURL)
	if err == nil && result.Scheme == "" {
		result, err = url.Parse("http://" + rawURL)
		if err == nil {
			result.Scheme = ""
		}
	}
	return result, err
}
```
TransformSchemelessUrlFromString returns string->*url.Url transform It allows
string to be a scheme-less url

```go
var TransformUrlFromString = url.Parse
```
TransformUrlFromString returns string->*url.Url transform

#### func  HavingHost

```go
func HavingHost(args ...any) types.BeMatcher
```
HavingHost succeeds if the actual value is a *url.URL and its Host matches the
provided one (via direct value or matchers)

#### func  HavingHostname

```go
func HavingHostname(args ...any) types.BeMatcher
```
HavingHostname succeeds if the actual value is a *url.URL and its Hostname
matches the provided one (via direct value or matchers)

#### func  HavingMultipleSearchParam

```go
func HavingMultipleSearchParam(searchParamName string, args ...any) types.BeMatcher
```
HavingMultipleSearchParam succeeds if the actual value is a *url.URL and its
specified search parameter (all its values via slice) matches the provided
arguments.

#### func  HavingPassword

```go
func HavingPassword(args ...any) types.BeMatcher
```
HavingPassword succeeds if the actual value is a *url.URL and its Password
matches the provided one.

#### func  HavingPath

```go
func HavingPath(args ...any) types.BeMatcher
```
HavingPath succeeds if the actual value is a *url.URL and its Path matches the
given one.

#### func  HavingPort

```go
func HavingPort(args ...any) types.BeMatcher
```
HavingPort succeeds if the actual value is a *url.URL and its Port matches the
provided one.

#### func  HavingRawQuery

```go
func HavingRawQuery(args ...any) types.BeMatcher
```
HavingRawQuery succeeds if the actual value is a *url.URL and its RawQuery
matches the given one.

#### func  HavingScheme

```go
func HavingScheme(args ...any) types.BeMatcher
```
HavingScheme succeeds if the actual value is a *url.URL and its Scheme matches
the provided one (via direct value or matchers)

#### func  HavingSearchParam

```go
func HavingSearchParam(searchParamName string, args ...any) types.BeMatcher
```
HavingSearchParam succeeds if the actual value is a *url.URL and its specified
search parameter matches the provided arguments.

#### func  HavingUserinfo

```go
func HavingUserinfo(args ...any) types.BeMatcher
```
HavingUserinfo succeeds if the actual value is a *url.URL and its User.String()
matches the provided one.

#### func  HavingUsername

```go
func HavingUsername(args ...any) types.BeMatcher
```
HavingUsername succeeds if the actual value is a *url.URL and its Username
matches the provided one.

#### func  NotHavingPort

```go
func NotHavingPort(args ...any) types.BeMatcher
```
NotHavingPort succeeds if the actual value is a *url.URL and its Port does not
match the given one. Example: `Expect(u).To(NotHavingPort())` matches port-less
url

#### func  NotHavingScheme

```go
func NotHavingScheme(args ...any) types.BeMatcher
```
NotHavingScheme succeeds if the actual value is a *url.URL and its Scheme
negatively matches given value Example: `Expect(u).To(NotHavingScheme())`
matches url without a scheme

#### func  URL

```go
func URL(args ...any) types.BeMatcher
```
URL matches actual value to be a valid URL corresponding to given inputs
Possible inputs: 1. Nil args -> so actual value MUST be any valid *url.URL 2.
Single arg <string>. Actual value MUST be a *url.URL, whose .String() compared
against args[0] 3. Single arg <*url.Url>. Actual value MUST be a *url.URL, whose
.String() compared against args[0].String() 4. List of Omega/Gomock/Psi
matchers, that are applied to *url.URL object

    - TransformUrlFromString() transform can be given as first argument, so string->*url.URL transform is applied

#### func  WithHttp

```go
func WithHttp() types.BeMatcher
```
WithHttp succeeds if the actual value is a *url.URL and its scheme is "http".

#### func  WithHttps

```go
func WithHttps() types.BeMatcher
```
WithHttps succeeds if the actual value is a *url.URL and its scheme is "https".
