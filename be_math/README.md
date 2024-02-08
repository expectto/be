# be_math
--
    import "github.com/expectto/be/be_math"

Package be_math provides Be matchers for mathematical operations

## Usage

#### func  Approx

```go
func Approx(compareTo, threshold any) types.BeMatcher
```
Approx succeeds if actual is numerically approximately equal to the passed-in
value within the specified threshold.

#### func  ApproxZero

```go
func ApproxZero() types.BeMatcher
```
ApproxZero succeeds if actual is numerically approximately equal to zero Any
type of int/float will work for comparison.

#### func  DivisibleBy

```go
func DivisibleBy(divisor any) types.BeMatcher
```
DivisibleBy succeeds if actual is numerically divisible by the passed-in value.

#### func  Even

```go
func Even() types.BeMatcher
```
Even succeeds if actual is an even numeric value.

#### func  GreaterThan

```go
func GreaterThan(arg any) types.BeMatcher
```
GreaterThan succeeds if actual is numerically greater than the passed-in value.

#### func  GreaterThanEqual

```go
func GreaterThanEqual(arg any) types.BeMatcher
```
GreaterThanEqual succeeds if actual is numerically greater than or equal to the
passed-in value.

#### func  Gt

```go
func Gt(arg any) types.BeMatcher
```
Gt is an alias for GreaterThan, succeeding if actual is numerically greater than
the passed-in value.

#### func  Gte

```go
func Gte(arg any) types.BeMatcher
```
Gte is an alias for GreaterThanEqual, succeeding if actual is numerically
greater than or equal to the passed-in value.

#### func  InRange

```go
func InRange(from any, fromInclusive bool, until any, untilInclusive bool) types.BeMatcher
```
InRange succeeds if actual is numerically within the specified range. The range
is defined by the 'from' and 'until' values, and inclusivity is determined by
the 'fromInclusive' and 'untilInclusive' flags.

#### func  Integral

```go
func Integral() types.BeMatcher
```
Integral succeeds if actual is an integral float, meaning it has zero decimal
places. This matcher checks if the numeric value has no fractional component.

#### func  LessThan

```go
func LessThan(arg any) types.BeMatcher
```
LessThan succeeds if actual is numerically less than the passed-in value.

#### func  LessThanEqual

```go
func LessThanEqual(arg any) types.BeMatcher
```
LessThanEqual succeeds if actual is numerically less than or equal to the
passed-in value.

#### func  Lt

```go
func Lt(arg any) types.BeMatcher
```
Lt is an alias for LessThan, succeeding if actual is numerically less than the
passed-in value.

#### func  Lte

```go
func Lte(arg any) types.BeMatcher
```
Lte is an alias for LessThanEqual, succeeding if actual is numerically less than
or equal to the passed-in value.

#### func  Negative

```go
func Negative() types.BeMatcher
```
Negative succeeds if actual is a negative numeric value.

#### func  Odd

```go
func Odd() types.BeMatcher
```
Odd succeeds if actual is an odd numeric value.

#### func  Positive

```go
func Positive() types.BeMatcher
```
Positive succeeds if actual is a positive numeric value.

#### func  Zero

```go
func Zero() types.BeMatcher
```
Zero succeeds if actual is numerically equal to zero. Any type of int/float will
work for comparison.
