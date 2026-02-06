package service

import (
	gosort "sort"

	"github.com/jotaen/klog/klog"
)

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
