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
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/parser"
	"github.com/jotaen/klog/klog/parser/txt"
	"strings"
)

// Reconciler is a mechanism to manipulate record data in a file.
type Reconciler struct {
	Record          klog.Record
	style           *style
	lastLinePointer int // Line index of the last entry
	lines           []txt.Line
	recordPointer   int
}

// Result is the result of an applied reconciler.
type Result struct {
	Record        klog.Record
	AllRecords    []klog.Record
	AllSerialised string
}

// Reconcile is a function interface for applying a reconciler.
type Reconcile func(*Reconciler) (*Result, error)

func countLines(es []klog.Entry) int {
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
	newRecords, _, errs := parser.NewSerialParser().Parse(text)
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
	return r.findLastEntry(func(e klog.Entry) bool {
		return klog.Unbox[bool](&e,
			func(klog.Range) bool { return false },
			func(klog.Duration) bool { return false },
			func(klog.OpenRange) bool { return true },
		)
	})
}

// findLastEntry finds the last entry that matches the predicate, or -1 if none match.
func (r *Reconciler) findLastEntry(match func(klog.Entry) bool) int {
	candidate := -1
	for i, e := range r.Record.Entries() {
		if match(e) {
			candidate = i
		}
	}
	return candidate
}

// concatenateSummary adds summary text to an existing entry that potentially already has one ore
// more lines of summary text.
func (r *Reconciler) concatenateSummary(entryIndex int, entryLineIndex int, additionalSummary klog.EntrySummary) {
	// Append additional summary text. Due to multiline entry summaries, that might
	// not be the same line as the time value.
	lineIndexOfLastSummaryLine := entryLineIndex + countLines([]klog.Entry{r.Record.Entries()[entryIndex]}) - 1
	if len(additionalSummary) > 0 {
		if len(additionalSummary[0]) > 0 {
			// If there is additional summary text, always prepend a space to delimit
			// the additional summary from either the time value or from an already
			// existing summary text.
			r.lines[lineIndexOfLastSummaryLine].Text += " "
		}
		r.lines[lineIndexOfLastSummaryLine].Text += additionalSummary[0]
	}

	if len(additionalSummary) > 1 {
		var subsequentSummaryLines []insertableText
		for _, nextLine := range additionalSummary[1:] {
			subsequentSummaryLines = append(subsequentSummaryLines, insertableText{nextLine, 2})
		}
		r.insert(lineIndexOfLastSummaryLine+1, subsequentSummaryLines)
	}
}

var blankLine = insertableText{"", 0}

type insertableText struct {
	text        string
	indentation int
}

func (r *Reconciler) insert(lineIndex int, texts []insertableText) {
	result := make([]txt.Line, len(r.lines)+len(texts))
	offset := 0
	for i := range result {
		if i >= lineIndex && offset < len(texts) {
			line := strings.Repeat(r.style.indentation.Get(), texts[offset].indentation)
			line += texts[offset].text + r.style.lineEnding.Get()
			result[i] = txt.NewLineFromString(line)
			offset++
		} else {
			result[i] = r.lines[i-offset]
		}
	}
	if lineIndex > 0 && result[lineIndex-1].LineEnding == "" {
		result[lineIndex-1].LineEnding = r.style.lineEnding.Get()
	}
	r.lines = result
}

func toMultilineEntryTexts(entryValue string, entrySummary klog.EntrySummary) []insertableText {
	var result []insertableText
	firstLine := func() string {
		text := entryValue
		// Make sure that there is a space between entry value and the subsequent
		// summary text. However, there shouldn’t be dangling spaces, in case either
		// value would be absent.
		if len(entrySummary) > 0 {
			if len(text) > 0 && len(entrySummary[0]) > 0 {
				text += " "
			}
			text += entrySummary[0]
			entrySummary = entrySummary[1:]
		}
		return text
	}()
	result = append(result, insertableText{firstLine, 1})
	for _, s := range entrySummary {
		result = append(result, insertableText{s, 2})
	}
	return result
}
