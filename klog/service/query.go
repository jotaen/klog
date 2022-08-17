package service

import (
	"github.com/jotaen/klog/klog"
	gosort "sort"
)

// FilterQry represents the filter clauses of a query.
type FilterQry struct {
	Tags          []klog.Tag
	BeforeOrEqual klog.Date
	AfterOrEqual  klog.Date
	AtDate        klog.Date
}

// Filter returns all records the matches the query.
// A matching record must satisfy *all* query clauses.
func Filter(rs []klog.Record, o FilterQry) []klog.Record {
	var records []klog.Record
	for _, r := range rs {
		if o.AtDate != nil && !o.AtDate.IsEqualTo(r.Date()) {
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
func Sort(rs []klog.Record, startWithOldest bool) []klog.Record {
	sorted := append([]klog.Record(nil), rs...)
	gosort.Slice(sorted, func(i, j int) bool {
		isLess := sorted[j].Date().IsAfterOrEqual(sorted[i].Date())
		if !startWithOldest {
			return !isLess
		}
		return isLess
	})
	return sorted
}

func reduceRecordToMatchingTags(queriedTags []klog.Tag, r klog.Record) (klog.Record, bool) {
	if isSubsetOf(queriedTags, r.Summary().Tags()) {
		return r, true
	}
	var matchingEntries []klog.Entry
	for _, e := range r.Entries() {
		allTags := klog.Merge(r.Summary().Tags(), e.Summary().Tags())
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

func isSubsetOf(queriedTags []klog.Tag, allTags klog.TagSet) bool {
	for _, t := range queriedTags {
		if !allTags.Contains(t) {
			return false
		}
	}
	return true
}
