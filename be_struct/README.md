# be_struct
--
    import "."


## Usage

#### func  HavingField

```go
func HavingField[StructT any](fieldName string, expectedValue ...any) types.BeMatcher
```
HavingField succeeds if the actual value is a struct and it has a field with the
given name. If an expected value is provided, it also succeeds if the actual
value's field has the same value.

Example:

    Expect(result).To(be_structs.HavingField[TestStruct]("Field1", "hello1"))
    Expect(result).To(be_structs.HavingField[TestStruct]("Field2"))
