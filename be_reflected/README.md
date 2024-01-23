# be_reflected
--
    import "github.com/expectto/be/be_reflected"

Package be_reflected provides Be matchers that use reflection, enabling
expressive assertions on values' reflect kinds and types.

## Usage

#### func  AsBytes

```go
func AsBytes() types.BeMatcher
```
AsBytes succeeds if actual is assignable to a slice of bytes ([]byte).

#### func  AsChan

```go
func AsChan() types.BeMatcher
```
AsChan succeeds if actual is of kind reflect.Chan.

#### func  AsFinalPointer

```go
func AsFinalPointer() types.BeMatcher
```
AsFinalPointer succeeds if the actual value is a final pointer, meaning it's a
pointer to a non-pointer type.

#### func  AsFloatish

```go
func AsFloatish() types.BeMatcher
```
AsFloatish succeeds if actual is a numeric value that represents a
floating-point value.

#### func  AsFloatishString

```go
func AsFloatishString() types.BeMatcher
```
AsFloatishString succeeds if actual is a string that can be parsed into a valid
floating-point value.

#### func  AsFunc

```go
func AsFunc() types.BeMatcher
```
AsFunc succeeds if actual is of kind reflect.Func.

#### func  AsIntish

```go
func AsIntish() types.BeMatcher
```
AsIntish succeeds if actual is a numeric value that represents an integer (from
reflect.Int up to reflect.Uint64).

#### func  AsIntishString

```go
func AsIntishString() types.BeMatcher
```
AsIntishString succeeds if actual is a string that can be parsed into a valid
integer value.

#### func  AsKind

```go
func AsKind(args ...any) types.BeMatcher
```
AsKind succeeds if actual is assignable to any of the specified kinds or matches
the provided matchers.

#### func  AsMap

```go
func AsMap() types.BeMatcher
```
AsMap succeeds if actual is of kind reflect.Map.

#### func  AsNumeric

```go
func AsNumeric() types.BeMatcher
```
AsNumeric succeeds if actual is a numeric value, supporting various integer
kinds: reflect.Int, ... reflect.Int64, and floating-point kinds:
reflect.Float32, reflect.Float64

#### func  AsNumericString

```go
func AsNumericString() types.BeMatcher
```
AsNumericString succeeds if actual is a string that can be parsed into a valid
numeric value.

#### func  AsObject

```go
func AsObject() types.BeMatcher
```
AsObject is more specific than AsMap. It checks if the given `actual` value is a
map with string keys and values of any type. This is particularly useful in the
context of BeJson matcher, where the term 'Object' aligns with JSON notation.

#### func  AsObjects

```go
func AsObjects() types.BeMatcher
```

#### func  AsPointer

```go
func AsPointer() types.BeMatcher
```
AsPointer succeeds if the actual value is a pointer.

#### func  AsPointerToMap

```go
func AsPointerToMap() types.BeMatcher
```
AsPointerToMap succeeds if actual is a pointer to a map.

#### func  AsPointerToObject

```go
func AsPointerToObject() types.BeMatcher
```
AsPointerToObject succeeds if actual is a pointer to a value that matches
AsObject after applying dereference.

#### func  AsPointerToSlice

```go
func AsPointerToSlice() types.BeMatcher
```
AsPointerToSlice succeeds if actual is a pointer to a slice.

#### func  AsPointerToStruct

```go
func AsPointerToStruct() types.BeMatcher
```
AsPointerToStruct succeeds if actual is a pointer to a struct.

#### func  AsReader

```go
func AsReader() types.BeMatcher
```
AsReader succeeds if actual implements the io.Reader interface.

#### func  AsSlice

```go
func AsSlice() types.BeMatcher
```
AsSlice succeeds if actual is of kind reflect.Slice.

#### func  AsSliceOf

```go
func AsSliceOf[T any]() types.BeMatcher
```
AsSliceOf succeeds if actual is of kind reflect.Slice and each element of the
slice is assignable to the specified type T.

#### func  AsString

```go
func AsString() types.BeMatcher
```
AsString succeeds if actual is of kind reflect.String.

#### func  AsStringer

```go
func AsStringer() types.BeMatcher
```
AsStringer succeeds if actual implements the fmt.Stringer interface.

#### func  AsStruct

```go
func AsStruct() types.BeMatcher
```
AsStruct succeeds if actual is of kind reflect.Struct.

#### func  AssignableTo

```go
func AssignableTo[T any]() types.BeMatcher
```
AssignableTo succeeds if actual is assignable to the specified type T.

#### func  Implementing

```go
func Implementing[T any]() types.BeMatcher
```
Implementing succeeds if actual implements the specified interface type T.
