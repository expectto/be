package be_time_test

import (
	"time"

	"github.com/expectto/be/be_time"
	"github.com/expectto/be/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TimeFormat used for tests is selected to be 1) allmighty 2) more readable
// It's a modification of time.RFC3339
const TimeFormat = "2006-01-02T15:04:05.999999999Z"

type TestCase struct {
	GetMatcher func(time.Time) types.BeMatcher
	Actual     string
	ActualTz   string // because of Actual will always be a string with `Z` suffix, we'll need to convert it back to proper tz
	Success    bool
}

func (tc *TestCase) SetTz(tz string) *TestCase {
	tc.ActualTz = tz
	return tc
}

// TC is a short way of creating positive TestCase
func TC_OK(m func(time.Time) types.BeMatcher, actual string) *TestCase {
	return &TestCase{Success: true, Actual: actual, GetMatcher: m}
}

// TC_NOT_OK is a short way of creating negative TestCase
func TC_NOT_OK(m func(time.Time) types.BeMatcher, actual string) *TestCase {
	return &TestCase{Success: false, Actual: actual, GetMatcher: m}
}

// compareTo is the reference "compare-to" date for most of our test cases
var compareTo = time.Date(2024, time.February, 02, 15, 30, 0, 0, time.UTC)
var (
	locationEST, _ = time.LoadLocation("EST")
)

var _ = DescribeTable("Be Time: Matching", func(tc *TestCase) {
	var loc = time.UTC
	if tc.ActualTz != "" {
		var err error
		loc, err = time.LoadLocation(tc.ActualTz)
		Expect(err).Should(Succeed(), "Invalid Test TZ given")
	}
	actual, err := time.ParseInLocation(TimeFormat, tc.Actual, loc)
	Expect(err).To(Succeed(), "Invalid Test given")

	success, err := tc.GetMatcher(compareTo).Match(actual)
	Expect(err).ToNot(HaveOccurred())

	Expect(success).To(Equal(tc.Success))
},
	// EarlierThan
	Entry("Earlier than: actual is earlier", TC_OK(be_time.EarlierThan, "2024-01-31T00:00:00.0Z")),
	Entry("Earlier than: actual is the same", TC_NOT_OK(be_time.EarlierThan, compareTo.Format(TimeFormat))),
	Entry("Earlier than: actual is later", TC_NOT_OK(be_time.EarlierThan, "2024-02-25T00:00:00.0Z")),

	// EarlierThanEqual
	Entry("Earlier than equal: actual is earlier", TC_OK(be_time.EarlierThanEqual, "2024-01-31T00:00:00.0Z")),
	Entry("Earlier than equal: actual is the same", TC_OK(be_time.EarlierThanEqual, compareTo.Format(TimeFormat))),
	Entry("Earlier than equal: actual is later", TC_NOT_OK(be_time.EarlierThanEqual, "2024-02-25T00:00:00.0Z")),

	// LaterThan
	Entry("Later than: actual is later", TC_OK(be_time.LaterThan, "2024-02-25T00:00:00.0Z")),
	Entry("Later than: actual is the same", TC_NOT_OK(be_time.LaterThan, compareTo.Format(TimeFormat))),
	Entry("Later than: actual is earlier", TC_NOT_OK(be_time.LaterThan, "2024-01-31T00:00:00.0Z")),

	// LaterThanEqual
	Entry("Later than equal: actual is later", TC_OK(be_time.LaterThanEqual, "2024-02-25T00:00:00.0Z")),
	Entry("Later than equal: actual is the same", TC_OK(be_time.LaterThanEqual, compareTo.Format(TimeFormat))),
	Entry("Later than equal: actual is earlier", TC_NOT_OK(be_time.LaterThanEqual, "2024-01-31T00:00:00.0Z")),

	// Eq
	Entry("Eq: actual is the same", TC_OK(be_time.Eq, compareTo.Format(TimeFormat))),
	Entry("Eq: actual differs by +5ns", TC_NOT_OK(be_time.Eq, compareTo.Add(5*time.Nanosecond).Format(TimeFormat))),
	Entry("Eq: actual differs by +1ns", TC_NOT_OK(be_time.Eq, compareTo.Add(+time.Nanosecond).Format(TimeFormat))),
	Entry("Eq: actual differs by -5ns", TC_NOT_OK(be_time.Eq, compareTo.Add(-5*time.Nanosecond).Format(TimeFormat))),
	Entry("Eq: actual differs by -1ns", TC_NOT_OK(be_time.Eq, compareTo.Add(-time.Nanosecond).Format(TimeFormat))),

	// SameExactMilli: here we KNOW that compareTO is 15:30:00.000
	Entry("Same Exact Milli: actual is the same", TC_OK(be_time.SameExactMilli, compareTo.Format(TimeFormat))),
	Entry("Same Exact Milli: actual differs by 1ns (within 1ms)", TC_OK(be_time.SameExactMilli, compareTo.Add(1*time.Nanosecond).Format(TimeFormat))),
	Entry("Same Exact Milli: actual differs by 1ns (outside 1ms)", TC_NOT_OK(be_time.SameExactMilli, compareTo.Add(-1*time.Nanosecond).Format(TimeFormat))),
	Entry("Same Exact Milli: actual differs by 999µs (within 1ms)", TC_OK(be_time.SameExactMilli, compareTo.Add(999*time.Microsecond).Format(TimeFormat))),
	Entry("Same Exact Milli: actual differs by 999µs (outside 1ms)", TC_NOT_OK(be_time.SameExactMilli, compareTo.Add(-999*time.Microsecond).Format(TimeFormat))),
	Entry("Same Exact Milli: actual differs by 5ms", TC_NOT_OK(be_time.SameExactMilli, compareTo.Add(5*time.Millisecond).Format(TimeFormat))),

	// SameExactSecond: here we KNOW that compareTO is 15:30:00.000
	Entry("Same Exact Second: actual is the same", TC_OK(be_time.SameExactSecond, compareTo.Format(TimeFormat))),
	Entry("Same Exact Second: actual differs by 1ns (within 1s)", TC_OK(be_time.SameExactSecond, compareTo.Add(1*time.Nanosecond).Format(TimeFormat))),
	Entry("Same Exact Second: actual differs by 1ns (outside 1s)", TC_NOT_OK(be_time.SameExactSecond, compareTo.Add(-1*time.Nanosecond).Format(TimeFormat))),
	Entry("Same Exact Second: actual differs by 999ms (within 1s)", TC_OK(be_time.SameExactSecond, compareTo.Add(999*time.Millisecond).Format(TimeFormat))),
	Entry("Same Exact Second: actual differs by 999ms (outside 1s)", TC_NOT_OK(be_time.SameExactSecond, compareTo.Add(-999*time.Millisecond).Format(TimeFormat))),
	Entry("Same Exact Second: actual differs by 5s", TC_NOT_OK(be_time.SameExactSecond, compareTo.Add(5*time.Second).Format(TimeFormat))),

	// SameExactMinute: here we KNOW that compareTO is 15:30:00.000
	Entry("Same Exact Minute: actual is the same", TC_OK(be_time.SameExactMinute, compareTo.Format(TimeFormat))),
	Entry("Same Exact Minute: actual differs by 1ns (within 1m)", TC_OK(be_time.SameExactMinute, compareTo.Add(1*time.Nanosecond).Format(TimeFormat))),
	Entry("Same Exact Minute: actual differs by 1ns (outside 1m)", TC_NOT_OK(be_time.SameExactMinute, compareTo.Add(-1*time.Nanosecond).Format(TimeFormat))),
	Entry("Same Exact Minute: actual differs by 59s (within 1m)", TC_OK(be_time.SameExactMinute, compareTo.Add(59*time.Second).Format(TimeFormat))),
	Entry("Same Exact Minute: actual differs by 59s (outside 1m)", TC_NOT_OK(be_time.SameExactMinute, compareTo.Add(-59*time.Second).Format(TimeFormat))),
	Entry("Same Exact Minute: actual differs by 61s", TC_NOT_OK(be_time.SameExactMinute, compareTo.Add(61*time.Second).Format(TimeFormat))),

	// SameExactHour: here we KNOW that compareTO is 15:30:00.000
	Entry("Same Exact Hour: actual is the same", TC_OK(be_time.SameExactHour, compareTo.Format(TimeFormat))),
	Entry("Same Exact Hour: actual differs by 1ns (within 1h)", TC_OK(be_time.SameExactHour, compareTo.Add(1*time.Nanosecond).Format(TimeFormat))),
	Entry("Same Exact Hour: actual differs by 29m (within 1h)", TC_OK(be_time.SameExactHour, compareTo.Add(29*time.Minute).Format(TimeFormat))),
	Entry("Same Exact Hour: actual differs by 30m (outside 1h)", TC_NOT_OK(be_time.SameExactHour, compareTo.Add(30*time.Minute).Format(TimeFormat))),
	Entry("Same Exact Hour: actual differs by 61m", TC_NOT_OK(be_time.SameExactHour, compareTo.Add(61*time.Minute).Format(TimeFormat))),

	// SameTimezone: here we KNOW that compareTO is in UTC timezone
	Entry("Same Timezone: actual is in UTC", TC_OK(be_time.SameTimezone, compareTo.Format(TimeFormat))),
	Entry("Same Timezone: actual is in EST", TC_NOT_OK(be_time.SameTimezone, compareTo.In(locationEST).Format(TimeFormat)).SetTz("EST")),
	Entry("Same Timezone: actual time is different but in the same timezone", TC_OK(be_time.SameTimezone, compareTo.Add(5*time.Minute).Format(TimeFormat))),
	Entry("Same Timezone: actual time is different and in different timezone", TC_NOT_OK(be_time.SameTimezone, compareTo.In(locationEST).Format(TimeFormat)).SetTz("EST")),

	// TODO: other matchers
)

var _ = Context("Time atomic matchers", func() {
	It("should pass on atomic matchers", func() {
		Expect(compareTo).To(be_time.Year(2024))
		Expect(compareTo).To(be_time.Month(time.February))
		Expect(compareTo).To(be_time.YearDay(33))
		Expect(compareTo).To(be_time.Weekday(time.Friday))
		Expect(compareTo).To(be_time.Day(2))
	})
})
