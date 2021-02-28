package parser

import (
	"errors"
	. "klog"
	"klog/parser/engine"
)

func (pr *ParseResult) AddEntry(createEntry func([]Record) (int, string)) (string, error) {
	recordIndex, newEntry := createEntry(pr.Records)
	if recordIndex > len(pr.Records)-1 {
		return engine.Join(pr.lines), errors.New("")
	}
	lineIndex := pr.fileInfo.recordLastLine[recordIndex]
	result := append(pr.lines, engine.Line{})
	copy(result[lineIndex+1:], result[lineIndex:])
	result[lineIndex] = engine.Line{
		Original: pr.fileInfo.indentation + newEntry + pr.fileInfo.lineEnding,
	}
	newFileText := engine.Join(result)
	_, err := Parse(newFileText)
	if err != nil {
		return engine.Join(pr.lines), err
	}
	return engine.Join(result), nil
}
