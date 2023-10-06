package reconciling

import (
	"errors"
	"github.com/jotaen/klog/klog"
	"regexp"
	"strings"
)

// AppendPause adds a new pause entry to a record that contains an open range.
func (r *Reconciler) AppendPause(summary klog.EntrySummary) error {
	if r.findOpenRangeIndex() == -1 {
		return errors.New("No open time range found")
	}
	entryValue := "-0m"
	if len(summary) == 0 {
		summary, _ = klog.NewEntrySummary("")
	}
	if len(summary[0]) > 0 {
		entryValue += " "
	}
	summary[0] = entryValue + summary[0]
	return r.AppendEntry(summary)
}

// ExtendPause extends an existing pause entry.
func (r *Reconciler) ExtendPause(increment klog.Duration, additionalSummary klog.EntrySummary) error {
	if r.findOpenRangeIndex() == -1 {
		return errors.New("No open time range found")
	}

	pauseEntryI := r.findLastEntry(func(e klog.Entry) bool {
		return klog.Unbox[bool](&e, func(_ klog.Range) bool {
			return false
		}, func(d klog.Duration) bool {
			return d.InMinutes() <= 0
		}, func(_ klog.OpenRange) bool {
			return false
		})
	})
	if pauseEntryI == -1 {
		return errors.New("Could not find existing pause to extend")
	}

	extendedPause := r.Record.Entries()[pauseEntryI].Duration().Plus(increment)
	pauseLineIndex := r.lastLinePointer - countLines(r.Record.Entries()[pauseEntryI:])
	durationPattern := regexp.MustCompile(`(-\w+)`)
	value := durationPattern.FindString(r.lines[pauseLineIndex].Text)
	if extendedPause.InMinutes() != 0 {
		r.lines[pauseLineIndex].Text = strings.Replace(r.lines[pauseLineIndex].Text, value, extendedPause.ToString(), 1)
	}

	r.concatenateSummary(pauseEntryI, pauseLineIndex, additionalSummary)
	return nil
}
