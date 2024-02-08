# be_strings
--
    import "github.com/expectto/be/be_strings"

Package be_strings provides Be matchers for string-related assertions.

## Usage

#### func  Alpha

```go
func Alpha() types.BeMatcher
```
Alpha succeeds if actual is a string containing only alphabetical characters.
Actual must be a string-like value (can be adjusted via SetStringFormat method).

#### func  AlphaNumeric

```go
func AlphaNumeric() types.BeMatcher
```
AlphaNumeric succeeds if actual is a string containing only alphanumeric
characters. Actual must be a string-like value (can be adjusted via
SetStringFormat method). As Numeric() matcher is considered to match on
integers, AlphaNumeric() doesn't match on dots So, consider
AlphaNumericWithDots() then

#### func  AlphaNumericWhitespace

```go
func AlphaNumericWhitespace() types.BeMatcher
```
AlphaNumericWhitespace succeeds if actual is a string containing only
alphanumeric characters and whitespace. Actual must be a string-like value (can
be adjusted via SetStringFormat method).

#### func  AlphaNumericWhitespaceWithPunctuation

```go
func AlphaNumericWhitespaceWithPunctuation() types.BeMatcher
```
AlphaNumericWhitespaceWithPunctuation succeeds if actual is a string containing
alphanumeric, whitespace and punctuation. Actual must be a string-like value
(can be adjusted via SetStringFormat method).

#### func  AlphaNumericWithDots

```go
func AlphaNumericWithDots() types.BeMatcher
```
AlphaNumericWithDots succeeds if actual is a string containing only alphanumeric
characters and dots. Actual must be a string-like value (can be adjusted via
SetStringFormat method).

#### func  AlphaNumericWithPunctuation

```go
func AlphaNumericWithPunctuation() types.BeMatcher
```
AlphaNumericWithPunctuation succeeds if actual is a string containing only
alphanumeric characters and punctuation. Actual must be a string-like value (can
be adjusted via SetStringFormat method).

#### func  AlphaWhitespace

```go
func AlphaWhitespace() types.BeMatcher
```
AlphaWhitespace succeeds if actual is a string containing only alphabetical and
whitespace characters. Actual must be a string-like value (can be adjusted via
SetStringFormat method).

#### func  AlphaWhitespaceWithPunctuation

```go
func AlphaWhitespaceWithPunctuation() types.BeMatcher
```
AlphaWhitespaceWithPunctuation succeeds if actual is a string containing only
alphabetical characters with whitespace and punctuation.. Actual must be a
string-like value (can be adjusted via SetStringFormat method).

#### func  AlphaWithPunctuation

```go
func AlphaWithPunctuation() types.BeMatcher
```
AlphaWithPunctuation succeeds if actual is a string containing only alphabetical
characters. Actual must be a string-like value (can be adjusted via
SetStringFormat method).

#### func  ContainingCharacters

```go
func ContainingCharacters(characters string) types.BeMatcher
```
ContainingCharacters succeeds if actual is a string containing all characters
from a given set

#### func  ContainingOnlyCharacters

```go
func ContainingOnlyCharacters(characters string) types.BeMatcher
```
ContainingOnlyCharacters succeeds if actual is a string containing only
characters from a given set

#### func  ContainingSubstring

```go
func ContainingSubstring(substr string) types.BeMatcher
```
ContainingSubstring succeeds if actual is a string containing only characters
from a given set

#### func  EmptyString

```go
func EmptyString() types.BeMatcher
```
EmptyString succeeds if actual is an empty string. Actual must be a string-like
value (can be adjusted via SetStringFormat method).

#### func  Float

```go
func Float() types.BeMatcher
```
Float succeeds if actual is a string representing a valid floating-point number.
Actual must be a string-like value (can be adjusted via SetStringFormat method).

#### func  LowerCaseOnly

```go
func LowerCaseOnly() types.BeMatcher
```
LowerCaseOnly succeeds if actual is a string containing only lowercase
characters. Actual must be a string-like value (can be adjusted via
SetStringFormat method).

#### func  MatchTemplate

```go
func MatchTemplate(template string, vars ...*V) types.BeMatcher
```
MatchTemplate succeeds if actual matches given template pattern. Provided
template must have `{{Field}}` placeholders. Each distinct placeholder from
template requires a var to be passed in list of `vars`. Var can be a raw value
or a matcher

E.g.

    Expect(someString).To(be_strings.MatchTemplate("Hello {{Name}}. Your number is {{Number}}", be_strings.Var("Name", "John"), be_strings.Var("Number", 3)))
    Expect(someString).To(be_strings.MatchTemplate("Hello {{Name}}. Good bye, {{Name}}.", be_strings.Var("Name", be_strings.Titled()))

#### func  MatchWildcard

```go
func MatchWildcard(pattern string) types.BeMatcher
```
MatchWildcard succeeds if actual matches given wildcard pattern. Actual must be
a string-like value (can be adjusted via SetStringFormat method).

#### func  NonEmptyString

```go
func NonEmptyString() types.BeMatcher
```
NonEmptyString succeeds if actual is not an empty string. Actual must be a
string-like value (can be adjusted via SetStringFormat method).

#### func  Numeric

```go
func Numeric() types.BeMatcher
```
Numeric succeeds if actual is a string representing a valid numeric integer.
Actual must be a string-like value (can be adjusted via SetStringFormat method).

#### func  NumericWhitespace

```go
func NumericWhitespace() types.BeMatcher
```
NumericWhitespace succeeds if actual is a string containing only number
characters and whitespace. Actual must be a string-like value (can be adjusted
via SetStringFormat method).

#### func  Titled

```go
func Titled(languageArg ...language.Tag) types.BeMatcher
```
Titled succeeds if actual is a string with the first letter of each word
capitalized. Actual must be a string-like value (can be adjusted via
SetStringFormat method).

#### func  UpperCaseOnly

```go
func UpperCaseOnly() types.BeMatcher
```
UpperCaseOnly succeeds if actual is a string containing only uppercase
characters. Actual must be a string-like value (can be adjusted via
SetStringFormat method).

#### func  ValidEmail

```go
func ValidEmail() types.BeMatcher
```
ValidEmail succeeds if actual is a valid email. Actual must be a string-like
value (can be adjusted via SetStringFormat method).

#### func  Whitespace

```go
func Whitespace() types.BeMatcher
```
Whitespace succeeds if actual is a string containing only whitespace characters.
Actual must be a string-like value (can be adjusted via SetStringFormat method).

#### type V

```go
type V struct {
	Name    string
	Matcher types.BeMatcher
}
```


#### func  Var

```go
func Var(name string, matching any) *V
```
Var creates a var used for replacing placeholders for templates in
`MatchTemplate`
