package workday

import (
	"errors"
	"klog/datetime"
)

type WorkDay interface {
	Date() datetime.Date
	SetDate(datetime.Date)
	Summary() string
	SetSummary(string) error
	Times() []datetime.Duration
	AddDuration(datetime.Duration) error
	Ranges() [][]datetime.Time // tuple of start and end time (end can be `nil`)
	AddRange(datetime.Time, datetime.Time) error
	AddOpenRange(datetime.Time) error
	TotalWorkTime() datetime.Duration
}

func Create(date datetime.Date) WorkDay {
	return &workday{
		date: date,
	}
}

type workday struct {
	date    datetime.Date
	summary string
	times   []datetime.Duration
	ranges  [][]datetime.Time
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
	if time < 0 {
		return errors.New("NEGATIVE_DURATION")
	}
	e.times = append(e.times, time)
	return nil
}

func (e *workday) Ranges() [][]datetime.Time {
	return e.ranges
}

func (e *workday) AddRange(start datetime.Time, end datetime.Time) error {
	e.ranges = append(e.ranges, []datetime.Time{start, end})
	return nil
}

func (e *workday) AddOpenRange(start datetime.Time) error {
	e.ranges = append(e.ranges, []datetime.Time{start, nil})
	return nil
}

func (e *workday) TotalWorkTime() datetime.Duration {
	total := datetime.Duration(0)
	for _, t := range e.times {
		total += t
	}
	for _, rs := range e.ranges {
		if rs[1] == nil {
			continue
		}
		start := rs[0].Minute() + 60*rs[0].Hour()
		end := rs[1].Minute() + 60*rs[1].Hour()
		total += datetime.Duration(end - start)
	}
	return total
}
