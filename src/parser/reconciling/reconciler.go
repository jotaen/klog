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
	"github.com/jotaen/klog/src/parser/engine"
	"regexp"
)

type Reconciler struct {
	records []Record
	blocks  []engine.Block
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

func NewReconciler(records []Record, blocks []engine.Block) Reconciler {
	return Reconciler{records, blocks}
}

// AppendEntry tries to find the matching record and append a new entry to it.
func (r *Reconciler) AppendEntry(matchRecord func(Record) bool, handler func(Record) string) (*Result, error) {
	recordIndex := findRecordIndex(r.records, matchRecord)
	if recordIndex == -1 {
		return nil, NotEligibleError{}
	}
	newEntry := handler(r.records[recordIndex])
	lastEntry := lastLine(r.blocks[recordIndex].SignificantLines())
	result := insert(
		flatten(r.blocks),
		lastEntry.LineNumber,
		[]InsertableText{{newEntry, 1}},
		stylePreferencesOrDefault(r.blocks[recordIndex]),
	)
	return makeResult(result, uint(recordIndex))
}

// CloseOpenRange tries to find the matching record and closes its open time range.
func (r *Reconciler) CloseOpenRange(matchRecord func(Record) bool, handler func(Record) (Time, EntrySummary)) (*Result, error) {
	recordIndex := findRecordIndex(r.records, matchRecord)
	if recordIndex == -1 {
		return nil, NotEligibleError{}
	}
	record := r.records[recordIndex]
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
	time, summary := handler(record)
	lastEntry := lastLine(r.blocks[recordIndex].SignificantLines())
	openRangeLineIndex := lastEntry.LineNumber - len(record.Entries()) + openRangeEntryIndex
	allLines := flatten(r.blocks)
	originalText := allLines[openRangeLineIndex].Text
	summaryText := func() string {
		if summary.IsEmpty() {
			return ""
		}
		return " " + summary[0]
	}()
	allLines[openRangeLineIndex].Text = regexp.MustCompile(`^(.*?)\?+(.*)$`).
		ReplaceAllString(originalText, "${1}"+time.ToString()+"${2}"+summaryText)
	return makeResult(allLines, uint(recordIndex))
}

// InsertRecord inserts a new record. For finding the right position, it assumes that
// the existing records are chronologically ordered.
func (r *Reconciler) InsertRecord(newDate Date, texts []InsertableText) (*Result, error) {
	recordIndex := findRecordIndexAfterDate(r.records, newDate)
	lineNumber, newRecordIndex, insertable := func() (int, uint, []InsertableText) {
		if recordIndex == -1 {
			if len(r.records) == 0 {
				return 0, 0, texts
			}
			return 0, 0, append(texts, blankLine)
		}
		lastEntry := lastLine(r.blocks[recordIndex].SignificantLines())
		return lastEntry.LineNumber,
			uint(recordIndex + 1),
			append([]InsertableText{blankLine}, texts...)
	}()
	var styleReferenceBlock engine.Block
	if len(r.blocks) > 0 {
		styleReferenceBlock = r.blocks[0]
	}
	lines := insert(
		flatten(r.blocks),
		lineNumber,
		insertable,
		stylePreferencesOrDefault(styleReferenceBlock),
	)
	return makeResult(lines, newRecordIndex)
}

// NotEligibleError is to indicate that a reconciler isn’t applicable.
type NotEligibleError struct{}

func (e NotEligibleError) Error() string { return "Boom" } // TODO

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
