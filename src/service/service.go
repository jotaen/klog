package service

import (
	. "klog/record"
	"sort"
)

func Total(rs ...Record) Duration {
	total := NewDuration(0, 0)
	for _, r := range rs {
		for _, e := range r.Entries() {
			total = total.Plus(e.Duration())
		}
	}
	return total
}

func HypotheticalTotal(r Record, until Time) Duration {
	_ = r.EndOpenRange(until)
	return Total(r)
}

func ShouldTotal(rs ...Record) Duration {
	total := NewDuration(0, 0)
	for _, r := range rs {
		total = total.Plus(r.ShouldTotal())
	}
	return total
}

func TotalEntries(es []Entry) Duration {
	total := NewDuration(0, 0)
	for _, e := range es {
		total = total.Plus(e.Duration())
	}
	return total
}

type Filter struct {
	Tags     []string
	BeforeEq Date
	AfterEq  Date
}

func Sort(rs []Record, startWithOldest bool) []Record {
	sorted := append([]Record(nil), rs...)
	sort.Slice(sorted, func(i, j int) bool {
		return !startWithOldest || rs[j].Date().IsAfterOrEqual(rs[i].Date())
	})
	return sorted
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

func FindRelevantOpenRangeAt(rs []Record, date Date) []Record {
	var result []Record
	for _, r := range rs {
		if r.OpenRange() == nil {
			continue
		}
		if r.Date() == date || r.Date().PlusDays(-1) == date {
			result = append(result, r)
		}
	}
	return result
}
