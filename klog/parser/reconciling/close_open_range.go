package reconciling

import (
	"errors"
	"github.com/jotaen/klog/klog"
	"regexp"
)

// CloseOpenRange tries to close the open time range.
func (r *Reconciler) CloseOpenRange(endTime klog.Time, additionalSummary klog.EntrySummary) (*Result, error) {
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
	r.lines[openRangeValueLineIndex].Text = regexp.MustCompile(`^(.*?)\?+(.*)$`).
		ReplaceAllString(
			r.lines[openRangeValueLineIndex].Text,
			"${1}"+endTime.ToStringWithFormat(r.style.timeFormat())+"${2}",
		)

	// Append additional summary text. Due to multiline entry summaries, that might
	// not be the same line as the time value.
	openRangeLastSummaryLineIndex := openRangeValueLineIndex + countLines([]klog.Entry{r.Record.Entries()[openRangeEntryIndex]}) - 1
	if len(additionalSummary) > 0 {
		if len(additionalSummary[0]) > 0 {
			// If there is additional summary text, always prepend a space to delimit
			// the additional summary from either the time value or from an already
			// existing summary text.
			r.lines[openRangeLastSummaryLineIndex].Text += " "
		}
		r.lines[openRangeLastSummaryLineIndex].Text += additionalSummary[0]
	}

	if len(additionalSummary) > 1 {
		var subsequentSummaryLines []insertableText
		for _, nextLine := range additionalSummary[1:] {
			subsequentSummaryLines = append(subsequentSummaryLines, insertableText{nextLine, 2})
		}
		r.insert(openRangeLastSummaryLineIndex+1, subsequentSummaryLines)
	}

	return r.MakeResult()
}
