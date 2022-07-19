package service

import (
	"errors"
	. "github.com/jotaen/klog/src"
	gotime "time"
)

// CloseOpenRanges closes open ranges at the time of `endTime`. Returns an error
// if a range is not closeable at that point in time.
// This method alters the provided records!
func CloseOpenRanges(endTime gotime.Time, rs ...Record) ([]Record, error) {
	thisDay := NewDateFromGo(endTime)
	theDayBefore := thisDay.PlusDays(-1)
	for _, r := range rs {
		if r.OpenRange() == nil {
			continue
		}
		end, tErr := func() (Time, error) {
			end := NewTimeFromGo(endTime)
			if r.Date().IsEqualTo(thisDay) {
				return end, nil
			}
			if r.Date().IsEqualTo(theDayBefore) {
				return end.Plus(NewDuration(24, 0))
			}
			return nil, errors.New("Encountered uncloseable open range")
		}()
		if tErr != nil {
			return nil, tErr
		}
		eErr := r.EndOpenRange(end)
		if eErr != nil {
			return nil, errors.New("Encountered uncloseable open range")
		}
	}
	return rs, nil
}
