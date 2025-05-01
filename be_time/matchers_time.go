// Package be_time provides Be matchers on time.Time
package be_time

import (
	"fmt"
	"time"

	"github.com/amberpixels/abu/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gcustom"
)

// LaterThan succeeds if actual time is later than the specified time `compareTo`.
func LaterThan(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally(">", compareTo))
}

// LaterThanEqual succeeds if actual time is later than or equal to the specified time `compareTo`.
func LaterThanEqual(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally(">=", compareTo))
}

// EarlierThan succeeds if actual time is earlier than the specified time `compareTo`.
func EarlierThan(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally("<", compareTo))
}

// EarlierThanEqual succeeds if actual time is earlier than or equal to the specified time `compareTo`.
func EarlierThanEqual(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally("<=", compareTo))
}

// Eq succeeds if actual time is equal to the specified time `compareTo`
// with the precision of one nanosecond.
func Eq(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally("==", compareTo))
}

// Approx succeeds if actual time is approximately equal to the specified time `compareTo`
// within the given time duration threshold.
func Approx(compareTo time.Time, threshold time.Duration) types.BeMatcher {
	return Psi(gomega.BeTemporally("~", compareTo, threshold))
}

//
// Atomic matching of time parts
//

func atomicTimePartMatcher[T comparable](actualGetter func(t time.Time) T, compareTo T, customMessageArg ...string) types.BeMatcher {
	message := fmt.Sprintf("%v", compareTo)
	if len(customMessageArg) > 0 {
		message = customMessageArg[0]
	}

	return Psi(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		return actualGetter(cast.AsTime(actual)) == compareTo, nil

		// TODO: would be great if Expect part of the message contains reference value as well
		//    e.g. Expected <time.Time>: .. (TUESDAY) to be Friday
	}, "be "+message)
}

func Year(v int) types.BeMatcher {
	return atomicTimePartMatcher(func(t time.Time) int { return t.Year() }, v)
}
func Month(monthCompareTo time.Month) types.BeMatcher {
	return atomicTimePartMatcher(func(t time.Time) time.Month { return t.Month() }, monthCompareTo)
}
func Day(v int) types.BeMatcher {
	return atomicTimePartMatcher(func(t time.Time) int { return t.Day() }, v, fmt.Sprintf("%d day of month", v))
}
func YearDay(v int) types.BeMatcher {
	return atomicTimePartMatcher(func(t time.Time) int { return t.YearDay() }, v, fmt.Sprintf("%d day of the year", v))
}
func Weekday(v time.Weekday) types.BeMatcher {
	return atomicTimePartMatcher(func(t time.Time) time.Weekday { return t.Weekday() }, v)
}
func Unix(v int64) types.BeMatcher {
	return atomicTimePartMatcher(func(t time.Time) int64 { return t.Unix() }, v, fmt.Sprintf("equal to %d Unix timestamp", v))
}

// TODO: more atomic matchers

//
// --- Same Exact * ---
//

// TODO: SameExact* matchers. Naming should note that it's 2 times comparison
//   so probably it should be Same...With())

// sameExactDuration is an internal matcher that succeeds if
// actual time falls within the same X duration as the specified time `compareTo`
// It's considered to avoid duplication in implementation of public matchers `SameExact*`
func sameExactDuration(compareTo time.Time, duration time.Duration) types.BeMatcher {
	return Psi(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		// Truncate both times to the given duration
		truncatedCompareTo := compareTo.Truncate(duration)
		truncatedActual := actualTime.Truncate(duration)

		// Compare truncated times
		return truncatedCompareTo.Equal(truncatedActual), nil
		// TODO: better message
	}, fmt.Sprintf("be same as %s", duration))
}

// SameExactMilli succeeds if the actual time falls within the same millisecond as the specified time `compareTo`.
func SameExactMilli(compareTo time.Time) types.BeMatcher {
	return sameExactDuration(compareTo, time.Millisecond)
}

// SameExactSecond succeeds if the actual time falls within the same second as the specified time `compareTo`.
func SameExactSecond(compareTo time.Time) types.BeMatcher {
	return sameExactDuration(compareTo, time.Second)
}

// SameExactMinute succeeds if the actual time falls within the same minute as the specified time `compareTo`.
func SameExactMinute(compareTo time.Time) types.BeMatcher {
	return sameExactDuration(compareTo, time.Minute)
}

// SameExactHour succeeds if the actual time falls within the same hour as the specified time `compareTo`.
func SameExactHour(compareTo time.Time) types.BeMatcher {
	return sameExactDuration(compareTo, time.Hour)
}

// SameExactDay succeeds if the actual time falls within the same day as the specified time `compareTo`.
func SameExactDay(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		// Different is less than 24 * Hours and they are within the same day
		return compareTo.Sub(actualTime).Abs() < 24*time.Hour && compareTo.Day() == actualTime.Day(), nil
	}))
}

// SameExactWeekday succeeds if the weekday component of the actual time is equal to the weekday component
// of the specified time `compareTo`.
func SameExactWeekday(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		return compareTo.Weekday() == actualTime.Weekday(), nil
	}))
}

// SameExactWeek succeeds if the actual time falls within the same ISO week as the specified time `compareTo`.
func SameExactWeek(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		compareToYear, compareToWeek := compareTo.ISOWeek()
		actualYear, actualWeek := actualTime.ISOWeek()
		return compareToYear == actualYear && compareToWeek == actualWeek, nil
	}))
}

// SameExactMonth succeeds if the actual time falls within the same month as the specified time `compareTo`.
func SameExactMonth(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		return compareTo.Year() == actualTime.Year() && compareTo.Month() == actualTime.Month(), nil
	}))
}

//
// --- Same * ---
//

// SameSecond checks if the .Second() component of the actual time matches
// the .Second() component of the specified time `compareTo`.
// It only verifies the second [0..59], disregarding other components such as minute, hour, day, etc.
func SameSecond(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		return compareTo.Second() == actualTime.Second(), nil
	}))
}

// SameMinute checks if the .Minute() component of the actual time matches
// the .Minute() component of the specified time `compareTo`.
// It only verifies the minute [0..59], disregarding other components such as second, hour, day, etc.
func SameMinute(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		return compareTo.Minute() == actualTime.Minute(), nil
	}))
}

// SameHour checks if the .Hour() component of the actual time matches
// the .Hour() component of the specified time `compareTo`.
// It only verifies the hour [0..59], disregarding other components such as second, minute, day, etc.
func SameHour(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		return compareTo.Hour() == actualTime.Hour(), nil
	}))
}

// SameDay checks if the .Day() component of the actual time matches
// the .Day() component of the specified time `compareTo`.
// It only verifies the day [1..30(31)], disregarding the month, year, etc.
func SameDay(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		return compareTo.Day() == actualTime.Day(), nil
	}))
}

// SameYearDay checks if the .YearDay() component of the actual time matches
// the .YearDay() component of the specified time `compareTo`.
// It only verifies the day [1..365], disregarding the year.
func SameYearDay(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		return compareTo.YearDay() == actualTime.YearDay(), nil
	}))
}

// SameWeek succeeds if the ISO week of the actual time
// are equal to the ISO week and year components of the specified time `compareTo`.
// It only verifies the week [1..53], disregarding of year.
// Note: use SameExactWeek to respect exact week of exact year
func SameWeek(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		_, compareToWeek := compareTo.ISOWeek()
		_, actualWeek := actualTime.ISOWeek()
		return compareToWeek == actualWeek, nil
	}))
}

// SameMonth checks if the .Month() component of the actual time matches
// the .Month() component of the specified time `compareTo`.
// It only verifies the month [1..12], disregarding the year.
func SameMonth(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		return compareTo.Month() == actualTime.Month(), nil
	}))
}

// SameYear succeeds if the year component of the actual time is equal to the year component
// of the specified time `compareTo`.
func SameYear(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		return compareTo.Year() == actualTime.Year(), nil
	}))
}

// SameTimezone checks if actual time is the same timezone as specified time `compareTo`
func SameTimezone(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		return compareTo.Location().String() == actualTime.Location().String(), nil
	}))
}

// SameOffset checks if actual time is the same timezone offset as specified time `compareTo`
// Note: times can have different timezone names, but same offset, e.g. America/New_York and Canada/Toronto
func SameOffset(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		_, offsetGiven := compareTo.Zone()
		_, offsetActual := actualTime.Zone()
		return offsetGiven == offsetActual, nil
	}))
}

// IsDST checks if actual time is DST
func IsDST(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)
		return actualTime.IsDST(), nil
	}))
}
