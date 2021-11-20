/*
Package reconciler contains logic to manipulate klog source text.
The idea of reconcilers in general is to add or modify serialised records
in a minimally invasive manner. Instead or re-serialising the record itself,
it tries to find the location in the original text and modify that directly.
While this approach might feel a little hacky, it avoids lots of other
complications and sources of bugs, that could potentially mess up user data.
*/
package reconciler

import (
	"errors"
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/parsing"
	"regexp"
)

// ReconcileResult is the return value of reconcilers and contains
// the modified record as Record and serialised.
type ReconcileResult struct {
	NewRecord Record
	NewText   string
}

// RecordReconciler is for inserting a new entry into a record.
type RecordReconciler struct {
	pr            *parser.ParseResult
	recordPointer uint // `-1` indicates to prepend
}

// BlockReconciler is for inserting a new record into a list of records.
type BlockReconciler struct {
	pr                 *parser.ParseResult
	maybeRecordPointer int
}

func NewRecordReconciler(pr *parser.ParseResult, matchRecord func(Record) bool) *RecordReconciler {
	index := -1
	for i, r := range pr.Records {
		if matchRecord(r) {
			index = i
			break
		}
	}
	if index == -1 {
		return nil
	}
	return &RecordReconciler{
		pr:            pr,
		recordPointer: uint(index),
	}
}

func (r *RecordReconciler) AppendEntry(handler func(Record) string) (*ReconcileResult, error) {
	newEntry := handler(r.pr.Records[r.recordPointer])
	result := Insert(
		r.pr.Lines,
		r.pr.LastLineOfRecord[r.recordPointer],
		[]Text{{newEntry, 1}},
		r.pr.Preferences,
	)
	return makeResult(result, r.recordPointer)
}

func (r *RecordReconciler) CloseOpenRange(handler func(Record) (Time, EntrySummary)) (*ReconcileResult, error) {
	record := r.pr.Records[r.recordPointer]
	if record.OpenRange() == nil {
		return nil, errors.New("No open time range found")
	}
	entryIndex := 0
	for i, e := range record.Entries() {
		e.Unbox(
			func(Range) interface{} { return nil },
			func(Duration) interface{} { return nil },
			func(OpenRange) interface{} {
				entryIndex = i
				return nil
			},
		)
	}
	time, summary := handler(record)
	openRangeLineIndex := r.pr.LastLineOfRecord[r.recordPointer] - len(record.Entries()) + entryIndex
	originalText := r.pr.Lines[openRangeLineIndex].Text
	summaryText := func() string {
		if summary.IsEmpty() {
			return ""
		}
		return " " + summary[0]
	}()
	r.pr.Lines[openRangeLineIndex].Text = regexp.MustCompile(`^(.*?)\?+(.*)$`).
		ReplaceAllString(originalText, "${1}"+time.ToString()+"${2}"+summaryText)
	return makeResult(r.pr.Lines, r.recordPointer)
}

func NewBlockReconciler(pr *parser.ParseResult, newDate Date) *BlockReconciler {
	index := -1
	for i, r := range pr.Records {
		if i == 0 && !newDate.IsAfterOrEqual(r.Date()) {
			break
		}
		if i == len(pr.Records)-1 {
			index = len(pr.Records) - 1
			break
		}
		if newDate.IsAfterOrEqual(r.Date()) && !newDate.IsAfterOrEqual(pr.Records[i+1].Date()) {
			index = i
			break
		}
	}
	return &BlockReconciler{
		pr:                 pr,
		maybeRecordPointer: index,
	}
}

var blankLine = Text{"", 0}

func (r *BlockReconciler) InsertBlock(texts []Text) (*ReconcileResult, error) {
	lineIndex, newRecordIndex, insertable := func() (int, uint, []Text) {
		if r.maybeRecordPointer == -1 {
			if len(r.pr.Records) == 0 {
				return 0, 0, texts
			}
			return 0, 0, append(texts, blankLine)
		}
		return r.pr.LastLineOfRecord[r.maybeRecordPointer],
			uint(r.maybeRecordPointer + 1),
			append([]Text{blankLine}, texts...)
	}()
	lines := Insert(
		r.pr.Lines,
		lineIndex,
		insertable,
		r.pr.Preferences,
	)
	return makeResult(lines, newRecordIndex)
}

func makeResult(ls []parsing.Line, recordIndex uint) (*ReconcileResult, error) {
	newText := parsing.Join(ls)
	newRecords, pErr := parser.Parse(newText)
	if pErr != nil {
		err := pErr.Get()[0]
		return nil, errors.New(err.Message())
	}
	return &ReconcileResult{
		newRecords.Records[recordIndex],
		newText,
	}, nil
}
