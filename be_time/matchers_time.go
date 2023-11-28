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

// TODO: not sure it will work, tests are required

func LaterThan(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally(">", compareTo))
}

func LaterThanEqual(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally(">=", compareTo))
}

func EarlierThan(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally("<", compareTo))
}

func EarlierThanEqual(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally("<=", compareTo))
}

func Approx(arg time.Time, threshold time.Duration) types.BeMatcher {
	return Psi(gomega.BeTemporally("~", arg, threshold))
}

func SameNano(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally("~", compareTo, time.Nanosecond))
}

func SameSecond(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally("~", compareTo, time.Second))
}

func SameMinute(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally("~", compareTo, time.Minute))
}

func SameHour(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally("~", compareTo, time.Hour))
}

// SameTimezone checks if actual time is the same timezone as given time
func SameTimezone(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		return compareTo.Location().String() == actualTime.Location().String(), nil
	}))
}

// SameOffset checks if actual time is the same timezone offset as given time
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

// SameDay checks if given and actual times are the same day (timezone is respected)
func SameDay(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		return compareTo.Day() == actualTime.Day(), nil
	}))
}

// SameWeekday checks if given and actual times are the same weekday (timezone is respected)
func SameWeekday(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		return compareTo.Weekday() == actualTime.Weekday(), nil
	}))
}

// SameWeek checks if given and actual times are the same week (timezone is respected)
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

// SameMonth checks if given and actual times are the same month (timezone is respected)
func SameMonth(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		return compareTo.Month() == actualTime.Month(), nil
	}))
}

// SameYear checks if given and actual times are the same year (timezone is respected)
func SameYear(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		if !cast.IsTime(actual) {
			return false, fmt.Errorf("invalid time type")
		}
		actualTime := cast.AsTime(actual)

		return compareTo.Year() == actualTime.Year(), nil
	}))
}
