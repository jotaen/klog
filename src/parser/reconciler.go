package parser

import (
	"errors"
	. "klog"
	"klog/parser/parsing"
)

func (pr *ParseResult) AddEntry(matchRecord func(Record) bool, createEntry func(Record) string) (string, error) {
	index := -1
	for i, r := range pr.Records {
		if matchRecord(r) {
			index = i
			break
		}
	}
	if index == -1 {
		return parsing.Join(pr.lines), errors.New("No such record")
	}
	newEntry := createEntry(pr.Records[index])
	lineIndex := pr.lastLineOfRecord[index]
	result := parsing.Insert(pr.lines, lineIndex, newEntry, true, pr.preferences)
	newFileText := parsing.Join(result)
	_, pErr := Parse(newFileText)
	if pErr != nil {
		err := pErr.Get()[0]
		return parsing.Join(pr.lines), errors.New(err.Message())
	}
	return parsing.Join(result), nil
}
