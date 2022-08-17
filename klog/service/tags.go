package service

import (
	"github.com/jotaen/klog/klog"
	"sort"
)

type TagTotal struct {
	Tag     klog.Tag
	Total   klog.Duration
	forSort string
}

// AggregateTotalsByTags returns a map for looking up matching entries for a given tag.
func AggregateTotalsByTags(rs ...klog.Record) []TagTotal {
	result := make(totalByTag)
	for _, r := range rs {
		for _, e := range r.Entries() {
			alreadyCounted := make(map[klog.Tag]bool)
			allTags := klog.Merge(r.Summary().Tags(), e.Summary().Tags())
			for tag := range allTags {
				if alreadyCounted[tag] {
					continue
				}
				result.put(tag, e.Duration())
			}
		}
	}
	return result.toSortedList()
}

type totalByTag map[string]map[string]klog.Duration

func (tbt totalByTag) put(t klog.Tag, d klog.Duration) {
	if tbt[t.Name()] == nil {
		tbt[t.Name()] = make(map[string]klog.Duration)
	}

	if tbt[t.Name()][t.Value()] == nil {
		tbt[t.Name()][t.Value()] = klog.NewDuration(0, 0)
	}
	tbt[t.Name()][t.Value()] = tbt[t.Name()][t.Value()].Plus(d)
}

func (tbt totalByTag) toSortedList() []TagTotal {
	var result []TagTotal
	for tagName, totalsByValue := range tbt {
		for tagValue, total := range totalsByValue {
			result = append(result, TagTotal{
				forSort: tagName + "=" + tagValue,
				Tag:     klog.NewTagOrPanic(tagName, tagValue),
				Total:   total,
			})
		}
	}
	sort.Slice(result, func(i int, j int) bool {
		return result[i].forSort < result[j].forSort
	})
	return result
}
