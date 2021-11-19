package parser

import (
	"errors"
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser/parsing"
	"regexp"
)

type ReconcileResult struct {
	NewRecord Record
	NewText   string
}

type RecordReconciler struct {
	pr            *ParseResult
	recordPointer uint // `-1` indicates to prepend
}

type BlockReconciler struct {
	pr                 *ParseResult
	maybeRecordPointer int
}

func NewRecordReconciler(pr *ParseResult, matchRecord func(Record) bool) *RecordReconciler {
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
	result := parsing.Insert(
		r.pr.lines,
		r.pr.lastLineOfRecord[r.recordPointer],
		[]parsing.Text{{newEntry, 1}},
		r.pr.preferences,
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
	openRangeLineIndex := r.pr.lastLineOfRecord[r.recordPointer] - len(record.Entries()) + entryIndex
	originalText := r.pr.lines[openRangeLineIndex].Text
	summaryText := func() string {
		if summary.IsEmpty() {
			return ""
		}
		return " " + summary[0]
	}()
	r.pr.lines[openRangeLineIndex].Text = regexp.MustCompile(`^(.*?)\?+(.*)$`).
		ReplaceAllString(originalText, "${1}"+time.ToString()+"${2}"+summaryText)
	return makeResult(r.pr.lines, r.recordPointer)
}

func NewBlockReconciler(pr *ParseResult, newDate Date) *BlockReconciler {
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

var blankLine = parsing.Text{"", 0}

func (r *BlockReconciler) InsertBlock(texts []parsing.Text) (*ReconcileResult, error) {
	lineIndex, newRecordIndex, insertable := func() (int, uint, []parsing.Text) {
		if r.maybeRecordPointer == -1 {
			if len(r.pr.Records) == 0 {
				return 0, 0, texts
			}
			return 0, 0, append(texts, blankLine)
		}
		return r.pr.lastLineOfRecord[r.maybeRecordPointer],
			uint(r.maybeRecordPointer + 1),
			append([]parsing.Text{blankLine}, texts...)
	}()
	lines := parsing.Insert(
		r.pr.lines,
		lineIndex,
		insertable,
		r.pr.preferences,
	)
	return makeResult(lines, newRecordIndex)
}

func makeResult(ls []parsing.Line, recordIndex uint) (*ReconcileResult, error) {
	newText := parsing.Join(ls)
	newRecords, pErr := Parse(newText)
	if pErr != nil {
		err := pErr.Get()[0]
		return nil, errors.New(err.Message())
	}
	return &ReconcileResult{
		newRecords.Records[recordIndex],
		newText,
	}, nil
}
