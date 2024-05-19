# be_time
--
    import "github.com/expectto/be/be_time"

Package be_time provides Be matchers on time.Time

## Usage

#### func  Approx

```go
func Approx(compareTo time.Time, threshold time.Duration) types.BeMatcher
```
Approx succeeds if actual time is approximately equal to the specified time
`compareTo` within the given time duration threshold.

#### func  EarlierThan

```go
func EarlierThan(compareTo time.Time) types.BeMatcher
```
EarlierThan succeeds if actual time is earlier than the specified time
`compareTo`.

#### func  EarlierThanEqual

```go
func EarlierThanEqual(compareTo time.Time) types.BeMatcher
```
EarlierThanEqual succeeds if actual time is earlier than or equal to the
specified time `compareTo`.

#### func  Eq

```go
func Eq(compareTo time.Time) types.BeMatcher
```
Eq succeeds if actual time is equal to the specified time `compareTo` with the
precision of one nanosecond.

#### func  IsDST

```go
func IsDST(compareTo time.Time) types.BeMatcher
```
IsDST checks if actual time is DST

#### func  LaterThan

```go
func LaterThan(compareTo time.Time) types.BeMatcher
```
LaterThan succeeds if actual time is later than the specified time `compareTo`.

#### func  LaterThanEqual

```go
func LaterThanEqual(compareTo time.Time) types.BeMatcher
```
LaterThanEqual succeeds if actual time is later than or equal to the specified
time `compareTo`.

#### func  SameDay

```go
func SameDay(compareTo time.Time) types.BeMatcher
```
SameDay checks if the .Day() component of the actual time matches the .Day()
component of the specified time `compareTo`. It only verifies the day
[1..30(31)], disregarding the month, year, etc.

#### func  SameExactDay

```go
func SameExactDay(compareTo time.Time) types.BeMatcher
```
SameExactDay succeeds if the actual time falls within the same day as the
specified time `compareTo`.

#### func  SameExactHour

```go
func SameExactHour(compareTo time.Time) types.BeMatcher
```
SameExactHour succeeds if the actual time falls within the same hour as the
specified time `compareTo`.

#### func  SameExactMilli

```go
func SameExactMilli(compareTo time.Time) types.BeMatcher
```
SameExactMilli succeeds if the actual time falls within the same millisecond as
the specified time `compareTo`.

#### func  SameExactMinute

```go
func SameExactMinute(compareTo time.Time) types.BeMatcher
```
SameExactMinute succeeds if the actual time falls within the same minute as the
specified time `compareTo`.

#### func  SameExactMonth

```go
func SameExactMonth(compareTo time.Time) types.BeMatcher
```
SameExactMonth succeeds if the actual time falls within the same month as the
specified time `compareTo`.

#### func  SameExactSecond

```go
func SameExactSecond(compareTo time.Time) types.BeMatcher
```
SameExactSecond succeeds if the actual time falls within the same second as the
specified time `compareTo`.

#### func  SameExactWeek

```go
func SameExactWeek(compareTo time.Time) types.BeMatcher
```
SameExactWeek succeeds if the actual time falls within the same ISO week as the
specified time `compareTo`.

#### func  SameExactWeekday

```go
func SameExactWeekday(compareTo time.Time) types.BeMatcher
```
SameExactWeekday succeeds if the weekday component of the actual time is equal
to the weekday component of the specified time `compareTo`.

#### func  SameHour

```go
func SameHour(compareTo time.Time) types.BeMatcher
```
SameHour checks if the .Hour() component of the actual time matches the .Hour()
component of the specified time `compareTo`. It only verifies the hour [0..59],
disregarding other components such as second, minute, day, etc.

#### func  SameMinute

```go
func SameMinute(compareTo time.Time) types.BeMatcher
```
SameMinute checks if the .Minute() component of the actual time matches the
.Minute() component of the specified time `compareTo`. It only verifies the
minute [0..59], disregarding other components such as second, hour, day, etc.

#### func  SameMonth

```go
func SameMonth(compareTo time.Time) types.BeMatcher
```
SameMonth checks if the .Month() component of the actual time matches the
.Month() component of the specified time `compareTo`. It only verifies the month
[1..12], disregarding the year.

#### func  SameOffset

```go
func SameOffset(compareTo time.Time) types.BeMatcher
```
SameOffset checks if actual time is the same timezone offset as specified time
`compareTo` Note: times can have different timezone names, but same offset, e.g.
America/New_York and Canada/Toronto

#### func  SameSecond

```go
func SameSecond(compareTo time.Time) types.BeMatcher
```
SameSecond checks if the .Second() component of the actual time matches the
.Second() component of the specified time `compareTo`. It only verifies the
second [0..59], disregarding other components such as minute, hour, day, etc.

#### func  SameTimezone

```go
func SameTimezone(compareTo time.Time) types.BeMatcher
```
SameTimezone checks if actual time is the same timezone as specified time
`compareTo`

#### func  SameWeek

```go
func SameWeek(compareTo time.Time) types.BeMatcher
```
SameWeek succeeds if the ISO week of the actual time are equal to the ISO week
and year components of the specified time `compareTo`. It only verifies the week
[1..53], disregarding of year. Note: use SameExactWeek to respect exact week of
exact year

#### func  SameYear

```go
func SameYear(compareTo time.Time) types.BeMatcher
```
SameYear succeeds if the year component of the actual time is equal to the year
component of the specified time `compareTo`.

#### func  SameYearDay

```go
func SameYearDay(compareTo time.Time) types.BeMatcher
```
SameYearDay checks if the .YearDay() component of the actual time matches the
.YearDay() component of the specified time `compareTo`. It only verifies the day
[1..365], disregarding the year.
