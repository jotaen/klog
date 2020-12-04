package workday

import (
	"cloud.google.com/go/civil"
	"time"
)

type Minutes int64

type Time struct {
	Hour   int
	Minute int
}

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

type WorkDay interface {
	Date() Date
	SetDate(Date) error
	Summary() string
	SetSummary(string) error
	Times() []Minutes
	AddTime(Minutes) error
	Ranges() [][]Time // tuple of [start, end]
	AddRange(Time, Time) error
	TotalTime() Minutes
}

func Create(date Date) (WorkDay, error) {
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
	times   []Minutes
	ranges  [][]civil.Time
}

func (e *workday) Date() Date {
	return Date{
		Year:  e.date.Year,
		Month: e.date.Month,
		Day:   e.date.Day,
	}
}

func (e *workday) SetDate(date Date) error {
	d := civil.Date{
		Year:  date.Year,
		Month: date.Month,
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

func (e *workday) Times() []Minutes {
	return e.times
}

func (e *workday) AddTime(time Minutes) error {
	if time < 0 {
		return &WorkDayError{Code: NEGATIVE_TIME}
	}
	e.times = append(e.times, time)
	return nil
}

func (e *workday) Ranges() [][]Time {
	ts := [][]Time{}
	for _, r := range e.ranges {
		ts = append(ts, []Time{
			Time{Hour: r[0].Hour, Minute: r[0].Minute},
			Time{Hour: r[1].Hour, Minute: r[1].Minute},
		})
	}
	return ts
}

func (e *workday) AddRange(start Time, end Time) error {
	e.ranges = append(e.ranges, []civil.Time{
		civil.Time{Hour: start.Hour, Minute: start.Minute, Second: 0, Nanosecond: 0},
		civil.Time{Hour: end.Hour, Minute: end.Minute, Second: 0, Nanosecond: 0},
	})
	return nil
}

func (e *workday) TotalTime() Minutes {
	total := Minutes(0)
	for _, t := range e.times {
		total += t
	}
	for _, rs := range e.ranges {
		start := rs[0].Minute + 60*rs[0].Hour
		end := rs[1].Minute + 60*rs[1].Hour
		total += Minutes(end - start)
	}
	return total
}
