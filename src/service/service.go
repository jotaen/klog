package service

import (
	"klog"
	"sort"
	"time"
)

func Total(rs ...klog.Record) klog.Duration {
	total, _ := HypotheticalTotal(time.Time{}, rs...)
	return total
}

func HypotheticalTotal(until time.Time, rs ...klog.Record) (klog.Duration, bool) {
	total := klog.NewDuration(0, 0)
	isCurrent := false
	void := time.Time{}
	thisDay := klog.NewDateFromTime(until)
	theDayBefore := thisDay.PlusDays(-1)
	for _, r := range rs {
		for _, e := range r.Entries() {
			t := (e.Unbox(
				func(r klog.Range) interface{} { return r.Duration() },
				func(d klog.Duration) interface{} { return d },
				func(o klog.OpenRange) interface{} {
					if until != void && (r.Date().IsEqualTo(thisDay) || r.Date().IsEqualTo(theDayBefore)) {
						end := klog.NewTimeFromTime(until)
						if r.Date().IsEqualTo(theDayBefore) {
							end, _ = klog.NewTimeTomorrow(end.Hour(), end.Minute())
						}
						tr, err := klog.NewRange(o.Start(), end)
						if err == nil {
							isCurrent = true
							return tr.Duration()
						}
					}
					return klog.NewDuration(0, 0)
				})).(klog.Duration)
			total = total.Plus(t)
		}
	}
	return total, isCurrent
}

func ShouldTotal(rs ...klog.Record) klog.ShouldTotal {
	total := klog.NewDuration(0, 0)
	for _, r := range rs {
		total = total.Plus(r.ShouldTotal())
	}
	return klog.NewShouldTotal(0, total.InMinutes())
}

func TotalEntries(es []klog.Entry) klog.Duration {
	total := klog.NewDuration(0, 0)
	for _, e := range es {
		total = total.Plus(e.Duration())
	}
	return total
}

type Filter struct {
	Tags     []string
	BeforeEq klog.Date
	AfterEq  klog.Date
	Dates    []klog.Date
}

func Sort(rs []klog.Record, startWithOldest bool) []klog.Record {
	sorted := append([]klog.Record(nil), rs...)
	sort.Slice(sorted, func(i, j int) bool {
		isLess := sorted[j].Date().IsAfterOrEqual(sorted[i].Date())
		if startWithOldest {
			return !isLess
		}
		return isLess
	})
	return sorted
}

func dateHash(d klog.Date) int {
	return d.Year()*(12*31) + d.Month()*31 + d.Day()
}

func FindFilter(rs []klog.Record, f Filter) ([]klog.Record, []klog.Entry) {
	tags := klog.NewTagSet(f.Tags...)
	dates := newDateSet(f.Dates)
	var records []klog.Record
	var entries []klog.Entry
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

func newDateSet(ds []klog.Date) map[int]bool {
	dict := make(map[int]bool, len(ds))
	for _, d := range ds {
		dict[dateHash(d)] = true
	}
	return dict
}

func FindEntriesWithHashtags(tags klog.TagSet, r klog.Record) ([]klog.Entry, bool) {
	if klog.ContainsOneOfTags(tags, r.Summary().ToString()) {
		return r.Entries(), true
	}
	var matches []klog.Entry
	for _, e := range r.Entries() {
		if klog.ContainsOneOfTags(tags, e.Summary().ToString()) {
			matches = append(matches, e)
		}
	}
	return matches, len(matches) > 0
}
