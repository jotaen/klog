package entry

import (
	"cloud.google.com/go/civil"
	"time"
)

type Minutes int64

type Time struct {
	Hour int
	Minute int
}

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

type Entry interface {
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

func Create(date Date) (Entry, error) {
	e := entry{}
	err := e.SetDate(date)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

type entry struct {
	date    civil.Date
	summary string
	times   []Minutes
	ranges  [][]civil.Time
}

func (e *entry) Date() Date {
	return Date{
		Year: e.date.Year,
		Month: e.date.Month,
		Day: e.date.Day,
	}
}

func (e *entry) SetDate(date Date) error {
	d := civil.Date{
		Year: date.Year,
		Month: date.Month,
		Day: date.Day,
	}
	if !d.IsValid() {
		return &EntryError{ Code: INVALID_DATE }
	}
	e.date = d
	return nil
}

func (e *entry) Summary() string {
	return e.summary
}

func (e *entry) SetSummary(summary string) error {
	e.summary = summary
	return nil
}

func (e *entry) Times() []Minutes {
	return e.times
}

func (e *entry) AddTime(time Minutes) error {
	if time < 0 {
		return &EntryError{ Code: NEGATIVE_TIME }
	}
	e.times = append(e.times, time)
	return nil
}

func (e *entry) Ranges() [][]Time {
	ts := [][]Time{}
	for _, r := range e.ranges {
		ts = append(ts, []Time{
			Time{ Hour: r[0].Hour, Minute: r[0].Minute },
			Time{ Hour: r[1].Hour, Minute: r[1].Minute },
		})
	}
	return ts
}

func (e *entry) AddRange(start Time, end Time) error {
	e.ranges = append(e.ranges, []civil.Time{
		civil.Time{ Hour: start.Hour, Minute: start.Minute, Second: 0, Nanosecond: 0 },
		civil.Time{ Hour: end.Hour, Minute: end.Minute, Second: 0, Nanosecond: 0 },
	})
	return nil
}

func (e *entry) TotalTime() Minutes {
	total := Minutes(0)
	for _, t := range e.times {
		total += t
	}
	for _, rs := range e.ranges {
		start := rs[0].Minute + 60 * rs[0].Hour
		end := rs[1].Minute + 60 * rs[1].Hour
		total += Minutes(end - start)
	}
	return total
}
