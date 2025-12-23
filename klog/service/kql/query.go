package kql

import (
	"github.com/jotaen/klog/klog"
)

func Query(p Predicate, rs []klog.Record) []klog.Record {
	var res []klog.Record
	for _, r := range rs {
		var es []klog.Entry
		for i, e := range r.Entries() {
			if p.Matches(queriedEntry{r, r.Entries()[i]}) {
				es = append(es, e)
			}
		}
		if len(es) == 0 {
			continue
		}
		r.SetEntries(es)
		res = append(res, r)
	}
	return res
}
