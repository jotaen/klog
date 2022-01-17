package service

import (
	. "github.com/jotaen/klog/src"
	gotime "time"
)

// DateTime represents a point in time with a normalized time value.
type DateTime struct {
	Date Date
	Time Time
}

func NewDateTime(d Date, t Time) DateTime {
	normalizedTime, _ := NewTime(t.Hour(), t.Minute())
	dayOffset := func() int {
		if t.IsTomorrow() {
			return 1
		} else if t.IsYesterday() {
			return -1
		}
		return 0
	}()
	return DateTime{
		Date: d.PlusDays(dayOffset),
		Time: normalizedTime,
	}
}

func NewDateTimeFromGo(reference gotime.Time) DateTime {
	date := NewDateFromGo(reference)
	time := NewTimeFromGo(reference)
	return NewDateTime(date, time)
}

func (dt DateTime) IsEqual(compare DateTime) bool {
	return dt.Date.IsEqualTo(compare.Date) && dt.Time.IsEqualTo(compare.Time)
}

func (dt DateTime) IsAfterOrEqual(compare DateTime) bool {
	if dt.Date.IsEqualTo(compare.Date) {
		return dt.Time.IsAfterOrEqual(compare.Time)
	}
	return dt.Date.IsAfterOrEqual(compare.Date)
}
