package be_time

import (
	"fmt"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gcustom"
	"time"
)

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

func SameSecond(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally("~", compareTo, time.Second))
}

func SameMinute(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally("~", compareTo, time.Minute))
}

func SameHour(compareTo time.Time) types.BeMatcher {
	return Psi(gomega.BeTemporally("~", compareTo, time.Hour))
}
func SameTimezone(compareTo time.Time) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		// todo: move into cast.As() helpers
		var actualTime time.Time
		var ok bool

		if actualTime, ok = actual.(time.Time); !ok {
			ptrActualTime, ok := actual.(*time.Time)
			if !ok {
				return false, fmt.Errorf("actual must be a valid time.Time")
			}
			actualTime = *ptrActualTime
		}

		return compareTo.Location().String() == actualTime.Location().String(), nil
	}))
}

// todo: SameOffset()
//       SameDay(), SameWeek(), SameMonth() SameYear()
