# be_json
--
    import "."

Package be_json provides Be matchers for expressive assertions on JSON TODO:
more detailed explanation what is considered to be JSON here

## Usage

#### func  HaveKeyValue

```go
func HaveKeyValue(key any, args ...any) types.BeMatcher
```
HaveKeyValue is a facade to gomega.HaveKey & gomega.HaveKeyWithValue

#### func  Matcher

```go
func Matcher(args ...any) types.BeMatcher
```
Matcher is a JSON matcher. "JSON" here means a []byte with JSON data in it By
default several input types are available: string(*) / []byte(*), fmt.Stringer,
io.Reader

    - custom string-based or []byte-based types are available as well

To make it stricter and to specify which format JSON we should expect, you must
pass one of transforms as first argument:

    - JsonAsBytes/ JsonAsString / JsonAsStringer  / JsonAsReader (for string-like representation)
    - JsonAsObject / JsonAsObjects (for map[string]any representation)

#### type JsonInputType

```go
type JsonInputType uint32
```


```go
const (
	JsonAsBytes JsonInputType = 1 << iota
	JsonAsString
	JsonAsStringer
	JsonAsReader
	JsonAsObject
	JsonAsObjects
	JsonAsStruct
)
```
