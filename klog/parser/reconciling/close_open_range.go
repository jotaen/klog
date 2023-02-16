package reconciling

import (
	"errors"
	"github.com/jotaen/klog/klog"
	"regexp"
)

// CloseOpenRange tries to close the open time range.
func (r *Reconciler) CloseOpenRange(endTime klog.Time, format ReformatDirective[klog.TimeFormat], additionalSummary klog.EntrySummary) (*Result, error) {
	openRangeEntryIndex := r.findOpenRangeIndex()
	if openRangeEntryIndex == -1 {
		return nil, errors.New("No open time range")
	}
	eErr := r.Record.EndOpenRange(endTime)
	if eErr != nil {
		return nil, errors.New("Start and end time must be in chronological order")
	}

	// Replace question mark with end time.
	openRangeValueLineIndex := r.lastLinePointer - countLines(r.Record.Entries()[openRangeEntryIndex:])
	endTimeValue := endTime.ToString()
	format.apply(r.style.timeFormat(), func(f klog.TimeFormat) {
		endTimeValue = endTime.ToStringWithFormat(f)
	})
	r.lines[openRangeValueLineIndex].Text = regexp.MustCompile(`^(.*?)\?+(.*)$`).
		ReplaceAllString(
			r.lines[openRangeValueLineIndex].Text,
			"${1}"+endTimeValue+"${2}",
		)

	r.concatenateSummary(openRangeEntryIndex, openRangeValueLineIndex, additionalSummary)
	return r.MakeResult()
}
