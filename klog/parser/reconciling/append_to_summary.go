package reconciling

import (
	"errors"
	"github.com/jotaen/klog/klog"
)

// AppendToOpenRange adds text to the summary of the open time range.
func (r *Reconciler) AppendToOpenRange(additionalSummary klog.EntrySummary) error {
	openRangeEntryIndex := r.findOpenRangeIndex()
	if openRangeEntryIndex == -1 {
		return errors.New("No open time range")
	}
	r.appendToSummary(openRangeEntryIndex, additionalSummary)
	return nil
}

// AppendToNthEntry adds text to the summary of the nth entry. If nth is positive, it
// counts from the start (beginning with 1); if negative, it counts from the end
// (beginning with -1).
func (r *Reconciler) AppendToNthEntry(nth int, additionalSummary klog.EntrySummary) error {
	entriesCount := len(r.Record.Entries())
	entryIndex := nth - 1
	if nth < 0 {
		entryIndex = entriesCount + nth
	}
	if nth == 0 || entryIndex < 0 || entryIndex >= entriesCount {
		return errors.New("No such entry")
	}
	r.appendToSummary(entryIndex, additionalSummary)
	return nil
}

func (r *Reconciler) appendToSummary(entryIndex int, additionalSummary klog.EntrySummary) {
	entryValueLineIndex := r.lastLinePointer - countLines(r.Record.Entries()[entryIndex:])
	r.concatenateSummary(entryIndex, entryValueLineIndex, additionalSummary)
}
