package service

import (
	"github.com/jotaen/klog/klog"
	gosort "sort"
)

type EntryType string

const (
	ENTRY_TYPE_DURATION          = EntryType("DURATION")
	ENTRY_TYPE_POSITIVE_DURATION = EntryType("DURATION_POSITIVE")
	ENTRY_TYPE_NEGATIVE_DURATION = EntryType("DURATION_NEGATIVE")
	ENTRY_TYPE_RANGE             = EntryType("RANGE")
	ENTRY_TYPE_OPEN_RANGE        = EntryType("OPEN_RANGE")
)

// FilterQry represents the filter clauses of a query.
type FilterQry struct {
	Tags          []klog.Tag
	BeforeOrEqual klog.Date
	AfterOrEqual  klog.Date
	AtDate        klog.Date
	EntryType     EntryType
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
		if o.EntryType != "" {
			reducedR, hasMatched := reduceRecordToMatchingEntryTypes(o.EntryType, r)
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
		if isSubsetOf(queriedTags, &allTags) {
			matchingEntries = append(matchingEntries, e)
		}
	}
	if len(matchingEntries) == 0 {
		return nil, false
	}
	r.SetEntries(matchingEntries)
	return r, true
}

func reduceRecordToMatchingEntryTypes(t EntryType, r klog.Record) (klog.Record, bool) {
	var matchingEntries []klog.Entry
	for _, e := range r.Entries() {
		isMatch := klog.Unbox(&e, func(r klog.Range) bool {
			return t == ENTRY_TYPE_RANGE
		}, func(duration klog.Duration) bool {
			if t == ENTRY_TYPE_DURATION {
				return true
			} else if t == ENTRY_TYPE_POSITIVE_DURATION && e.Duration().InMinutes() >= 0 {
				return true
			} else if t == ENTRY_TYPE_NEGATIVE_DURATION && e.Duration().InMinutes() < 0 {
				return true
			}
			return false
		}, func(openRange klog.OpenRange) bool {
			return t == ENTRY_TYPE_OPEN_RANGE
		})
		if isMatch {
			matchingEntries = append(matchingEntries, e)
		}
	}
	if len(matchingEntries) == 0 {
		return nil, false
	}
	r.SetEntries(matchingEntries)
	return r, true
}

func isSubsetOf(queriedTags []klog.Tag, allTags *klog.TagSet) bool {
	for _, t := range queriedTags {
		if !allTags.Contains(t) {
			return false
		}
	}
	return true
}
