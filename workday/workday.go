package workday

import (
	"klog/datetime"
)

type WorkDay interface {
	Date() datetime.Date
	SetDate(datetime.Date)
	Summary() string
	SetSummary(string) error
	Times() []datetime.Duration
	AddDuration(datetime.Duration) error
	Ranges() []datetime.TimeRange
	AddRange(datetime.TimeRange) error
	OpenRange() datetime.TimeRange
	TotalWorkTime() datetime.Duration
}

func Create(date datetime.Date) WorkDay {
	return &workday{
		date: date,
	}
}

type workday struct {
	date      datetime.Date
	summary   string
	times     []datetime.Duration
	ranges    []datetime.TimeRange
	openRange datetime.TimeRange
}

func (e *workday) Date() datetime.Date {
	return e.date
}

func (e *workday) SetDate(date datetime.Date) {
	e.date = date
}

func (e *workday) Summary() string {
	return e.summary
}

func (e *workday) SetSummary(summary string) error {
	e.summary = summary
	return nil
}

func (e *workday) Times() []datetime.Duration {
	return e.times
}

func (e *workday) AddDuration(time datetime.Duration) error {
	e.times = append(e.times, time)
	return nil
}

func (e *workday) Ranges() []datetime.TimeRange {
	return e.ranges
}

func (e *workday) AddRange(r datetime.TimeRange) error {
	e.ranges = append(e.ranges, r)
	return nil
}

func (e *workday) OpenRange() datetime.TimeRange {
	var res datetime.TimeRange
	for _, r := range e.ranges {
		if r.IsOpen() {
			return r
		}
	}
	return res
}

func (e *workday) TotalWorkTime() datetime.Duration {
	total := datetime.Duration(0)
	for _, t := range e.times {
		total += t
	}
	for _, r := range e.ranges {
		if r.IsOpen() {
			continue
		}
		total += datetime.Duration(r.End().MinutesSinceMidnight() - r.Start().MinutesSinceMidnight())
	}
	return total
}
