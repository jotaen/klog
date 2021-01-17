package service

import (
	"errors"
	. "klog/record"
)

func Total(r Record) Duration {
	total := NewDuration(0, 0)
	for _, e := range r.Entries() {
		switch v := e.Value().(type) {
		case Duration:
			total = total.Add(v)
			break
		case Range:
			total = total.Add(v.Duration())
			break
		}
	}
	return total
}

func TotalEntries(es []Entry) Duration {
	total := NewDuration(0, 0)
	for _, e := range es {
		total = total.Add(e.Duration())
	}
	return total
}

type Filter struct {
	Tags     []string
	BeforeEq Date
	AfterEq  Date
}

func FindFilter(rs []Record, f Filter) ([]Record, []Entry) {
	tags := NewTagSet(f.Tags...)
	var records []Record
	var entries []Entry
	for _, r := range rs {
		if f.BeforeEq != nil && !f.BeforeEq.IsAfterOrEqual(r.Date()) {
			continue
		}
		if f.AfterEq != nil && !r.Date().IsAfterOrEqual(f.AfterEq) {
			continue
		}
		es := r.Entries()
		if len(tags) > 0 {
			matchingEs, hasMatched := FindEntriesWithHashtags(tags, r)
			if !hasMatched {
				continue
			}
			es = matchingEs
		}
		entries = append(entries, es...)
		records = append(records, r)
	}
	return records, entries
}

func FindEntriesWithHashtags(tags TagSet, r Record) ([]Entry, bool) {
	if ContainsOneOfTags(tags, r.Summary().ToString()) {
		return r.Entries(), true
	}
	var matches []Entry
	for _, e := range r.Entries() {
		if ContainsOneOfTags(tags, e.Summary().ToString()) {
			matches = append(matches, e)
		}
	}
	return matches, len(matches) > 0
}

func QuickStartAt(rs []Record, date Date, time Time) (Record, error) {
	var recordToAlter *Record
	for _, r := range rs {
		if r.Date() == date {
			recordToAlter = &r
		}
	}
	if recordToAlter == nil {
		r := NewRecord(date)
		recordToAlter = &r
	}
	(*recordToAlter).StartOpenRange(time, "")
	return *recordToAlter, nil
}

func QuickStopAt(rs []Record, date Date, time Time) (Record, error) {
	var recordToAlter *Record
	for _, r := range rs {
		if r.Date() == date && r.OpenRange() != nil {
			recordToAlter = &r
		}
	}
	if recordToAlter == nil {
		return nil, errors.New("NO_OPEN_RANGE")
	}
	newRange, err := NewRange((*recordToAlter).OpenRange().Start(), time)
	if err != nil {
		return nil, err
	}
	(*recordToAlter).AddRange(newRange, "") // TODO take over summary
	(*recordToAlter).StartOpenRange(time, "")
	return *recordToAlter, nil
}
