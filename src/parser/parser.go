package parser

import (
	. "klog"
	. "klog/parser/parsing"
)

type ParseResult struct {
	Records          []Record
	lines            []Line
	lastLineOfRecord []int
	preferences      Preferences
}

// Parse parses a text with records into Record data structures.
func Parse(recordsAsText string) (*ParseResult, Errors) {
	parseResult := ParseResult{
		Records:          nil,
		lines:            Split(recordsAsText),
		lastLineOfRecord: nil,
		preferences:      DefaultPreferences(),
	}
	var allErrs []Error
	blocks := GroupIntoBlocks(parseResult.lines)
	for _, block := range blocks {
		r, errs := parseRecord(block)
		if len(errs) > 0 {
			allErrs = append(allErrs, errs...)
		}
		parseResult.Records = append(parseResult.Records, r)
		parseResult.lastLineOfRecord = append(
			parseResult.lastLineOfRecord,
			block[len(block)-1].LineNumber,
		)
		for _, l := range block {
			parseResult.preferences.Adapt(&l)
		}
	}
	if len(allErrs) > 0 {
		return nil, NewErrors(allErrs)
	}
	return &parseResult, nil
}

func parseRecord(block []Line) (Record, []Error) {
	var errs []Error

	// ========== HEADLINE ==========
	record := func(headline Parseable) Record {
		if headline.IndentationLevel() != 0 {
			errs = append(errs, ErrorIllegalIndentation(NewError(headline.Line, 0, headline.Length())))
			return nil
		}
		dateText, _ := headline.PeekUntil(func(r rune) bool { return IsWhitespace(r) })
		date, err := NewDateFromString(dateText.ToString())
		if err != nil {
			errs = append(errs, ErrorInvalidDate(NewError(headline.Line, headline.PointerPosition, dateText.Length())))
			return nil
		}
		headline.Advance(dateText.Length())
		headline.SkipWhitespace()
		r := NewRecord(date)
		if headline.Peek() == '(' {
			headline.Advance(1) // '('
			headline.SkipWhitespace()
			allPropsText, hasClosingParenthesis := headline.PeekUntil(func(r rune) bool { return r == ')' })
			if !hasClosingParenthesis {
				errs = append(errs, ErrorMalformedPropertiesSyntax(NewError(headline.Line, headline.Length(), 1)))
				return r
			}
			if allPropsText.Length() == 0 {
				errs = append(errs, ErrorMalformedPropertiesSyntax(NewError(headline.Line, headline.PointerPosition, 1)))
				return r
			}
			shouldTotalText, hasExclamationMark := headline.PeekUntil(func(r rune) bool { return r == '!' })
			if !hasExclamationMark {
				errs = append(errs, ErrorUnrecognisedProperty(NewError(headline.Line, headline.PointerPosition, shouldTotalText.Length()-1)))
				return r
			}
			shouldTotal, err := NewDurationFromString(shouldTotalText.ToString())
			if err != nil {
				errs = append(errs, ErrorMalformedShouldTotal(NewError(headline.Line, headline.PointerPosition, shouldTotalText.Length())))
				return r
			}
			r.SetShouldTotal(shouldTotal)
			headline.Advance(shouldTotalText.Length())
			headline.Advance(1) // '!'
			headline.SkipWhitespace()
			if headline.Peek() != ')' {
				errs = append(errs, ErrorUnrecognisedProperty(NewError(headline.Line, headline.PointerPosition, headline.RemainingLength()-1)))
				return r
			}
			headline.Advance(1) // ')'
		}
		headline.SkipWhitespace()
		if headline.RemainingLength() > 0 {
			errs = append(errs, ErrorUnrecognisedTextInHeadline(NewError(headline.Line, headline.PointerPosition, headline.RemainingLength())))
		}
		return r
	}(NewParseable(block[0]))
	block = block[1:]

	// In case there was an error, generate dummy record to ensure that we have something to
	// work with during parsing. That allows us to continue even if there are errors early on.
	if record == nil {
		dummyDate, _ := NewDate(0, 0, 0)
		record = NewRecord(dummyDate)
	}

	// ========== SUMMARY LINES ==========
	for i, s := range block {
		summary := NewParseable(s)
		if summary.IndentationLevel() > 0 {
			break
		} else if summary.IndentationLevel() < 0 {
			errs = append(errs, ErrorIllegalIndentation(NewError(summary.Line, 0, summary.Length())))
		}
		lineBreak := ""
		if i > 0 {
			lineBreak = "\n"
		}
		err := record.SetSummary(record.Summary().ToString() + lineBreak + summary.ToString())
		block = block[1:]
		if err != nil {
			errs = append(errs, ErrorMalformedSummary(NewError(summary.Line, 0, summary.Length())))
		}
	}

	// ========== ENTRIES ==========
entries:
	for _, e := range block {
		entry := NewParseable(e)
		if entry.IndentationLevel() != 1 {
			errs = append(errs, ErrorIllegalIndentation(NewError(entry.Line, 0, entry.Length())))
			continue
		}
		durationCandidate, _ := entry.PeekUntil(func(r rune) bool { return IsWhitespace(r) })
		duration, err := NewDurationFromString(durationCandidate.ToString())
		if err == nil {
			entry.Advance(durationCandidate.Length())
			entry.SkipWhitespace()
			summaryText, _ := entry.PeekUntil(func(r rune) bool { return false })
			record.AddDuration(duration, Summary(summaryText.ToString()))
			continue
		}
		startCandidate, _ := entry.PeekUntil(func(r rune) bool { return r == '-' || IsWhitespace(r) })
		if startCandidate.Length() == 0 {
			// Handle case where `-` appears right at the beginning of the line
			firstToken, _ := entry.PeekUntil(func(r rune) bool { return IsWhitespace(r) })
			errs = append(errs, ErrorMalformedEntry(NewError(entry.Line, entry.PointerPosition, firstToken.Length())))
			continue
		}
		start, err := NewTimeFromString(startCandidate.ToString())
		if err != nil {
			errs = append(errs, ErrorMalformedEntry(NewError(entry.Line, entry.PointerPosition, startCandidate.Length())))
			continue
		}
		entry.Advance(startCandidate.Length())
		entry.SkipWhitespace()
		if entry.Peek() != '-' {
			errs = append(errs, ErrorMalformedEntry(NewError(entry.Line, entry.PointerPosition, 1)))
			continue
		}
		entry.Advance(1) // '-'
		entry.SkipWhitespace()
		if entry.Peek() == '?' {
			entry.Advance(1)
			placeholder, _ := entry.PeekUntil(func(r rune) bool { return IsWhitespace(r) })
			for _, p := range placeholder.Chars {
				if p != '?' {
					errs = append(errs, ErrorMalformedEntry(NewError(entry.Line, entry.PointerPosition, placeholder.Length())))
					continue entries
				}
			}
			entry.Advance(placeholder.Length())
			entry.SkipWhitespace()
			summaryText, _ := entry.PeekUntil(func(r rune) bool { return false })
			err := record.StartOpenRange(start, Summary(summaryText.ToString()))
			if err != nil {
				errs = append(errs, ErrorDuplicateOpenRange(NewError(entry.Line, 0, entry.PointerPosition)))
				continue
			}
		} else {
			endCandidate, _ := entry.PeekUntil(func(r rune) bool { return IsWhitespace(r) })
			if endCandidate.Length() == 0 {
				errs = append(errs, ErrorMalformedEntry(NewError(entry.Line, entry.PointerPosition, 1)))
				continue
			}
			end, err := NewTimeFromString(endCandidate.ToString())
			if err != nil {
				errs = append(errs, ErrorMalformedEntry(NewError(entry.Line, entry.PointerPosition, endCandidate.Length())))
				continue
			}
			entry.Advance(endCandidate.Length())
			timeRange, err := NewRange(start, end)
			if err != nil {
				errs = append(errs, ErrorIllegalRange(NewError(entry.Line, 0, entry.PointerPosition)))
				continue
			}
			entry.SkipWhitespace()
			summaryText, _ := entry.PeekUntil(func(r rune) bool { return false })
			record.AddRange(timeRange, Summary(summaryText.ToString()))
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}
	return record, nil
}
