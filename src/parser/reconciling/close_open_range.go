package reconciling

import (
	"errors"
	. "github.com/jotaen/klog/src"
	"regexp"
	"strings"
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

	additionalSummary = strings.ReplaceAll(additionalSummary, "\\n", "\n")
	summaryLines := strings.Split(additionalSummary, "\n")

	// Append additional summary text. Due to multiline entry summaries, that might
	// not be the same line as the time value.
	openRangeLastSummaryLineIndex := openRangeValueLineIndex + countLines([]Entry{r.record.Entries()[openRangeEntryIndex]}) - 1
	firstSummaryLine := summaryLines[0] // Index `0` will always exist
	if len(firstSummaryLine) > 0 {
		// If there is additional summary text, always prepend a space to delimit
		// the additional summary from either the time value or from an already
		// existing summary text.
		r.lines[openRangeLastSummaryLineIndex].Text += " "
	}
	r.lines[openRangeLastSummaryLineIndex].Text += firstSummaryLine

	if len(summaryLines) > 1 {
		var subsequentSummaryLines []insertableText
		for _, nextLine := range summaryLines[1:] {
			subsequentSummaryLines = append(subsequentSummaryLines, insertableText{nextLine, 2})
		}
		r.insert(openRangeLastSummaryLineIndex+1, subsequentSummaryLines)
	}

	return r.MakeResult()
}
