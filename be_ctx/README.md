<p align="center">
  <img src="logo.svg" alt="be_ctx" width="204">
</p>

<div align="center">

Part of [`expectto/be`](../README.md) - composable test matchers for Go.

</div>

---

```go
import "github.com/expectto/be/be_ctx"
```

Package be_ctx provides Be matchers on context.Context

## Usage

#### func  Ctx

```go
func Ctx(args ...any) types.BeMatcher
```
Ctx succeeds if the actual value is a context.Context. If no arguments are
provided, it matches any context.Context. Otherwise it first enforces that the
actual value is a context.Context and then applies the provided matchers (e.g.
CtxWithValue/CtxWithDeadline) to it.

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
