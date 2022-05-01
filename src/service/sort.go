package service

import (
	. "github.com/jotaen/klog/src"
	"sort"
)

// Sort orders the records by date.
func Sort(rs []Record, startWithOldest bool) []Record {
	sorted := append([]Record(nil), rs...)
	sort.Slice(sorted, func(i, j int) bool {
		isLess := sorted[j].Date().IsAfterOrEqual(sorted[i].Date())
		if !startWithOldest {
			return !isLess
		}
		return isLess
	})
	return sorted
}
