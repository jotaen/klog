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
	newEntryLine := startTime.ToStringWithFormat(r.style.TimeFormat()) + r.style.SpacingInRange() + "-" + r.style.SpacingInRange() + "?"
	if len(entrySummary) > 0 {
		newEntryLine += " " + entrySummary
	}
	return r.AppendEntry(newEntryLine)
}
