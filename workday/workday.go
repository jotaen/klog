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
	AddTime(datetime.Duration) error
	Ranges() [][]datetime.Time // tuple of [start, end]
	AddRange(datetime.Time, datetime.Time) error
	TotalTime() datetime.Duration
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

func (e *workday) AddTime(time datetime.Duration) error {
	if time < 0 {
		return errors.New("NEGATIVE_TIME")
	}
	e.times = append(e.times, time)
	return nil
}

func (e *workday) Ranges() [][]datetime.Time {
	ts := [][]datetime.Time{}
	for _, r := range e.ranges {
		ts = append(ts, []datetime.Time{r[0], r[1]})
	}
	return ts
}

func (e *workday) AddRange(start datetime.Time, end datetime.Time) error {
	e.ranges = append(e.ranges, []datetime.Time{start, end})
	return nil
}

func (e *workday) TotalTime() datetime.Duration {
	total := datetime.Duration(0)
	for _, t := range e.times {
		total += t
	}
	for _, rs := range e.ranges {
		start := rs[0].Minute() + 60*rs[0].Hour()
		end := rs[1].Minute() + 60*rs[1].Hour()
		total += datetime.Duration(end - start)
	}
	return total
}
