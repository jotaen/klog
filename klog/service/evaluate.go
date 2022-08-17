package service

import (
	"github.com/jotaen/klog/klog"
)

// Total calculates the overall time spent in records.
// It disregards open ranges.
func Total(rs ...klog.Record) klog.Duration {
	total := klog.NewDuration(0, 0)
	for _, r := range rs {
		for _, e := range r.Entries() {
			total = total.Plus(e.Duration())
		}
	}
	return total
}

// ShouldTotalSum calculates the overall should-total time of records.
func ShouldTotalSum(rs ...klog.Record) klog.ShouldTotal {
	total := klog.NewDuration(0, 0)
	for _, r := range rs {
		total = total.Plus(r.ShouldTotal())
	}
	return klog.NewShouldTotal(0, total.InMinutes())
}

// Diff calculates the difference between should-total and actual total
func Diff(should klog.ShouldTotal, actual klog.Duration) klog.Duration {
	return actual.Minus(should)
}
