package parser

import (
	"errors"
	. "klog"
	"klog/parser/parsing"
)

func (pr *ParseResult) AddEntry(createEntry func([]Record) (int, string)) (string, error) {
	recordIndex, newEntry := createEntry(pr.Records)
	if recordIndex > len(pr.Records)-1 || recordIndex < 0 {
		return parsing.Join(pr.lines), errors.New("No such record")
	}
	lineIndex := pr.lastLineOfRecord[recordIndex]
	result := parsing.Insert(pr.lines, lineIndex, newEntry, true, pr.preferences)
	newFileText := parsing.Join(result)
	_, pErr := Parse(newFileText)
	if pErr != nil {
		err := pErr.Get()[0]
		return parsing.Join(pr.lines), errors.New(err.Message())
	}
	return parsing.Join(result), nil
}
