package parser

import (
	"errors"
	. "klog"
	"klog/parser/parsing"
)

func (pr *ParseResult) AddEntry(
	errorMessage string,
	matchRecord func(Record) bool,
	createEntry func(Record) string,
) (Record, string, error) {
	index := -1
	for i, r := range pr.Records {
		if matchRecord(r) {
			index = i
			break
		}
	}
	if index == -1 {
		return nil, parsing.Join(pr.lines), errors.New(errorMessage)
	}
	newEntry := createEntry(pr.Records[index])
	lineIndex := pr.lastLineOfRecord[index]
	result := parsing.Insert(pr.lines, lineIndex, newEntry, true, pr.preferences)
	newFileText := parsing.Join(result)
	newRecords, pErr := Parse(newFileText)
	if pErr != nil {
		err := pErr.Get()[0]
		return nil, parsing.Join(pr.lines), errors.New(err.Message())
	}
	return newRecords.Records[index], parsing.Join(result), nil
}
