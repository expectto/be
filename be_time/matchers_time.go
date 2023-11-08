package be_time

import (
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
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

// todo:
// SameTimezone()
// SameDay, SameWeek, SameMonth, SameYear
