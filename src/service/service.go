package service

import (
	. "klog"
	"sort"
	"time"
)

func Total(rs ...Record) Duration {
	total, _ := HypotheticalTotal(time.Time{}, rs...)
	return total
}

func HypotheticalTotal(until time.Time, rs ...Record) (Duration, bool) {
	total := NewDuration(0, 0)
	isCurrent := false
	void := time.Time{}
	thisDay := NewDateFromTime(until)
	theDayBefore := thisDay.PlusDays(-1)
	for _, r := range rs {
		for _, e := range r.Entries() {
			t := (e.Unbox(
				func(r Range) interface{} { return r.Duration() },
				func(d Duration) interface{} { return d },
				func(o OpenRange) interface{} {
					if until != void && (r.Date().IsEqualTo(thisDay) || r.Date().IsEqualTo(theDayBefore)) {
						end := NewTimeFromTime(until)
						if r.Date().IsEqualTo(theDayBefore) {
							end, _ = NewTimeTomorrow(end.Hour(), end.Minute())
						}
						tr, err := NewRange(o.Start(), end)
						if err == nil {
							isCurrent = true
							return tr.Duration()
						}
					}
					return NewDuration(0, 0)
				})).(Duration)
			total = total.Plus(t)
		}
	}
	return total, isCurrent
}

func ShouldTotalSum(rs ...Record) ShouldTotal {
	total := NewDuration(0, 0)
	for _, r := range rs {
		total = total.Plus(r.ShouldTotal())
	}
	return NewShouldTotal(0, total.InMinutes())
}

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
			if !hasTagMatchesAndReduce(tags, r) {
				continue
			}
		}
		records = append(records, r)
	}
	return records
}

func hasTagMatchesAndReduce(tags TagSet, r Record) bool {
	remainder := func() TagSet {
		rs := NewTagSet()
		matches := r.Summary().MatchTags(tags)
		for t := range tags {
			if !matches[t] {
				rs[t] = true
			}
		}
		return rs
	}()
	if len(remainder) == 0 {
		return true
	}
	var matchingEntries []Entry
	for _, e := range r.Entries() {
		matches := e.Summary().MatchTags(remainder)
		if len(matches) == len(remainder) {
			matchingEntries = append(matchingEntries, e)
		}
	}
	r.SetEntries(matchingEntries)
	return len(matchingEntries) > 0
}

func newDateSet(ds []Date) map[DateHash]bool {
	dict := make(map[DateHash]bool, len(ds))
	for _, d := range ds {
		dict[d.Hash()] = true
	}
	return dict
}
