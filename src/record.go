package klog

import (
	"errors"
	"regexp"
)

type ShouldTotal Duration
type shouldTotal struct {
	Duration
}

type Record interface {
	Date() Date

	ShouldTotal() ShouldTotal
	SetShouldTotal(Duration)

	Summary() Summary
	SetSummary(string) error

	Entries() []Entry
	SetEntries([]Entry)
	AddDuration(Duration, Summary)
	AddRange(Range, Summary)
	OpenRange() OpenRange
	StartOpenRange(Time, Summary) error
	EndOpenRange(Time) error
}

func NewRecord(date Date) Record {
	return &record{
		date: date,
	}
}

func NewShouldTotal(hours int, minutes int) ShouldTotal {
	return shouldTotal{NewDuration(hours, minutes)}
}

func (s shouldTotal) ToString() string {
	return s.Duration.ToString() + "!"
}

type record struct {
	date        Date
	shouldTotal ShouldTotal
	summary     Summary
	entries     []Entry
}

func (r *record) Date() Date {
	return r.date
}

func (r *record) ShouldTotal() ShouldTotal {
	if r.shouldTotal == nil {
		return NewDuration(0, 0)
	}
	return r.shouldTotal
}

func (r *record) SetShouldTotal(t Duration) {
	r.shouldTotal = NewShouldTotal(0, t.InMinutes())
}

func (r *record) Summary() Summary {
	return r.summary
}

func (r *record) SetSummary(summary string) error {
	if regexp.MustCompile(`(^|\n) `).MatchString(summary) {
		return errors.New("MALFORMED_SUMMARY")
	}
	r.summary = Summary(summary)
	return nil
}

func (r *record) Entries() []Entry {
	return r.entries
}

func (r *record) SetEntries(es []Entry) {
	r.entries = es
}

func (r *record) AddDuration(d Duration, s Summary) {
	r.entries = append(r.entries, NewEntry(d, s))
}

func (r *record) AddRange(tr Range, s Summary) {
	r.entries = append(r.entries, NewEntry(tr, s))
}

func (r *record) OpenRange() OpenRange {
	for _, e := range r.entries {
		t, isOpenRange := e.value.(*openRange)
		if isOpenRange {
			return t
		}
	}
	return nil
}

func (r *record) StartOpenRange(t Time, s Summary) error {
	if r.OpenRange() != nil {
		return errors.New("DUPLICATE_OPEN_RANGE")
	}
	r.entries = append(r.entries, NewEntry(NewOpenRange(t), s))
	return nil
}

func (r *record) EndOpenRange(end Time) error {
	for i, e := range r.entries {
		t, isOpenRange := e.value.(*openRange)
		if isOpenRange {
			tr, err := NewRange(t.Start(), end)
			if err != nil {
				return err
			}
			r.entries[i] = NewEntry(tr, e.summary)
			return nil
		}
	}
	return errors.New("NO_OPEN_RANGE")
}
