# be_jwt
--
    import "github.com/expectto/be/be_jwt"

Package be_jwt provides Be matchers for handling JSON Web Tokens (JWT). It
includes matchers for transforming and validating JWT tokens. Matchers
corresponds to specific golang jwt implementation:
https://github.com/golang-jwt/jwt/v5

## Usage

```go
var TransformJwtFromString = func(input string) any {
	p := jwt.NewParser()

	t, parts, err := p.ParseUnverified(input, jwt.MapClaims{})
	if err != nil {
		return NewTransformError(err, input)
	}

	t.Signature, err = p.DecodeSegment(parts[2])
	if err != nil {
		return NewTransformError(fmt.Errorf("corrupted signature part: %w", err), input)
	}

	return t
}
```
TransformJwtFromString is a transform function (string->*jwt.Token) without a
secret. It parses the input string as a JWT and returns the resulting
*jwt.Token.

```go
var TransformSignedJwtFromString = func(secret string) func(string) any {
	return func(input string) any {
		parsed, err := jwt.Parse(input, func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		if err != nil {
			return NewTransformError(fmt.Errorf("to parse jwt token (with secret=%s): %w", secret, err), input)
		}

		return parsed
	}
}
```
TransformSignedJwtFromString returns a transform function (string->*jwt.Token)
for a given secret.

#### func  HavingClaims

```go
func HavingClaims(args ...any) types.BeMatcher
```
HavingClaims succeeds if the actual value is a JWT token and its claims match
the provided value or matchers.

#### func  HavingMethodAlg

```go
func HavingMethodAlg(args ...any) types.BeMatcher
```
HavingMethodAlg succeeds if the actual value is a JWT token and its method
algorithm match the provided value or matchers.

#### func  SignedVia

```go
func SignedVia(secret string) types.BeMatcher
```
SignedVia succeeds if the actual value is a valid and signed JWT token, verified
using the specified secret key. It's intended for matching against a secret-less
token and applying the secret only for this specific matching.

Example:

Token(TransformJwtFromString, SignedVia(secret)) // works similar to:
Token(TransformSignedJwtFromString(secret), Valid())

#### func  Token

```go
func Token(args ...any) types.BeMatcher
```
Token matches the actual value to be a valid *jwt.Token corresponding to given
inputs. Possible inputs: 1. No args -> the actual value MUST be any valid
*jwt.Token. 2. Single arg <string>. The actual value MUST be a *jwt.Token, whose
.String() is compared against args[0]. 3. Single arg <*jwt.Token>. The actual
value MUST be a *jwt.Token. 4. List of Omega/Gomock/Psi matchers that are
applied to *jwt.Token object.

    - TransformJwtFromString/TransformSignedJwtFromString(secret) transforms can be given as the first argument,
      so the string->*jwt.Token transform is applied.

#### func  Valid

```go
func Valid() types.BeMatcher
```
Valid succeeds if the actual value is a JWT token and it's valid
