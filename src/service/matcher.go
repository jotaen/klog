package service

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/service/period"
)

type Matcher interface {
	Apply(Record) Record
}

type atDateMatcher struct {
	date Date
}

func (m *atDateMatcher) Apply(r Record) Record {
	if r.Date().IsEqualTo(m.date) {
		return r
	}
	return nil
}

type upToDateMatcher struct {
	date Date
}

func (m *upToDateMatcher) Apply(r Record) Record {
	if m.date.IsAfterOrEqual(r.Date()) {
		return r
	}
	return nil
}

type fromDateMatcher struct {
	date Date
}

func (m *fromDateMatcher) Apply(r Record) Record {
	if r.Date().IsAfterOrEqual(m.date) {
		return r
	}
	return nil
}

type inPeriodMatcher struct {
	period period.Period
}

func (m *inPeriodMatcher) Apply(r Record) Record {
	if r.Date().IsAfterOrEqual(m.period.Since()) && m.period.Until().IsAfterOrEqual(r.Date()) {
		return r
	}
	return nil
}

type tagMatcher struct {
	tag Tag
}

func (m *tagMatcher) Apply(r Record) Record {
	reducedR, hasMatched := reduceRecordToMatchingTags([]Tag{m.tag}, r)
	if !hasMatched {
		return nil
	}
	return reducedR
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

type andMatcher struct {
	left  Matcher
	right Matcher
}

func (m *andMatcher) Apply(r Record) Record {
	r1 := m.left.Apply(r)
	if r1 == nil {
		return nil
	}
	return m.right.Apply(r1)
}

//lint:ignore U1000 Ignore unused code
type orMatcher struct {
	left  Matcher
	right Matcher
}

func (m *orMatcher) Apply(r Record) Record {
	r1 := m.left.Apply(r)
	if r1 != nil {
		return r1
	}
	return m.right.Apply(r)
}

//lint:ignore U1000 Ignore unused code
type notMatcher struct {
	matcher Matcher
}

func (m *notMatcher) Apply(r Record) Record {
	if m.matcher.Apply(r) == nil {
		return r
	}
	return nil
}

type identityMatcher struct{}

func (m *identityMatcher) Apply(r Record) Record {
	return r
}
