<p align="center">
  <img src="logo.svg" alt="be_struct" width="366">
</p>

<div align="center">

Part of [`expectto/be`](../README.md) - composable test matchers for Go.

</div>

---

```go
import "github.com/expectto/be/be_struct"
```

Package be_struct provides Be matchers on struct fields.

## Usage

#### func  HavingField

```go
func HavingField[StructT any](fieldName string, expectedValue ...any) types.BeMatcher
```
HavingField succeeds if the actual value is a struct of exactly type StructT and
it has a field with the given name. If an expected value (raw value or matcher)
is provided, the field's value must match it too.

Example:

    Expect(result).To(be_struct.HavingField[TestStruct]("Field1", "hello1"))
    Expect(result).To(be_struct.HavingField[TestStruct]("Field2"))

This is the generic, compile-time-typed variant: it pins the struct type at the
call site. For everyday assertions prefer the root be.HaveField — it needs no
type parameter, and supports dotted paths and "Method()" specs. The naming
wobble (HaveField vs HavingField) is deliberate: it separates the default
matcher from this typed variant.
