package reconciling

import (
	"errors"
	. "github.com/jotaen/klog/src"
	"regexp"
)

// CloseOpenRange tries to close the open time range.
func (r *Reconciler) CloseOpenRange(endTime Time, additionalSummary string) (*Result, error) {
	openRangeEntryIndex := r.findOpenRangeIndex()
	if openRangeEntryIndex == -1 {
		return nil, errors.New("No open time range")
	}
	eErr := r.record.EndOpenRange(endTime)
	if eErr != nil {
		return nil, errors.New("Start and end time must be in chronological order")
	}

	// Replace question mark with end time.
	openRangeValueLineIndex := r.lastLinePointer - countLines(r.record.Entries()[openRangeEntryIndex:])
	r.lines[openRangeValueLineIndex].Text = regexp.MustCompile(`^(.*?)\?+(.*)$`).
		ReplaceAllString(
			r.lines[openRangeValueLineIndex].Text,
			"${1}"+endTime.ToStringWithFormat(r.style.TimeFormat())+"${2}",
		)

	// Append additional summary text. Due to multiline entry summaries, that might
	// not be the same line as the time value.
	openRangeLastSummaryLineIndex := openRangeValueLineIndex + countLines([]Entry{r.record.Entries()[openRangeEntryIndex]}) - 1
	if len(additionalSummary) > 0 {
		// If there is additional summary text, always prepend a space to delimit
		// the additional summary from either the time value or from an already
		// existing summary text.
		additionalSummary = " " + additionalSummary
	}
	r.lines[openRangeLastSummaryLineIndex].Text += additionalSummary

	return r.MakeResult()
}
