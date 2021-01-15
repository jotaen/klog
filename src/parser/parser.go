package parser

import (
	"errors"
	"klog/parser/engine"
	. "klog/record"
	"strings"
)

type Result []Record

func Parse(recordsAsText string) ([]Record, error) {
	var result Result
	err := engine.Parse(recordsAsText, result.parseRecord)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (result *Result) parseRecord(c engine.Chunk) error {
	// Date
	headline := c[0]
	date, err := NewDateFromString(headline.Peek(10))
	if err != nil {
		return errors.New("UNEXPECTED_CHARACTER")
	}
	headline.Advance(10)
	r := NewRecord(date)
	headline.SkipWhitespace()

	// Properties
	if headline.Peek(1) == "(" {
		headline.Advance(1)
		text, err := headline.PeekUntil(')')
		if err != nil {
			return err
		}
		// TODO process properties
		headline.Advance(len(text) + 1)
		headline.SkipWhitespace()
	}
	if headline.Peek(1) != "" {
		return errors.New("UNEXPECTED_CHARACTER")
	}

	// Summary
	var summaryLines []string
	for _, s := range c[1:] {
		summaryLines = append(summaryLines, string(s.Text))
	}
	r.SetSummary(strings.Join(summaryLines, "\n"))

	*result = append(*result, r)
	return nil
}
