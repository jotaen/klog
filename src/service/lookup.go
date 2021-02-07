package service

import (
	. "klog"
	"sort"
)

type Filter struct {
	Tags     []string
	BeforeEq Date
	AfterEq  Date
	Dates    []Date
}

func Sort(rs []Record, startWithOldest bool) []Record {
	sorted := append([]Record(nil), rs...)
	sort.Slice(sorted, func(i, j int) bool {
		isLess := sorted[j].Date().IsAfterOrEqual(sorted[i].Date())
		if startWithOldest {
			return !isLess
		}
		return isLess
	})
	return sorted
}

func FindFilter(rs []Record, f Filter) []Record {
	tags := NewTagSet(f.Tags...)
	dates := newDateSet(f.Dates)
	var records []Record
	for _, r := range rs {
		if len(dates) > 0 && !dates[r.Date().Hash()] {
			continue
		}
		if f.BeforeEq != nil && !f.BeforeEq.IsAfterOrEqual(r.Date()) {
			continue
		}
		if f.AfterEq != nil && !r.Date().IsAfterOrEqual(f.AfterEq) {
			continue
		}
		if len(tags) > 0 {
			reducedR, hasMatched := reduceRecordToMatchingTags(tags, r)
			if !hasMatched {
				continue
			}
			r = reducedR
		}
		records = append(records, r)
	}
	return records
}

func reduceRecordToMatchingTags(tags TagSet, r Record) (Record, bool) {
	if isSubsetOf(tags, r.Summary().Tags()) {
		return r, true
	}
	_, tagsByEntry := EntryTagLookup(r)
	var matchingEntries []Entry
	for _, e := range r.Entries() {
		if isSubsetOf(tags, tagsByEntry[e]) {
			matchingEntries = append(matchingEntries, e)
		}
	}
	if len(matchingEntries) == 0 {
		return nil, false
	}
	r.SetEntries(matchingEntries)
	return r, true
}

func isSubsetOf(sub TagSet, super TagSet) bool {
	for t := range sub {
		if !super[t] {
			return false
		}
	}
	return true
}

func EntryTagLookup(rs ...Record) (map[Tag][]Entry, map[Entry]TagSet) {
	entriesByTag := make(map[Tag][]Entry)
	tagsByEntry := make(map[Entry]TagSet)
	for _, r := range rs {
		alreadyAdded := make(map[Tag]bool)
		for t := range r.Summary().Tags() {
			entriesByTag[t] = append(entriesByTag[t], r.Entries()...)
			alreadyAdded[t] = true
		}
		for _, e := range r.Entries() {
			tagsByEntry[e] = func() TagSet {
				result := r.Summary().Tags()
				for t := range e.Summary().Tags() {
					result[t] = true
				}
				return result
			}()
			for t := range e.Summary().Tags() {
				if alreadyAdded[t] {
					continue
				}
				entriesByTag[t] = append(entriesByTag[t], e)
			}
		}
	}
	return entriesByTag, tagsByEntry
}

func newDateSet(ds []Date) map[DateHash]bool {
	dict := make(map[DateHash]bool, len(ds))
	for _, d := range ds {
		dict[d.Hash()] = true
	}
	return dict
}
