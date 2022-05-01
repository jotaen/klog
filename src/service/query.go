package service

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/service/period"
)

// Filter returns all records that satisfy the given matcher(s).
// The records might be reduced to the matching entries.
func Filter(matcher Matcher, rs []Record) []Record {
	var records []Record
	for _, r := range rs {
		reducedRecord := matcher.Apply(r)
		if reducedRecord != nil {
			records = append(records, reducedRecord)
		}
	}
	return records
}

type Query struct {
	AtDate   Date
	UpToDate Date
	FromDate Date
	InPeriod []period.Period
	WithTags []Tag
}

func (q *Query) ToMatcher() Matcher {
	var result Matcher = &identityMatcher{}
	if q.AtDate != nil {
		result = &andMatcher{result, &atDateMatcher{q.AtDate}}
	}
	if q.UpToDate != nil {
		result = &andMatcher{result, &upToDateMatcher{q.UpToDate}}
	}
	if q.FromDate != nil {
		result = &andMatcher{result, &fromDateMatcher{q.FromDate}}
	}
	for _, p := range q.InPeriod {
		result = &andMatcher{result, &inPeriodMatcher{p}}
	}
	for _, t := range q.WithTags {
		result = &andMatcher{result, &tagMatcher{t}}
	}
	return result
}
