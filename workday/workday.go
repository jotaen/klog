package workday

import (
	"errors"
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

	OpenRange() datetime.Time
	StartOpenRange(datetime.Time)
	EndOpenRange(datetime.Time) error

	TotalWorkTime() datetime.Duration
	TotalWorkTimeWithOpenRange(time datetime.Time) (datetime.Duration, error)
}

func NewWorkDay(date datetime.Date) WorkDay {
	return &workday{
		date: date,
	}
}

type workday struct {
	date           datetime.Date
	summary        string
	times          []datetime.Duration
	ranges         []datetime.TimeRange
	openRangeBegin datetime.Time
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

func (e *workday) OpenRange() datetime.Time {
	return e.openRangeBegin
}

func (e *workday) StartOpenRange(start datetime.Time) {
	e.openRangeBegin = start
}

func (e *workday) EndOpenRange(end datetime.Time) error {
	r, _ := datetime.NewTimeRange(e.openRangeBegin, end)
	e.openRangeBegin = nil
	e.AddRange(r)
	return nil // TODO
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

func (e *workday) TotalWorkTimeWithOpenRange(end datetime.Time) (datetime.Duration, error) {
	if e.openRangeBegin == nil {
		return 0, errors.New("NO_OPEN_RANGE")
	}
	total := e.TotalWorkTime()
	r, _ := datetime.NewTimeRange(e.openRangeBegin, end)
	total += r.Duration()
	return total, nil
}
