package workday

import (
	"cloud.google.com/go/civil"
	"klog/datetime"
	"time"
)

type WorkDay interface {
	Date() datetime.Date
	SetDate(datetime.Date) error
	Summary() string
	SetSummary(string) error
	Times() []datetime.Minutes
	AddTime(datetime.Minutes) error
	Ranges() [][]datetime.Time // tuple of [start, end]
	AddRange(datetime.Time, datetime.Time) error
	TotalTime() datetime.Minutes
}

func Create(date datetime.Date) (WorkDay, error) {
	e := workday{}
	err := e.SetDate(date)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

type workday struct {
	date    civil.Date
	summary string
	times   []datetime.Minutes
	ranges  [][]civil.Time
}

func (e *workday) Date() datetime.Date {
	return datetime.Date{
		Year:  e.date.Year,
		Month: int(e.date.Month),
		Day:   e.date.Day,
	}
}

func (e *workday) SetDate(date datetime.Date) error {
	d := civil.Date{
		Year:  date.Year,
		Month: time.Month(date.Month),
		Day:   date.Day,
	}
	if !d.IsValid() {
		return &WorkDayError{Code: INVALID_DATE}
	}
	e.date = d
	return nil
}

func (e *workday) Summary() string {
	return e.summary
}

func (e *workday) SetSummary(summary string) error {
	e.summary = summary
	return nil
}

func (e *workday) Times() []datetime.Minutes {
	return e.times
}

func (e *workday) AddTime(time datetime.Minutes) error {
	if time < 0 {
		return &WorkDayError{Code: NEGATIVE_TIME}
	}
	e.times = append(e.times, time)
	return nil
}

func (e *workday) Ranges() [][]datetime.Time {
	ts := [][]datetime.Time{}
	for _, r := range e.ranges {
		ts = append(ts, []datetime.Time{
			datetime.Time{Hour: r[0].Hour, Minute: r[0].Minute},
			datetime.Time{Hour: r[1].Hour, Minute: r[1].Minute},
		})
	}
	return ts
}

func (e *workday) AddRange(start datetime.Time, end datetime.Time) error {
	e.ranges = append(e.ranges, []civil.Time{
		civil.Time{Hour: start.Hour, Minute: start.Minute, Second: 0, Nanosecond: 0},
		civil.Time{Hour: end.Hour, Minute: end.Minute, Second: 0, Nanosecond: 0},
	})
	return nil
}

func (e *workday) TotalTime() datetime.Minutes {
	total := datetime.Minutes(0)
	for _, t := range e.times {
		total += t
	}
	for _, rs := range e.ranges {
		start := rs[0].Minute + 60*rs[0].Hour
		end := rs[1].Minute + 60*rs[1].Hour
		total += datetime.Minutes(end - start)
	}
	return total
}
