package service

import (
	"klog"
	"sort"
	"time"
)

func Total(rs ...src.Record) src.Duration {
	total, _ := HypotheticalTotal(time.Time{}, rs...)
	return total
}

func HypotheticalTotal(until time.Time, rs ...src.Record) (src.Duration, bool) {
	total := src.NewDuration(0, 0)
	isCurrent := false
	void := time.Time{}
	thisDay := src.NewDateFromTime(until)
	theDayBefore := thisDay.PlusDays(-1)
	for _, r := range rs {
		for _, e := range r.Entries() {
			t := (e.Unbox(
				func(r src.Range) interface{} { return r.Duration() },
				func(d src.Duration) interface{} { return d },
				func(o src.OpenRange) interface{} {
					if until != void && (r.Date().IsEqualTo(thisDay) || r.Date().IsEqualTo(theDayBefore)) {
						end := src.NewTimeFromTime(until)
						if r.Date().IsEqualTo(theDayBefore) {
							end, _ = src.NewTimeTomorrow(end.Hour(), end.Minute())
						}
						tr, err := src.NewRange(o.Start(), end)
						if err == nil {
							isCurrent = true
							return tr.Duration()
						}
					}
					return src.NewDuration(0, 0)
				})).(src.Duration)
			total = total.Plus(t)
		}
	}
	return total, isCurrent
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

func dateHash(d src.Date) int {
	return d.Year()*(12*31) + d.Month()*31 + d.Day()
}

func FindFilter(rs []src.Record, f Filter) ([]src.Record, []src.Entry) {
	tags := src.NewTagSet(f.Tags...)
	dates := newDateSet(f.Dates)
	var records []src.Record
	var entries []src.Entry
	for _, r := range rs {
		if len(dates) > 0 && !dates[dateHash(r.Date())] {
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

func newDateSet(ds []src.Date) map[int]bool {
	dict := make(map[int]bool, len(ds))
	for _, d := range ds {
		dict[dateHash(d)] = true
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
