package datetime

import (
	"klog/datetime"
)

func Date_(year int, month int, day int) datetime.Date {
	date, err := datetime.NewDate(year, month, day)
	if err != nil {
		panic("Operation failed!")
	}
	return date
}

func Time_(hour int, minute int) datetime.Time {
	time, err := datetime.NewTime(hour, minute)
	if err != nil {
		panic("Operation failed!")
	}
	return time
}

func Range_(start datetime.Time, end datetime.Time) datetime.TimeRange {
	timeRange, err := datetime.NewTimeRange(start, end)
	if err != nil {
		panic("Operation failed!")
	}
	return timeRange
}

func OverlappingRange_(start datetime.Time, isStartYesterday bool, end datetime.Time, isEndTomorrow bool) datetime.TimeRange {
	timeRange, err := datetime.NewOverlappingTimeRange(start, isStartYesterday, end, isEndTomorrow)
	if err != nil {
		panic("Operation failed!")
	}
	return timeRange
}
