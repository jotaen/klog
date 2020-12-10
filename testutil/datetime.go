package testutil

import (
	"klog/datetime"
)

func Date_(year int, month int, day int) datetime.Date {
	date, err := datetime.CreateDate(year, month, day)
	if err != nil {
		panic("Operation failed!")
	}
	return date
}

func Time_(hour int, minute int) datetime.Time {
	time, err := datetime.CreateTime(hour, minute)
	if err != nil {
		panic("Operation failed!")
	}
	return time
}

func Range_(start datetime.Time, end datetime.Time) datetime.TimeRange {
	timeRange, err := datetime.CreateTimeRange(start, end)
	if err != nil {
		panic("Operation failed!")
	}
	return timeRange
}
