package parser

import (
	. "klog"
	. "klog/parser/engine"
)

// Parse parses a text with records into Record data structures.
func Parse(recordsAsText string) ([]Record, Errors) {
	var records []Record
	var allErrs []Error
	lines := Split(recordsAsText)
	blocks := GroupIntoBlocks(lines)
	for _, block := range blocks {
		r, errs := parseRecord(block)
		if len(errs) > 0 {
			allErrs = append(allErrs, errs...)
		}
		records = append(records, r)
	}
	if len(allErrs) > 0 {
		return nil, NewErrors(allErrs)
	}
	return records, nil
}

func parseRecord(block []Line) (Record, []Error) {
	var errs []Error

	// ========== HEADLINE ==========
	record := func(headline Line) Record {
		if headline.IndentationLevel != 0 {
			errs = append(errs, ErrorIllegalIndentation(NewError(headline, 0, headline.Length())))
			return nil
		}
		dateText, _ := headline.PeekUntil(func(r rune) bool { return IsWhitespace(r) })
		date, err := NewDateFromString(dateText.ToString())
		if err != nil {
			errs = append(errs, ErrorInvalidDate(NewError(headline, headline.PointerPosition, dateText.Length())))
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
				errs = append(errs, ErrorMalformedPropertiesSyntax(NewError(headline, headline.Length(), 1)))
				return r
			}
			if allPropsText.Length() == 0 {
				errs = append(errs, ErrorMalformedPropertiesSyntax(NewError(headline, headline.PointerPosition, 1)))
				return r
			}
			shouldTotalText, hasExclamationMark := headline.PeekUntil(func(r rune) bool { return r == '!' })
			if !hasExclamationMark {
				errs = append(errs, ErrorUnrecognisedProperty(NewError(headline, headline.PointerPosition, shouldTotalText.Length()-1)))
				return r
			}
			shouldTotal, err := NewDurationFromString(shouldTotalText.ToString())
			if err != nil {
				errs = append(errs, ErrorMalformedShouldTotal(NewError(headline, headline.PointerPosition, shouldTotalText.Length())))
				return r
			}
			r.SetShouldTotal(shouldTotal)
			headline.Advance(shouldTotalText.Length())
			headline.Advance(1) // '!'
			headline.SkipWhitespace()
			if headline.Peek() != ')' {
				errs = append(errs, ErrorUnrecognisedProperty(NewError(headline, headline.PointerPosition, headline.RemainingLength()-1)))
				return r
			}
			headline.Advance(1) // ')'
		}
		headline.SkipWhitespace()
		if headline.RemainingLength() > 0 {
			errs = append(errs, ErrorUnrecognisedTextInHeadline(NewError(headline, headline.PointerPosition, headline.RemainingLength())))
		}
		return r
	}(block[0])
	block = block[1:]

	// In case there was an error, generate dummy record to ensure that we have something to
	// work with during parsing. That allows us to continue even if there are errors early on.
	if record == nil {
		dummyDate, _ := NewDate(0, 0, 0)
		record = NewRecord(dummyDate)
	}

	// ========== SUMMARY LINES ==========
	for i, sLine := range block {
		if sLine.IndentationLevel > 0 {
			break
		} else if sLine.IndentationLevel < 0 {
			errs = append(errs, ErrorIllegalIndentation(NewError(sLine, 0, sLine.Length())))
		}
		lineBreak := ""
		if i > 0 {
			lineBreak = "\n"
		}
		err := record.SetSummary(record.Summary().ToString() + lineBreak + sLine.ToString())
		block = block[1:]
		if err != nil {
			errs = append(errs, ErrorMalformedSummary(NewError(sLine, 0, sLine.Length())))
		}
	}

	// ========== ENTRIES ==========
entries:
	for _, eLine := range block {
		if eLine.IndentationLevel != 1 {
			errs = append(errs, ErrorIllegalIndentation(NewError(eLine, 0, eLine.Length())))
			continue
		}
		durationCandidate, _ := eLine.PeekUntil(func(r rune) bool { return IsWhitespace(r) })
		duration, err := NewDurationFromString(durationCandidate.ToString())
		if err == nil {
			eLine.Advance(durationCandidate.Length())
			eLine.SkipWhitespace()
			summaryText, _ := eLine.PeekUntil(func(r rune) bool { return false })
			record.AddDuration(duration, Summary(summaryText.ToString()))
			continue
		}
		startCandidate, _ := eLine.PeekUntil(func(r rune) bool { return r == '-' || IsWhitespace(r) })
		if startCandidate.Length() == 0 {
			// Handle case where `-` appears right at the beginning of the line
			firstToken, _ := eLine.PeekUntil(func(r rune) bool { return IsWhitespace(r) })
			errs = append(errs, ErrorMalformedEntry(NewError(eLine, eLine.PointerPosition, firstToken.Length())))
			continue
		}
		start, err := NewTimeFromString(startCandidate.ToString())
		if err != nil {
			errs = append(errs, ErrorMalformedEntry(NewError(eLine, eLine.PointerPosition, startCandidate.Length())))
			continue
		}
		eLine.Advance(startCandidate.Length())
		eLine.SkipWhitespace()
		if eLine.Peek() != '-' {
			errs = append(errs, ErrorMalformedEntry(NewError(eLine, eLine.PointerPosition, 1)))
			continue
		}
		eLine.Advance(1) // '-'
		eLine.SkipWhitespace()
		if eLine.Peek() == '?' {
			eLine.Advance(1)
			placeholder, _ := eLine.PeekUntil(func(r rune) bool { return IsWhitespace(r) })
			for _, p := range placeholder.Value {
				if p != '?' {
					errs = append(errs, ErrorMalformedEntry(NewError(eLine, eLine.PointerPosition, placeholder.Length())))
					continue entries
				}
			}
			eLine.Advance(placeholder.Length())
			eLine.SkipWhitespace()
			summaryText, _ := eLine.PeekUntil(func(r rune) bool { return false })
			err := record.StartOpenRange(start, Summary(summaryText.ToString()))
			if err != nil {
				errs = append(errs, ErrorDuplicateOpenRange(NewError(eLine, 0, eLine.PointerPosition)))
				continue
			}
		} else {
			endCandidate, _ := eLine.PeekUntil(func(r rune) bool { return IsWhitespace(r) })
			if endCandidate.Length() == 0 {
				errs = append(errs, ErrorMalformedEntry(NewError(eLine, eLine.PointerPosition, 1)))
				continue
			}
			end, err := NewTimeFromString(endCandidate.ToString())
			if err != nil {
				errs = append(errs, ErrorMalformedEntry(NewError(eLine, eLine.PointerPosition, endCandidate.Length())))
				continue
			}
			eLine.Advance(endCandidate.Length())
			timeRange, err := NewRange(start, end)
			if err != nil {
				errs = append(errs, ErrorIllegalRange(NewError(eLine, 0, eLine.PointerPosition)))
				continue
			}
			eLine.SkipWhitespace()
			summaryText, _ := eLine.PeekUntil(func(r rune) bool { return false })
			record.AddRange(timeRange, Summary(summaryText.ToString()))
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}
	return record, nil
}
