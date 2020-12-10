package datetime

import (
	"errors"
)

type TimeRange interface {
	Start() Time
	End() Time
	Duration() Duration
}

type timeRange struct {
	start Time
	end   Time
}

func CreateTimeRange(start Time, end Time) (TimeRange, error) {
	startMinutes := start.Hour()*60 + start.Minute()
	endMinutes := end.Hour()*60 + end.Minute()
	if endMinutes < startMinutes {
		return nil, errors.New("ILLEGAL_RANGE")
	}
	return timeRange{
		start: start,
		end:   end,
	}, nil
}

func (tr timeRange) Start() Time {
	return tr.start
}

func (tr timeRange) End() Time {
	return tr.end
}

func (tr timeRange) Duration() Duration {
	start := tr.Start().MinutesSinceMidnight()
	end := tr.End().MinutesSinceMidnight()
	return Duration(end - start)
}
