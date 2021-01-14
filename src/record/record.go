package record

import (
	"errors"
	"klog/datetime"
)

type Record interface {
	Date() datetime.Date

	Summary() string
	SetSummary(string)

	Durations() []datetime.Duration
	AddDuration(datetime.Duration)

	Ranges() []datetime.TimeRange
	AddRange(datetime.TimeRange)

	OpenRange() datetime.Time
	StartOpenRange(datetime.Time)
	EndOpenRange(datetime.Time) error
}

func NewRecord(date datetime.Date) Record {
	return &record{
		date: date,
	}
}

type record struct {
	date           datetime.Date
	summary        string
	entries        []interface{}
}

func (r *record) Date() datetime.Date {
	return r.date
}

func (r *record) Summary() string {
	return r.summary
}

func (r *record) SetSummary(summary string) {
	r.summary = summary
}

func (r *record) Durations() []datetime.Duration {
	var durations []datetime.Duration
	for _, e := range r.entries {
		d, isDuration := e.(datetime.Duration)
		if isDuration {
			durations = append(durations, d)
		}
	}
	return durations
}

func (r *record) AddDuration(d datetime.Duration) {
	r.entries = append(r.entries, d)
}

func (r *record) Ranges() []datetime.TimeRange {
	var ranges []datetime.TimeRange
	for _, e := range r.entries {
		tr, isTimeRange := e.(datetime.TimeRange)
		if isTimeRange {
			ranges = append(ranges, tr)
		}
	}
	return ranges
}

func (r *record) AddRange(tr datetime.TimeRange) {
	r.entries = append(r.entries, tr)
}

func (r *record) OpenRange() datetime.Time {
	for _, e := range r.entries {
		t, isStartTime := e.(datetime.Time)
		if isStartTime {
			return t
		}
	}
	return nil
}

func (r *record) StartOpenRange(t datetime.Time) {
	r.entries = append(r.entries, t)
}

func (r *record) EndOpenRange(end datetime.Time) error {
	for i, e := range r.entries {
		t, isStartTime := e.(datetime.Time)
		if isStartTime {
			tr, err := datetime.NewTimeRange(t, end)
			if err != nil {
				return err
			}
			r.entries[i] = tr
			return nil
		}
	}
	return errors.New("NO_OPEN_RANGE")
}
