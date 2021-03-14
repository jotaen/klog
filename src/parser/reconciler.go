package parser

import (
	"errors"
	. "klog"
	"klog/parser/parsing"
)

type Reconciler struct {
	pr    *ParseResult
	index int
}

func NewReconciler(
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
		pr:    pr,
		index: index,
	}, nil
}

func (r *Reconciler) AppendEntry(handler func(Record) string) (Record, string, error) {
	newEntry := handler(r.pr.Records[r.index])
	lineIndex := r.pr.lastLineOfRecord[r.index]
	result := parsing.Insert(r.pr.lines, lineIndex, newEntry, true, r.pr.preferences)
	newFileText := parsing.Join(result)
	newRecords, pErr := Parse(newFileText)
	if pErr != nil {
		err := pErr.Get()[0]
		return nil, parsing.Join(r.pr.lines), errors.New(err.Message())
	}
	return newRecords.Records[r.index], parsing.Join(result), nil
}
