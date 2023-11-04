package be

import (
	"github.com/expectto/be/internal/psi"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
	"time"
)

func BeLaterThan(compareTo time.Time) types.BeMatcher {
	return psi.Psi(gomega.BeTemporally(">", compareTo))
}

func BeLaterThanEqual(compareTo time.Time) types.BeMatcher {
	return psi.Psi(gomega.BeTemporally(">=", compareTo))
}

func BeEarlierThan(compareTo time.Time) types.BeMatcher {
	return psi.Psi(gomega.BeTemporally("<", compareTo))
}

func BeEarlierThanEqual(compareTo time.Time) types.BeMatcher {
	return psi.Psi(gomega.BeTemporally("<=", compareTo))
}

func BeApproxTime(arg time.Time, threshold time.Duration) types.BeMatcher {
	return psi.Psi(gomega.BeTemporally("~", arg, threshold))
}

// todo:
// SameSecond, SameMinute, SameHour, tz: SameDay, SameWeek, SameMonth, SameYear
// BeSameTimezone()
