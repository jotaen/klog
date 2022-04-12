package service

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/service/period"
	gosort "sort"
)

// FilterQry represents the filter clauses of a query.
type FilterQry struct {
	Tags          []Tag
	BeforeOrEqual Date
	AfterOrEqual  Date
	Dates         []Date
}

// Filter returns all records the matches the query.
// A matching record must satisfy *all* query clauses.
func Filter(rs []Record, o FilterQry) []Record {
	dates := newDateSet(o.Dates)
	var records []Record
	for _, r := range rs {
		if len(dates) > 0 && !dates[period.NewDayFromDate(r.Date()).Hash()] {
			continue
		}
		if o.BeforeOrEqual != nil && !o.BeforeOrEqual.IsAfterOrEqual(r.Date()) {
			continue
		}
		if o.AfterOrEqual != nil && !r.Date().IsAfterOrEqual(o.AfterOrEqual) {
			continue
		}
		if len(o.Tags) > 0 {
			reducedR, hasMatched := reduceRecordToMatchingTags(o.Tags, r)
			if !hasMatched {
				continue
			}
			r = reducedR
		}
		records = append(records, r)
	}
	return records
}

// Sort orders the records by date.
func Sort(rs []Record, startWithOldest bool) []Record {
	sorted := append([]Record(nil), rs...)
	gosort.Slice(sorted, func(i, j int) bool {
		isLess := sorted[j].Date().IsAfterOrEqual(sorted[i].Date())
		if !startWithOldest {
			return !isLess
		}
		return isLess
	})
	return sorted
}

func reduceRecordToMatchingTags(queriedTags []Tag, r Record) (Record, bool) {
	if isSubsetOf(queriedTags, r.Summary().Tags()) {
		return r, true
	}
	var matchingEntries []Entry
	for _, e := range r.Entries() {
		allTags := Merge(r.Summary().Tags(), e.Summary().Tags())
		if isSubsetOf(queriedTags, allTags) {
			matchingEntries = append(matchingEntries, e)
		}
	}
	if len(matchingEntries) == 0 {
		return nil, false
	}
	r.SetEntries(matchingEntries)
	return r, true
}

func isSubsetOf(queriedTags []Tag, allTags TagSet) bool {
	for _, t := range queriedTags {
		if !allTags.Contains(t) {
			return false
		}
	}
	return true
}

func newDateSet(ds []Date) map[period.DayHash]bool {
	dict := make(map[period.DayHash]bool, len(ds))
	for _, d := range ds {
		dict[period.NewDayFromDate(d).Hash()] = true
	}
	return dict
}
