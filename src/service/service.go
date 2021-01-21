package service

import (
	"klog"
	"sort"
)

func Total(rs ...src.Record) src.Duration {
	total := src.NewDuration(0, 0)
	for _, r := range rs {
		for _, e := range r.Entries() {
			total = total.Plus(e.Duration())
		}
	}
	return total
}

func HypotheticalTotal(r src.Record, until src.Time) src.Duration {
	total := Total(r)
	or := r.OpenRange()
	if or != nil {
		tr, err := src.NewRange(or.Start(), until)
		if err == nil {
			total = total.Plus(tr.Duration())
		}
	}
	return total
}

func ShouldTotal(rs ...src.Record) src.Duration {
	total := src.NewDuration(0, 0)
	for _, r := range rs {
		total = total.Plus(r.ShouldTotal())
	}
	return total
}

func TotalEntries(es []src.Entry) src.Duration {
	total := src.NewDuration(0, 0)
	for _, e := range es {
		total = total.Plus(e.Duration())
	}
	return total
}

type Filter struct {
	Tags     []string
	BeforeEq src.Date
	AfterEq  src.Date
	Dates    []src.Date
}

func Sort(rs []src.Record, startWithOldest bool) []src.Record {
	sorted := append([]src.Record(nil), rs...)
	sort.Slice(sorted, func(i, j int) bool {
		return !startWithOldest || rs[j].Date().IsAfterOrEqual(rs[i].Date())
	})
	return sorted
}

func FindFilter(rs []src.Record, f Filter) ([]src.Record, []src.Entry) {
	tags := src.NewTagSet(f.Tags...)
	dates := newDateSet(f.Dates)
	var records []src.Record
	var entries []src.Entry
	for _, r := range rs {
		if len(dates) > 0 && !dates[r.Date().ToString()] {
			continue
		}
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

func newDateSet(ds []src.Date) map[string]bool {
	dict := make(map[string]bool, len(ds))
	for _, d := range ds {
		dict[d.ToString()] = true
	}
	return dict
}

func FindEntriesWithHashtags(tags src.TagSet, r src.Record) ([]src.Entry, bool) {
	if src.ContainsOneOfTags(tags, r.Summary().ToString()) {
		return r.Entries(), true
	}
	var matches []src.Entry
	for _, e := range r.Entries() {
		if src.ContainsOneOfTags(tags, e.Summary().ToString()) {
			matches = append(matches, e)
		}
	}
	return matches, len(matches) > 0
}

func FindRelevantOpenRangeAt(rs []src.Record, date src.Date) []src.Record {
	var result []src.Record
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
