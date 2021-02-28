package parser

import (
	"errors"
	. "klog"
	"klog/parser/parsing"
)

func (pr *ParseResult) AddEntry(createEntry func([]Record) (int, string)) (string, error) {
	recordIndex, newEntry := createEntry(pr.Records)
	if recordIndex > len(pr.Records)-1 {
		return parsing.Join(pr.lines), errors.New("")
	}
	lineIndex := pr.fileInfo.recordLastLine[recordIndex]
	result := parsing.Insert(pr.lines, lineIndex, pr.fileInfo.indentation+newEntry+pr.fileInfo.lineEnding)
	newFileText := parsing.Join(result)
	_, err := Parse(newFileText)
	if err != nil {
		return parsing.Join(pr.lines), err
	}
	return parsing.Join(result), nil
}
