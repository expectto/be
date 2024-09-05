# be_string
--
    import "."

Package be_string provides Be matchers for string-related assertions.

## Usage

```go
var V = psi_matchers.V
```

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
func MatchTemplate(template string, values ...*psi_matchers.Value) types.BeMatcher
```
MatchTemplate succeeds if actual matches given template pattern. Provided
template must have `{{Field}}` placeholders. Each distinct placeholder from
template requires a var to be passed in list of `vars`. Value (V) can be a raw
value or a matcher

E.g.

    Expect(someString).To(be_string.MatchTemplate("Hello {{Name}}. Your number is {{Number}}", be_string.V("Name", "John"), be_string.V("Number", 3)))
    Expect(someString).To(be_string.MatchTemplate("Hello {{Name}}. Good bye, {{Name}}.", be_string.V("Name", be_string.Titled()))

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

#### func  Only

```go
func Only(option StringOption) types.BeMatcher
```
Only succeeds if actual is a string containing only characters described by
given options Only() defaults to empty string matching Only(Alpha|Numeric)
succeeds if string contains only from alphabetic and numeric characters
Available options are: Alpha, Numeric, Whitespace, Dots, Punctuation,
SpecialCharacters TODO: special-characters are not supported yet

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
