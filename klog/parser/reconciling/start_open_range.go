package reconciling

import (
	"errors"
	"github.com/jotaen/klog/klog"
)

// StartOpenRange appends a new open range entry in a record.
func (r *Reconciler) StartOpenRange(startTime Styled[klog.Time], entrySummary klog.EntrySummary) (*Result, error) {
	if r.findOpenRangeIndex() != -1 {
		return nil, errors.New("There is already an open range in this record")
	}
	if startTime.AutoStyle {
		// Re-parse time to apply styling from reconciler.
		reformattedTime, err := klog.NewTimeFromString(startTime.Value.ToStringWithFormat(r.style.timeFormat()))
		startTime.Value = reformattedTime
		if err != nil {
			panic("INVALID_TIME")
		}
	}
	or := klog.NewOpenRangeWithFormat(startTime.Value, r.style.openRangeFormat())
	entryValue := or.ToString()
	r.insert(r.lastLinePointer, toMultilineEntryTexts(entryValue, entrySummary))
	return r.MakeResult()
}
