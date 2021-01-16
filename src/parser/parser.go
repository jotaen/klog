package parser

import (
	"errors"
	. "klog/parser/engine"
	. "klog/record"
	"strings"
)

type Records []Record

func Parse(recordsAsText string) ([]Record, []error) {
	var result Records
	var errs []error
	cs := SplitIntoChunksOfLines(recordsAsText)
	for _, c := range cs {
		err := result.parseRecord(c)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return nil, errs
	}
	return result, nil
}

func (records *Records) parseRecord(c Chunk) error {
	// Date
	headline := c[0]
	dateText := headline.PeekUntil(func(r rune) bool { return r == ' ' })
	date, err := NewDateFromString(dateText.ToString())
	if err != nil {
		return errors.New("UNEXPECTED_CHARACTER")
	}
	headline.Advance(dateText.Length())
	r := NewRecord(date)
	headline.SkipWhitespace()

	// Properties
	if headline.Peek() == '(' {
		headline.Advance(1)
		headline.SkipWhitespace()
		shouldTotalText := headline.PeekUntil(func(r rune) bool { return r == '!' })
		shouldTotal, err := NewDurationFromString(shouldTotalText.ToString())
		if err != nil {
			return errors.New("INVALID_VALUE")
		}
		r.SetShouldTotal(shouldTotal)
		headline.Advance(shouldTotalText.Length() + 1)
		headline.SkipWhitespace()
		if headline.Peek() != ')' {
			return errors.New("UNEXPECTED_CHARACTER")
		}
		headline.Advance(1)
		headline.SkipWhitespace()
	}
	if headline.Peek() != END_OF_TEXT {
		return errors.New("UNEXPECTED_CHARACTER")
	}
	c = c[1:] // Done with headline

	// Summary
	var summaryLines []string
	for _, sLine := range c {
		if sLine.IndentationLevel > 0 {
			break
		}
		summaryLines = append(summaryLines, sLine.ToString())
	}
	err = r.SetSummary(strings.Join(summaryLines, "\n"))
	if err != nil {
		return err
	}
	c = c[len(summaryLines):] // Done with Summary

	// Entries
	for _, eLine := range c {
		durationCandidate := eLine.PeekUntil(func(r rune) bool { return r == ' ' })
		duration, err := NewDurationFromString(durationCandidate.ToString())
		if err == nil {
			eLine.Advance(durationCandidate.Length())
			eLine.SkipWhitespace()
			summaryText := eLine.PeekUntil(func(r rune) bool { return false })
			r.AddDuration(duration, Summary(summaryText.ToString()))
			continue
		}
		startCandidate := eLine.PeekUntil(func(r rune) bool { return r == ' ' || r == '-' })
		start, err := NewTimeFromString(startCandidate.ToString())
		if err != nil {
			return errors.New("INVALID_VALUE")
		}
		eLine.Advance(startCandidate.Length())
		eLine.SkipWhitespace()
		if eLine.Peek() != '-' {
			return errors.New("UNEXPECTED_TOKEN")
		}
		eLine.Advance(1)
		eLine.SkipWhitespace()
		endCandidate := eLine.PeekUntil(func(r rune) bool { return r == ' ' })
		end, _ := NewTimeFromString(endCandidate.ToString())
		if end == nil {
			eLine.SkipWhitespace()
			summaryText := eLine.PeekUntil(func(r rune) bool { return false })
			r.StartOpenRange(start, Summary(summaryText.ToString()))
			continue
		}
		timeRange, err := NewRange(start, end)
		if err != nil {
			return err
		}
		eLine.Advance(endCandidate.Length())
		eLine.SkipWhitespace()
		summaryText := eLine.PeekUntil(func(r rune) bool { return false })
		r.AddRange(timeRange, Summary(summaryText.ToString()))
	}

	*records = append(*records, r)
	return nil
}
