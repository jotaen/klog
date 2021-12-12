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
	"regexp"
)

type Reconciler struct {
	parsedRecords []parser.ParsedRecord
}

// Result contains the result of applied reconcilers.
type Result struct {
	record Record
	text   string
}

func (r *Result) Record() Record {
	return r.record
}

func (r *Result) FileContents() string {
	return r.text
}

// InsertableText is for inserting lines of text into a list of Line’s,
// without needing to know anything about indentation or line ending style.
type InsertableText struct {
	Text        string
	Indentation int
}

func NewReconciler(parsedRecords []parser.ParsedRecord) Reconciler {
	return Reconciler{parsedRecords}
}

// AppendEntry tries to find the matching record and append a new entry to it.
func (r *Reconciler) AppendEntry(matchRecord func(Record) bool, handler func(Record) (string, error)) (*Result, error) {
	recordIndex := findRecordIndex(parser.ToRecords(r.parsedRecords), matchRecord)
	if recordIndex == -1 {
		return nil, NotEligibleError{}
	}
	newEntry, err := handler(r.parsedRecords[recordIndex])
	if err != nil {
		return nil, err
	}
	lastEntry := lastLine(r.parsedRecords[recordIndex].Block.SignificantLines())
	result := insert(
		flatten(parser.ToBlocks(r.parsedRecords)),
		lastEntry.LineNumber,
		[]InsertableText{{newEntry, 1}},
		r.parsedRecords[recordIndex].Style,
	)
	return makeResult(result, uint(recordIndex))
}

// CloseOpenRange tries to find the matching record and closes its open time range.
func (r *Reconciler) CloseOpenRange(matchRecord func(Record) bool, handler func(Record) (Time, EntrySummary)) (*Result, error) {
	recordIndex := findRecordIndex(parser.ToRecords(r.parsedRecords), matchRecord)
	if recordIndex == -1 {
		return nil, NotEligibleError{}
	}
	record := r.parsedRecords[recordIndex]
	openRangeEntryIndex := -1
	for i, e := range record.Entries() {
		e.Unbox(
			func(Range) interface{} { return nil },
			func(Duration) interface{} { return nil },
			func(OpenRange) interface{} {
				openRangeEntryIndex = i
				return nil
			},
		)
	}
	if openRangeEntryIndex == -1 {
		return nil, errors.New("No open time range found")
	}
	time, additionalSummary := handler(record)
	lastEntry := lastLine(r.parsedRecords[recordIndex].Block.SignificantLines())
	openRangeLineIndex := lastEntry.LineNumber - len(record.Entries()) + openRangeEntryIndex
	allLines := flatten(parser.ToBlocks(r.parsedRecords))
	summaryText := func() string {
		if additionalSummary.IsEmpty() {
			return ""
		}
		// If there is additional summary text, always prepend a space to delimit
		// the additional summary from either the time value or from an already
		// existing summary text.
		return " " + additionalSummary[0]
	}()
	allLines[openRangeLineIndex].Text = regexp.MustCompile(`^(.*?)\?+(.*)$`).
		ReplaceAllString(
			allLines[openRangeLineIndex].Text,
			"${1}"+time.ToString()+"${2}"+summaryText,
		)
	return makeResult(allLines, uint(recordIndex))
}

// InsertRecord inserts a new record. For finding the right position, it assumes that
// the existing records are chronologically ordered.
func (r *Reconciler) InsertRecord(newDate Date, texts []InsertableText) (*Result, error) {
	recordIndex := findRecordIndexAfterDate(parser.ToRecords(r.parsedRecords), newDate)
	lineNumber, newRecordIndex, insertable := func() (int, uint, []InsertableText) {
		if recordIndex == -1 {
			if len(r.parsedRecords) == 0 {
				return 0, 0, texts
			}
			// If the new record is dated before the existing ones, prepend it.
			return 0, 0, append(texts, blankLine)
		}
		lastEntry := lastLine(r.parsedRecords[recordIndex].Block.SignificantLines())
		return lastEntry.LineNumber,
			uint(recordIndex + 1),
			append([]InsertableText{blankLine}, texts...)
	}()
	style := parser.DefaultStyle()
	if len(r.parsedRecords) > 0 {
		// If there are records in the file, take over the style preferences
		// from the last one.
		style = r.parsedRecords[len(r.parsedRecords)-1].Style
	}
	lines := insert(
		flatten(parser.ToBlocks(r.parsedRecords)),
		lineNumber,
		insertable,
		style,
	)
	return makeResult(lines, newRecordIndex)
}

// NotEligibleError is to indicate that a reconciler isn’t applicable.
type NotEligibleError struct{}

func (e NotEligibleError) Error() string { return "No eligible record found." }

var blankLine = InsertableText{"", 0}

func findRecordIndex(records []Record, matchRecord func(Record) bool) int {
	index := -1
	for i, r := range records {
		if matchRecord(r) {
			index = i
			break
		}
	}
	if index == -1 {
		return -1
	}
	return index
}

func findRecordIndexAfterDate(records []Record, newDate Date) int {
	index := -1
	for i, r := range records {
		if i == 0 && !newDate.IsAfterOrEqual(r.Date()) {
			break
		}
		if i == len(records)-1 {
			index = len(records) - 1
			break
		}
		if newDate.IsAfterOrEqual(r.Date()) && !newDate.IsAfterOrEqual(records[i+1].Date()) {
			index = i
			break
		}
	}
	return index
}
