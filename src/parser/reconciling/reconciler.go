/*
Package reconciling contains logic to manipulate klog source text.
The idea of the reconciler generally is to add or modify serialised records
in a minimally invasive manner. Instead or re-serialising the record itself,
it tries to find the location in the original text and modify that directly.
While this approach might feel a little hacky, it avoids lots of other
complications and sources of bugs, that could potentially mess up user data.
*/
package reconciling

import (
	"errors"
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/engine"
	"regexp"
)

// Reconciler is a mechanism to manipulate record data in a file.
type Reconciler struct {
	record          Record
	style           *parser.Style
	lastLinePointer int // Line index of the last entry
	lines           []engine.Line
	recordPointer   int
}

// Result is the result of an applied reconciler.
type Result struct {
	Record        Record
	AllRecords    []parser.ParsedRecord
	AllSerialised string
}

// Reconcile is a function interface for applying a reconciler.
type Reconcile func(*Reconciler) (*Result, error)

// AppendEntry adds a new entry to the end of the record.
func (r *Reconciler) AppendEntry(newEntry string) (*Result, error) {
	r.insert(r.lastLinePointer, []insertableText{{newEntry, 1}})
	return r.MakeResult()
}

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

func countLines(es []Entry) int {
	result := 0
	for _, e := range es {
		result += len(e.Summary())
	}
	return result
}

// MakeResult returns the reconciled data.
func (r *Reconciler) MakeResult() (*Result, error) {
	text := ""
	for _, l := range r.lines {
		text += l.Original()
	}

	// As a safeguard, make sure the result is parseable.
	newRecords, errs := parser.Parse(text)
	if errs != nil {
		return nil, errors.New("This operation wouldn’t result in a valid record")
	}

	return &Result{
		Record:        newRecords[r.recordPointer],
		AllRecords:    newRecords,
		AllSerialised: text,
	}, nil
}

// findOpenRangeIndex returns the index of the open range entry, or -1 if no open range.
func (r *Reconciler) findOpenRangeIndex() int {
	openRangeEntryIndex := -1
	for i, e := range r.record.Entries() {
		e.Unbox(
			func(Range) interface{} { return nil },
			func(Duration) interface{} { return nil },
			func(OpenRange) interface{} {
				openRangeEntryIndex = i
				return nil
			},
		)
	}
	return openRangeEntryIndex
}

var blankLine = insertableText{"", 0}

type insertableText struct {
	text        string
	indentation int
}

func (r *Reconciler) insert(lineIndex int, texts []insertableText) {
	result := make([]engine.Line, len(r.lines)+len(texts))
	offset := 0
	for i := range result {
		if i >= lineIndex && offset < len(texts) {
			line := ""
			if texts[offset].indentation > 0 {
				line += r.style.Indentation()
			}
			line += texts[offset].text + r.style.LineEnding()
			result[i] = engine.NewLineFromString(line, -1)
			offset++
		} else {
			result[i] = r.lines[i-offset]
		}
		result[i].LineNumber = i + 1
	}
	if lineIndex > 0 && result[lineIndex-1].LineEnding == "" {
		result[lineIndex-1].LineEnding = r.style.LineEnding()
	}
	r.lines = result
}
