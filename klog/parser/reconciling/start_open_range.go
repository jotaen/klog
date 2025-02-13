package reconciling

import (
	"errors"
	"github.com/jotaen/klog/klog"
)

// StartOpenRange appends a new open range entry in a record.
func (r *Reconciler) StartOpenRange(startTime klog.Time, format ReformatDirective[klog.TimeFormat], entrySummary klog.EntrySummary) error {
	if r.findOpenRangeIndex() != -1 {
		return errors.New("There is already an open range in this record")
	}
	format.apply(r.style.timeFormat(), func(f klog.TimeFormat) {
		// Re-parse time to apply format.
		reformattedTime, err := klog.NewTimeFromString(startTime.ToStringWithFormat(f))
		if err != nil {
			panic("Invalid time")
		}
		startTime = reformattedTime
	})
	or := klog.NewOpenRangeWithFormat(startTime, r.style.openRangeFormat())
	entryValue := or.ToString()
	r.insert(r.lastLinePointer, toMultilineEntryTexts(entryValue, entrySummary))
	return nil
}
