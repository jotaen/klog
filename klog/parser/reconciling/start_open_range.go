package reconciling

import (
	"errors"
	"github.com/jotaen/klog/klog"
)

// StartOpenRange appends a new open range entry in a record.
func (r *Reconciler) StartOpenRange(startTime klog.Time, entrySummary klog.EntrySummary) (*Result, error) {
	if r.findOpenRangeIndex() != -1 {
		return nil, errors.New("There is already an open range in this record")
	}
	// Re-parse time to apply styling from reconciler.
	reformattedTime, err := klog.NewTimeFromString(startTime.ToStringWithFormat(r.style.timeFormat()))
	if err != nil {
		panic("INVALID_TIME")
	}
	or := klog.NewOpenRangeWithFormat(reformattedTime, r.style.openRangeFormat())
	entryValue := or.ToString()
	r.insert(r.lastLinePointer, toMultilineEntryTexts(entryValue, entrySummary))
	return r.MakeResult()
}
