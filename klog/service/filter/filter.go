package filter

import (
	"github.com/jotaen/klog/klog"
)

// Filter goes through a list of records and only keeps those that match the
// given predicate. The records may be returned partially, keeping only those
// entries that match the predicate. The second return value indicates whether
// there are partial records with ShouldTotal set, as this may yield nonsensical
// results in a subsequent evaluation.
func Filter(p Predicate, rs []klog.Record) ([]klog.Record, bool) {
	var res []klog.Record
	hasPartialRecordsWithShouldTotal := false
	for _, r := range rs {
		if len(r.Entries()) == 0 && p.MatchesEmptyRecord(r) {
			res = append(res, r)
		} else {
			var es []klog.Entry
			for i, e := range r.Entries() {
				if p.Matches(r, r.Entries()[i]) {
					es = append(es, e)
				}
			}
			if len(es) == 0 {
				continue
			}
			if len(es) != len(r.Entries()) && r.ShouldTotal() != nil {
				hasPartialRecordsWithShouldTotal = true
			}
			r.SetEntries(es)
			res = append(res, r)
		}
	}
	return res, hasPartialRecordsWithShouldTotal
}
