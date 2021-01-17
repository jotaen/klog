package parser

import (
	. "klog/parser/engine"
	. "klog/record"
)

func Parse(recordsAsText string) ([]Record, []error) {
	var records []Record
	var errs []error
	cs := SplitIntoChunksOfLines(recordsAsText)
	for _, c := range cs {
		r, err := parseRecord(c)
		if err != nil {
			errs = append(errs, err)
		}
		records = append(records, r)
	}
	if len(errs) > 0 {
		return nil, errs
	}
	return records, nil
}

func parseRecord(c Chunk) (Record, error) {
	// Date
	headline := c[0]
	headline.SkipWhitespace()
	if headline.PointerPosition != 0 {
		return nil, ErrorIllegalWhitespace(NewError(headline, 0, headline.PointerPosition))
	}
	dateText := headline.PeekUntil(func(r rune) bool { return IsWhitespace(r) })
	date, err := NewDateFromString(dateText.ToString())
	if err != nil {
		return nil, ErrorMalformedDate(NewError(headline, headline.PointerPosition, dateText.Length()))
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
			return nil, ErrorMalformedShouldTotal(NewError(headline, headline.PointerPosition, shouldTotalText.Length()))
		}
		r.SetShouldTotal(shouldTotal)
		headline.Advance(shouldTotalText.Length() + 1)
		headline.SkipWhitespace()
		if headline.Peek() != ')' {
			return nil, ErrorMalformedShouldTotal(NewError(headline, headline.PointerPosition, 1))
		}
		headline.Advance(1)
		headline.SkipWhitespace()
	}
	if headline.Peek() != END_OF_TEXT {
		return nil, ErrorMalformedShouldTotal(NewError(headline, headline.PointerPosition, headline.RemainingLength()))
	}
	c = c[1:] // Done with headline

	// Summary
	for i, sLine := range c {
		if sLine.IndentationLevel > 0 {
			break
		}
		lineBreak := ""
		if i > 0 {
			lineBreak = "\n"
		}
		err := r.SetSummary(r.Summary().ToString() + lineBreak + sLine.ToString())
		c = c[1:]
		if err != nil {
			return nil, ErrorMalformedSummary(NewError(sLine, 0, sLine.Length()))
		}
	}

	// Entries
	for _, eLine := range c {
		if eLine.IndentationLevel != 1 {
			return nil, ErrorIllegalIndentation(NewError(eLine, 0, eLine.Length()))
		}
		durationCandidate := eLine.PeekUntil(func(r rune) bool { return IsWhitespace(r) })
		duration, err := NewDurationFromString(durationCandidate.ToString())
		if err == nil {
			eLine.Advance(durationCandidate.Length())
			eLine.SkipWhitespace()
			summaryText := eLine.PeekUntil(func(r rune) bool { return false })
			r.AddDuration(duration, Summary(summaryText.ToString()))
			continue
		}
		startCandidate := eLine.PeekUntil(func(r rune) bool { return r == '-' || IsWhitespace(r) })
		start, err := NewTimeFromString(startCandidate.ToString())
		if err != nil {
			return nil, ErrorMalformedEntry(NewError(eLine, eLine.PointerPosition, startCandidate.Length()))
		}
		eLine.Advance(startCandidate.Length())
		eLine.SkipWhitespace()
		if eLine.Peek() != '-' {
			return nil, ErrorMalformedEntry(NewError(eLine, eLine.PointerPosition, startCandidate.Length()))
		}
		eLine.Advance(1)
		eLine.SkipWhitespace()
		endCandidate := eLine.PeekUntil(func(r rune) bool { return IsWhitespace(r) })
		end, err := NewTimeFromString(endCandidate.ToString())
		if err != nil { // if there’s an “error” here we assume this entry to be an open range
			eLine.SkipWhitespace()
			summaryText := eLine.PeekUntil(func(r rune) bool { return false })
			err := r.StartOpenRange(start, Summary(summaryText.ToString()))
			if err != nil {
				return nil, ErrorDuplicateOpenRange(NewError(eLine, 0, eLine.PointerPosition))
			}
			continue
		}
		timeRange, err := NewRange(start, end)
		eLine.Advance(endCandidate.Length())
		if err != nil {
			return nil, ErrorIllegalRange(NewError(eLine, 0, eLine.PointerPosition))
		}
		eLine.SkipWhitespace()
		summaryText := eLine.PeekUntil(func(r rune) bool { return false })
		r.AddRange(timeRange, Summary(summaryText.ToString()))
	}

	return r, nil
}
