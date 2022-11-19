package service

import (
	"github.com/jotaen/klog/klog"
	"sort"
)

type TagStats struct {
	Tag klog.Tag

	// Total is the total duration alloted to the tag.
	Total klog.Duration

	// Count is the total number of matching entries for that tag.
	// I.e., this is *not* how often a tag appears in the record text.
	Count int

	keyForSort string
}

// AggregateTotalsByTags returns a list of tags (sorted by tag, alphanumerically)
// that contains statistics about the tags appearing in the data.
func AggregateTotalsByTags(rs ...klog.Record) []*TagStats {
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

// Structure: "tagName":"tagValue":TagStats
type totalByTag map[string]map[string]*TagStats

func (tbt totalByTag) put(t klog.Tag, d klog.Duration) {
	if tbt[t.Name()] == nil {
		tbt[t.Name()] = make(map[string]*TagStats)
	}

	if tbt[t.Name()][t.Value()] == nil {
		tbt[t.Name()][t.Value()] = &TagStats{
			Tag:        t,
			Total:      klog.NewDuration(0, 0),
			Count:      0,
			keyForSort: t.Name() + "=" + t.Value(),
		}
	}

	stats := tbt[t.Name()][t.Value()]
	stats.Total = stats.Total.Plus(d)
	stats.Count++
}

func (tbt totalByTag) toSortedList() []*TagStats {
	var result []*TagStats
	for _, ts := range tbt {
		for _, t := range ts {
			result = append(result, t)
		}
	}
	sort.Slice(result, func(i int, j int) bool {
		return result[i].keyForSort < result[j].keyForSort
	})
	return result
}
