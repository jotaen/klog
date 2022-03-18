package reconciling

import (
	"errors"
	. "github.com/jotaen/klog/src"
)

// StartOpenRange appends a new open range entry in a record.
func (r *Reconciler) StartOpenRange(startTime Time, entrySummary string) (*Result, error) {
	if r.findOpenRangeIndex() != -1 {
		return nil, errors.New("There is already an open range in this record")
	}
	entryValue := startTime.ToStringWithFormat(r.style.TimeFormat.Get()) + r.style.SpacingInRange.Get() + "-" + r.style.SpacingInRange.Get() + "?"
	r.insert(r.lastLinePointer, toMultilineEntryTexts(entryValue, entrySummary))
	return r.MakeResult()
}
