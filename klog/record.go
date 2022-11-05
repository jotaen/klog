/*
Package klog is the implementation of the domain logic of klog.
It is essentially the code representation of the concepts as they are defined
in the file format specification.
*/
package klog

import (
	"errors"
)

// SPEC_VERSION contains the version number of the file format
// specification which this implementation is based on.
const SPEC_VERSION = "1.4"

// Record is a self-contained data container that holds the time tracking
// information associated with a certain date.
type Record interface {
	Date() Date

	ShouldTotal() ShouldTotal
	SetShouldTotal(Duration)

	Summary() RecordSummary
	SetSummary(RecordSummary)

	// Entries returns a list of all entries that are associated with this record.
	Entries() []Entry

	// SetEntries associates new entries with the record.
	SetEntries([]Entry)
	AddDuration(Duration, EntrySummary)
	AddRange(Range, EntrySummary)

	// OpenRange returns the open time range, or `nil` if there is none.
	OpenRange() OpenRange

	// Start starts a new open time range. It returns an error if there is
	// already an open time range present. (There can only be one per record.)
	Start(OpenRange, EntrySummary) error

	// EndOpenRange ends the open time range. It returns an error if there is
	// no open time range present, or if start and end time cannot be converted
	// into a valid time range.
	EndOpenRange(Time) error
}

func NewRecord(date Date) Record {
	return &record{
		date: date,
	}
}

type record struct {
	date        Date
	shouldTotal ShouldTotal
	summary     RecordSummary
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

func (r *record) Summary() RecordSummary {
	return r.summary
}

func (r *record) SetSummary(summary RecordSummary) {
	r.summary = summary
}

func (r *record) Entries() []Entry {
	return r.entries
}

func (r *record) SetEntries(es []Entry) {
	r.entries = es
}

func (r *record) AddDuration(d Duration, s EntrySummary) {
	r.entries = append(r.entries, NewEntryFromDuration(d, s))
}

func (r *record) AddRange(tr Range, s EntrySummary) {
	r.entries = append(r.entries, NewEntryFromRange(tr, s))
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

func (r *record) Start(or OpenRange, s EntrySummary) error {
	if r.OpenRange() != nil {
		return errors.New("DUPLICATE_OPEN_RANGE")
	}
	r.entries = append(r.entries, NewEntryFromOpenRange(or, s))
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
			r.entries[i] = NewEntryFromRange(tr, e.summary)
			return nil
		}
	}
	return errors.New("NO_OPEN_RANGE")
}
