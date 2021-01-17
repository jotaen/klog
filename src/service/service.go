package service

import (
	"errors"
	. "klog/record"
	"strings"
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

func Find(date Date, rs []Record) Record {
	for _, r := range rs {
		if r.Date() == date {
			return r
		}
	}
	return nil
}

func FindEntriesWithHashtags(tags map[string]bool, r Record) []Entry {
	if ContainsOneOfTags(tags, r.Summary().ToString()) {
		return r.Entries()
	}
	var matches []Entry
	for _, e := range r.Entries() {
		if ContainsOneOfTags(tags, e.Summary().ToString()) {
			matches = append(matches, e)
		}
	}
	return matches
}

func ContainsOneOfTags(tags map[string]bool, searchText string) bool {
	for _, t := range HashTagPattern.FindAllStringSubmatch(searchText, -1) {
		if tags[strings.ToLower(t[1])] == true {
			return true
		}
	}
	return false
}

func TagList(tags ...string) map[string]bool {
	result := map[string]bool{}
	for _, t := range tags {
		result[strings.ToLower(t)] = true
	}
	return result
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
