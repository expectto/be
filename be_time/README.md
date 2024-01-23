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
SameDay succeeds if the day component of the actual time is equal to the day
component of the specified time `compareTo`.

#### func  SameHour

```go
func SameHour(compareTo time.Time) types.BeMatcher
```
SameHour succeeds if actual time is approximately equal to the specified time
`compareTo` with the precision of one hour.

#### func  SameMinute

```go
func SameMinute(compareTo time.Time) types.BeMatcher
```
SameMinute succeeds if actual time is approximately equal to the specified time
`compareTo` with the precision of one minute.

#### func  SameMonth

```go
func SameMonth(compareTo time.Time) types.BeMatcher
```
SameMonth succeeds if the month component of the actual time is equal to the
month component of the specified time `compareTo`.

#### func  SameNano

```go
func SameNano(compareTo time.Time) types.BeMatcher
```
SameNano succeeds if actual time is approximately equal to the specified time
`compareTo` with the precision of one nanosecond.

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
SameSecond succeeds if actual time is approximately equal to the specified time
`compareTo` with the precision of one second.

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
SameWeek succeeds if the ISO week and year components of the actual time are
equal to the ISO week and year components of the specified time `compareTo`.

#### func  SameWeekday

```go
func SameWeekday(compareTo time.Time) types.BeMatcher
```
SameWeekday succeeds if the weekday component of the actual time is equal to the
weekday component of the specified time `compareTo`.

#### func  SameYear

```go
func SameYear(compareTo time.Time) types.BeMatcher
```
SameYear succeeds if the year component of the actual time is equal to the year
component of the specified time `compareTo`.
