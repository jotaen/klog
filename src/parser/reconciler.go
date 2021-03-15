package parser

import (
	"errors"
	. "klog"
	"klog/parser/parsing"
	"regexp"
)

type ReconcileResult struct {
	NewRecord Record
	NewText   string
}

type RecordReconciler struct {
	pr            *ParseResult
	recordPointer uint // can be `-1` for nothing found
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

func (r *RecordReconciler) CloseOpenRange(handler func(Record) Time) (*ReconcileResult, error) {
	record := r.pr.Records[r.recordPointer]
	if record.OpenRange() == nil {
		return nil, errors.New("NO_OPEN_RANGE")
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
	time := handler(record)
	openRangeLineIndex := r.pr.lastLineOfRecord[r.recordPointer] - len(record.Entries()) + entryIndex
	originalText := r.pr.lines[openRangeLineIndex].Text
	r.pr.lines[openRangeLineIndex].Text = regexp.MustCompile(`^(.*?)\?+(.*)$`).
		ReplaceAllString(originalText, "${1}"+time.ToString()+"${2}")
	return makeResult(r.pr.lines, r.recordPointer)
}

func NewBlockReconciler(pr *ParseResult, findPosition func(Record, Record) bool) *BlockReconciler {
	index := len(pr.Records) - 1
	for i, r := range pr.Records {
		if i == index {
			break
		}
		if findPosition(r, pr.Records[i+1]) {
			index = i
			break
		}
	}
	return &BlockReconciler{
		pr:                 pr,
		maybeRecordPointer: index,
	}
}

func (r *BlockReconciler) AddNewRecord(texts []parsing.Text) (*ReconcileResult, error) {
	lineIndex, newRecordIndex, appendable := func() (int, uint, []parsing.Text) {
		if r.maybeRecordPointer == -1 {
			return 0, 0, texts
		}
		return r.pr.lastLineOfRecord[r.maybeRecordPointer],
			uint(r.maybeRecordPointer + 1),
			append([]parsing.Text{{"", 0}}, texts...)
	}()
	lines := parsing.Insert(
		r.pr.lines,
		lineIndex,
		appendable,
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
