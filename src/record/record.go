package record

import (
	"errors"
)

type Summary string

type Record interface {
	Date() Date

	ShouldTotal() Duration
	SetShouldTotal(Duration)

	Summary() string
	SetSummary(string) // TODO no blank lines; no indentation

	Entries() []Entry
	AddDuration(Duration, Summary)
	AddRange(Range, Summary)
	OpenRange() OpenRangeStart
	StartOpenRange(OpenRangeStart, Summary)
	EndOpenRange(Time) error

	//FindEntriesWithHashtags([]string) []Entry // TODO
}

func NewRecord(date Date) Record {
	return &record{
		date: date,
	}
}

type record struct {
	date        Date
	shouldTotal Duration
	summary     string
	entries     []Entry
}

func (r *record) Date() Date {
	return r.date
}

func (r *record) ShouldTotal() Duration {
	return r.shouldTotal
}

func (r *record) SetShouldTotal(t Duration) {
	r.shouldTotal = t
}

func (r *record) Summary() string {
	return r.summary
}

func (r *record) SetSummary(summary string) {
	r.summary = summary
}

func (r *record) Entries() []Entry {
	return r.entries
}

func (r *record) Durations() []Duration {
	var durations []Duration
	for _, e := range r.entries {
		d, isDuration := e.Value().(Duration)
		if isDuration {
			durations = append(durations, d)
		}
	}
	return durations
}

func (r *record) AddDuration(d Duration, s Summary) {
	r.entries = append(r.entries, Entry{value: d, summary: s})
}

func (r *record) Ranges() []Range {
	var ranges []Range
	for _, e := range r.entries {
		tr, isRange := e.Value().(Range)
		if isRange {
			ranges = append(ranges, tr)
		}
	}
	return ranges
}

func (r *record) AddRange(tr Range, s Summary) {
	r.entries = append(r.entries, Entry{value: tr, summary: s})
}

func (r *record) OpenRange() OpenRangeStart {
	for _, e := range r.entries {
		t, isStartTime := e.Value().(Time)
		if isStartTime {
			return t
		}
	}
	return nil
}

func (r *record) StartOpenRange(t OpenRangeStart, s Summary) {
	r.entries = append(r.entries, Entry{value: t, summary: s})
}

func (r *record) EndOpenRange(end Time) error {
	for i, e := range r.entries {
		t, isStartTime := e.Value().(Time)
		if isStartTime {
			tr, err := NewRange(t, end)
			if err != nil {
				return err
			}
			r.entries[i] = Entry{value: tr, summary: ""}
			return nil
		}
	}
	return errors.New("NO_OPEN_RANGE")
}
