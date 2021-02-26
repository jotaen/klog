package service

import (
	. "klog"
	"time"
)

// Total calculates the overall time spent in records.
// It disregards open ranges.
func Total(rs ...Record) Duration {
	total, _ := HypotheticalTotal(time.Time{}, rs...)
	return total
}

// TotalEntries calculates the overall of entries.
// It disregards open ranges.
func TotalEntries(es ...Entry) Duration {
	total := NewDuration(0, 0)
	for _, e := range es {
		total = total.Plus(e.Duration())
	}
	return total
}

// HypotheticalTotal calculates the overall total time of records,
// assuming all open ranges would be closed at the `until` time.
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

// ShouldTotalSum calculates the overall should total time of records.
func ShouldTotalSum(rs ...Record) ShouldTotal {
	total := NewDuration(0, 0)
	for _, r := range rs {
		total = total.Plus(r.ShouldTotal())
	}
	return NewShouldTotal(0, total.InMinutes())
}
