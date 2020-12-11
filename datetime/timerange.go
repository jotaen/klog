package datetime

import (
	"errors"
)

type TimeRange interface {
	Start() Time
	End() Time
	IsStartYesterday() bool
	IsEndTomorrow() bool
	Duration() Duration
}

type timeRange struct {
	start            Time
	end              Time
	isStartYesterday bool
	isEndTomorrow    bool
}

func CreateTimeRange(start Time, end Time) (TimeRange, error) {
	return CreateOverlappingTimeRange(start, false, end, false)
}

func CreateOverlappingTimeRange(start Time, isStartYesterday bool, end Time, isEndTomorrow bool) (TimeRange, error) {
	startMinutes := start.Hour()*60 + start.Minute()
	endMinutes := end.Hour()*60 + end.Minute()
	if !isStartYesterday && !isEndTomorrow && endMinutes < startMinutes {
		return nil, errors.New("ILLEGAL_RANGE")
	}
	return timeRange{
		start:            start,
		end:              end,
		isStartYesterday: isStartYesterday,
		isEndTomorrow:    isEndTomorrow,
	}, nil
}

func (tr timeRange) Start() Time {
	return tr.start
}

func (tr timeRange) End() Time {
	return tr.end
}

func (tr timeRange) IsStartYesterday() bool {
	return tr.isStartYesterday
}

func (tr timeRange) IsEndTomorrow() bool {
	return tr.isEndTomorrow
}

func (tr timeRange) Duration() Duration {
	start := tr.Start().MinutesSinceMidnight()
	end := tr.End().MinutesSinceMidnight()
	return Duration(end - start)
}
