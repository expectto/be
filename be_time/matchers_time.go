// Package be_time provides Be matchers on time.Time
package be_time

import (
	"fmt"
	"github.com/expectto/be/internal/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gcustom"
	"time"
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

// Approx succeeds if actual time is approximately equal to the specified time `compareTo`
// within the given time duration threshold.
func Approx(compareTo time.Time, threshold time.Duration) types.BeMatcher {
	return Psi(gomega.BeTemporally("~", compareTo, threshold))
}

// SameNano succeeds if actual time is approximately equal to the specified time `compareTo`
// with the precision of one nanosecond.
func SameNano(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally("~", compareTo, time.Nanosecond))
}

// SameSecond succeeds if actual time is approximately equal to the specified time `compareTo`
// with the precision of one second.
func SameSecond(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally("~", compareTo, time.Second))
}

// SameMinute succeeds if actual time is approximately equal to the specified time `compareTo`
// with the precision of one minute.
func SameMinute(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally("~", compareTo, time.Minute))
}

// SameHour succeeds if actual time is approximately equal to the specified time `compareTo`
// with the precision of one hour.
func SameHour(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally("~", compareTo, time.Hour))
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

// SameDay succeeds if the day component of the actual time is equal to the day component
// of the specified time `compareTo`.
func SameDay(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		return compareTo.Day() == actualTime.Day(), nil
	}))
}

// SameWeekday succeeds if the weekday component of the actual time is equal to the weekday component
// of the specified time `compareTo`.
func SameWeekday(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		return compareTo.Weekday() == actualTime.Weekday(), nil
	}))
}

// SameWeek succeeds if the ISO week and year components of the actual time
// are equal to the ISO week and year components of the specified time `compareTo`.
func SameWeek(compareTo time.Time) types.BeMatcher {
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

// SameMonth succeeds if the month component of the actual time is equal to the month component
// of the specified time `compareTo`.
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

// TODO: add aliases, but ensure confusion is not introduced
// 		as we can't do Lt for LaterThan and Et for EarlierThan
// 		Aliases should be consistent with be_math. So Lt is LessThan, meaning EarlierThan
//              and Gt is GreaterThan, meaning LaterThan
