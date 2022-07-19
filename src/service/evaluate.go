package service

import (
	. "github.com/jotaen/klog/src"
)

// Total calculates the overall time spent in records.
// It disregards open ranges.
func Total(rs ...Record) Duration {
	total := NewDuration(0, 0)
	for _, r := range rs {
		for _, e := range r.Entries() {
			total = total.Plus(e.Duration())
		}
	}
	return total
}

// ShouldTotalSum calculates the overall should-total time of records.
func ShouldTotalSum(rs ...Record) ShouldTotal {
	total := NewDuration(0, 0)
	for _, r := range rs {
		total = total.Plus(r.ShouldTotal())
	}
	return NewShouldTotal(0, total.InMinutes())
}

// Diff calculates the difference between should-total and actual total
func Diff(should ShouldTotal, actual Duration) Duration {
	return actual.Minus(should)
}
