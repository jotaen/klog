package parser

import (
	"errors"
	. "klog"
	"klog/parser/parsing"
	"regexp"
)

type Reconciler struct {
	pr            *ParseResult
	recordPointer int
}

func NewRecordReconciler(
	pr *ParseResult,
	notFoundError error,
	matchRecord func(Record) bool,
) (*Reconciler, error) {
	index := -1
	for i, r := range pr.Records {
		if matchRecord(r) {
			index = i
			break
		}
	}
	if index == -1 {
		return nil, notFoundError
	}
	return &Reconciler{
		pr:            pr,
		recordPointer: index,
	}, nil
}

func NewBlockReconciler(
	pr *ParseResult,
	findPosition func(Record, Record) bool,
) (*Reconciler, error) {
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
	return &Reconciler{
		pr:            pr,
		recordPointer: index,
	}, nil
}

func (r *Reconciler) AppendEntry(handler func(Record) string) (Record, string, error) {
	newEntry := handler(r.pr.Records[r.recordPointer])
	result := parsing.Insert(
		r.pr.lines,
		r.pr.lastLineOfRecord[r.recordPointer],
		[]parsing.Text{{newEntry, 1}},
		r.pr.preferences,
	)
	return makeResult(result, r.recordPointer)
}

func (r *Reconciler) CloseOpenRange(handler func(Record) Time) (Record, string, error) {
	record := r.pr.Records[r.recordPointer]
	if record.OpenRange() == nil {
		return nil, "", errors.New("NO_OPEN_RANGE")
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

func (r *Reconciler) AddBlock(texts []parsing.Text) (Record, string, error) {
	lineIndex, appendable := func() (int, []parsing.Text) {
		if r.recordPointer == -1 {
			return 0, texts
		}
		return r.pr.lastLineOfRecord[r.recordPointer], append([]parsing.Text{{"", 0}}, texts...)
	}()
	lines := parsing.Insert(
		r.pr.lines,
		lineIndex,
		appendable,
		r.pr.preferences,
	)
	return makeResult(lines, r.recordPointer+1)
}

func makeResult(ls []parsing.Line, recordIndex int) (Record, string, error) {
	newText := parsing.Join(ls)
	newRecords, pErr := Parse(newText)
	if pErr != nil {
		err := pErr.Get()[0]
		return nil, "", errors.New(err.Message())
	}
	return newRecords.Records[recordIndex], newText, nil
}
