package workday

import (
	"klog/datetime"
)

type WorkDay interface {
	Date() datetime.Date
	Summary() string
	SetSummary(string)
	Times() []datetime.Duration
	AddDuration(datetime.Duration)
	Ranges() []datetime.TimeRange
	AddRange(datetime.TimeRange)
	OpenRangeStart() datetime.Time
	SetOpenRangeStart(datetime.Time)
	TotalWorkTime() datetime.Duration
}

func Create(date datetime.Date) WorkDay {
	return &workday{
		date: date,
	}
}

type workday struct {
	date           datetime.Date
	summary        string
	times          []datetime.Duration
	ranges         []datetime.TimeRange
	openRangeStart datetime.Time
}

func (e *workday) Date() datetime.Date {
	return e.date
}

func (e *workday) Summary() string {
	return e.summary
}

func (e *workday) SetSummary(summary string) {
	e.summary = summary
}

func (e *workday) Times() []datetime.Duration {
	return e.times
}

func (e *workday) AddDuration(time datetime.Duration) {
	e.times = append(e.times, time)
}

func (e *workday) Ranges() []datetime.TimeRange {
	return e.ranges
}

func (e *workday) AddRange(r datetime.TimeRange) {
	e.ranges = append(e.ranges, r)
}

func (e *workday) OpenRangeStart() datetime.Time {
	return e.openRangeStart
}

func (e *workday) SetOpenRangeStart(start datetime.Time) {
	e.openRangeStart = start
}

func (e *workday) TotalWorkTime() datetime.Duration {
	total := datetime.Duration(0)
	for _, t := range e.times {
		total += t
	}
	for _, r := range e.ranges {
		total += r.Duration()
	}
	return total
}
