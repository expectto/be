# be_ctx
--
    import "."

Package be_ctx provides Be matchers on context.Context

## Usage

#### func  Ctx

```go
func Ctx(args ...any) types.BeMatcher
```
Ctx succeeds if the actual value is a context.Context. If no arguments are
provided, it matches any context.Context. Otherwise, it uses the Psi matcher to
match the provided arguments against the actual context's values.

#### func  CtxWithDeadline

```go
func CtxWithDeadline(deadline any) types.BeMatcher
```
CtxWithDeadline succeeds if the actual value is a context.Context and its
deadline matches the provided deadline.

#### func  CtxWithError

```go
func CtxWithError(err any) types.BeMatcher
```
CtxWithError succeeds if the actual value is a context.Context and its error
matches the provided error value.

#### func  CtxWithValue

```go
func CtxWithValue(key any, vs ...any) types.BeMatcher
```
CtxWithValue succeeds if the actual value is a context.Context and contains a
key-value pair where the key matches the provided key and the value matches the
provided arguments using any other matchers.
