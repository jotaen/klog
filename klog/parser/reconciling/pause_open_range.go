package reconciling

import (
	"errors"
	"github.com/jotaen/klog/klog"
	"regexp"
	"strings"
)

// PauseOpenRange adds/extends a pause entry after an open-ended time range.
// If the next entry is a negative duration and has the same summary, its
// value is extended. Otherwise, a new entry is created.
// The pause duration must be negative and can’t be 0.
func (r *Reconciler) PauseOpenRange(pause klog.Duration, summary klog.EntrySummary) (*Result, error) {
	if pause.InMinutes() > 0 {
		return nil, errors.New("Invalid pause duration")
	}
	openRangeEntryIndex := r.findOpenRangeIndex()
	if openRangeEntryIndex == -1 {
		return nil, errors.New("No open time range")
	}
	openRangeEntryLastLineIndex := r.lastLinePointer - countLines(r.Record.Entries()[openRangeEntryIndex:])

	existingPause := func() klog.Duration {
		// The open range is the last entry.
		if openRangeEntryIndex == len(r.Record.Entries())-1 {
			return nil
		}
		nextEntry := r.Record.Entries()[openRangeEntryIndex+1]
		if !nextEntry.Summary().Equals(summary) {
			// Summaries don’t match.
			return nil
		}
		// Find next duration entry.
		pauseCandidate := klog.Unbox[klog.Duration](&nextEntry,
			func(r klog.Range) klog.Duration { return nil },
			func(d klog.Duration) klog.Duration { return d },
			func(or klog.OpenRange) klog.Duration { return nil },
		)
		// Only return it if it’s negative.
		if pauseCandidate.InMinutes() < 0 {
			return pauseCandidate
		}
		return nil
	}()

	// If there is no existing pause that matches, create a new entry underneath
	// the open range entry.
	if existingPause == nil {
		r.insert(openRangeEntryLastLineIndex+1, toMultilineEntryTexts(pause.ToString(), summary))
		return r.MakeResult()
	}

	// If there is an existing pause that matches, replace it’s duration
	// with the extended value.
	extendedPause := existingPause.Plus(pause)
	pauseEntryLineIndex := openRangeEntryLastLineIndex + countLines([]klog.Entry{r.Record.Entries()[openRangeEntryIndex]})
	durationPattern := regexp.MustCompile(`(-\w+)`)
	value := durationPattern.FindString(r.lines[pauseEntryLineIndex].Text)
	if value != "" {
		r.lines[pauseEntryLineIndex].Text = strings.Replace(r.lines[pauseEntryLineIndex].Text, value, extendedPause.ToString(), 1)
	}
	return r.MakeResult()
}
