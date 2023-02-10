package service

import (
	"errors"
	"github.com/jotaen/klog/klog"
	gotime "time"
)

// CloseOpenRanges closes open ranges at the time of `endTime`. Returns an error
// if a range is not closeable at that point in time.
// This method alters the provided records!
// The bool return value indicates whether any open ranges have been closed.
func CloseOpenRanges(endTime gotime.Time, rs ...klog.Record) (bool, error) {
	thisDay := klog.NewDateFromGo(endTime)
	theDayBefore := thisDay.PlusDays(-1)
	hasClosedAnyRange := false
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
			return false, tErr
		}
		eErr := r.EndOpenRange(end)
		hasClosedAnyRange = true
		if eErr != nil {
			return false, errors.New("Encountered uncloseable open range")
		}
	}
	return hasClosedAnyRange, nil
}
