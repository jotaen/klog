package service

import (
	"errors"
	gotime "time"

	"github.com/jotaen/klog/klog"
)

// CloseOpenRanges closes open ranges at the time of `endTime`. Returns an error
// if a range is not closeable at that point in time.
// This method alters the provided records!
// The `map` return value indicates which open ranges have been closed.
func CloseOpenRanges(endTime gotime.Time, rs ...klog.Record) (map[klog.Record]int, error) {
	thisDay := klog.NewDateFromGo(endTime)
	theDayBefore := thisDay.PlusDays(-1)
	closedEntries := map[klog.Record]int{}
	for _, r := range rs {
		if r.OpenRange() == nil {
			continue
		}
		end, tErr := func() (klog.Time, error) {
			end := klog.NewTimeFromGo(endTime)
			if r.Date().IsEqualTo(thisDay) {
				return end, nil
			}
			if r.Date().IsEqualTo(theDayBefore) {
				return end.Plus(klog.NewDuration(24, 0))
			}
			return nil, errors.New("Encountered uncloseable open range")
		}()
		if tErr != nil {
			return nil, tErr
		}
		closedEntries[r] = openRangeIndex(r)
		eErr := r.EndOpenRange(end)
		if eErr != nil {
			return nil, errors.New("Encountered uncloseable open range")
		}
	}
	return closedEntries, nil
}

func openRangeIndex(r klog.Record) int {
	openEndedEntry := -1
	for entryI, e := range r.Entries() {
		klog.Unbox[any](&e,
			func(klog.Range) any { return nil },
			func(klog.Duration) any { return nil },
			func(klog.OpenRange) any {
				openEndedEntry = entryI
				return nil
			},
		)
	}
	return openEndedEntry
}
