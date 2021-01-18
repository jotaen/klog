package parser

import (
	. "klog/parser/engine"
	. "klog/record"
)

func Parse(recordsAsText string) ([]Record, Errors) {
	var records []Record
	var allErrs []Error
	cs := SplitIntoChunksOfLines(recordsAsText)
	for _, c := range cs {
		r, errs := parseRecord(c)
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

func parseRecord(c Chunk) (Record, []Error) {
	var errs []Error

	// ========== HEADLINE ==========
	r := func(headline Text) Record {
		headline.SkipWhitespace()
		if headline.PointerPosition != 0 {
			errs = append(errs, ErrorIllegalWhitespace(NewError(headline, 0, headline.PointerPosition)))
		}
		dateText, _ := headline.PeekUntil(func(r rune) bool { return IsWhitespace(r) })
		date, err := NewDateFromString(dateText.ToString())
		if err != nil {
			errs = append(errs, ErrorMalformedDate(NewError(headline, headline.PointerPosition, dateText.Length())))
			// Generate dummy record to ensure that we have something to work with
			// during parsing. That allows us to continue even if there are errors early on.
			dummyDate, _ := NewDate(0, 0, 0)
			return NewRecord(dummyDate)
		}
		headline.Advance(dateText.Length())
		headline.SkipWhitespace()
		r := NewRecord(date)
		if headline.Peek() == '(' {
			headline.Advance(1) // '('
			headline.SkipWhitespace()
			if _, hasClosingParenthesis := headline.PeekUntil(func(r rune) bool { return r == ')' }); !hasClosingParenthesis {
				errs = append(errs, ErrorMalformedPropertiesSyntax(NewError(headline, headline.Length(), 1)))
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
			headline.Advance(1) // ')'
		}
		headline.SkipWhitespace()
		if headline.RemainingLength() > 0 {
			errs = append(errs, ErrorUnrecognisedTextInHeadline(NewError(headline, headline.PointerPosition, headline.RemainingLength())))
		}
		return r
	}(c[0])
	c = c.Pop()

	// ========== SUMMARY LINES ==========
	for i, sLine := range c {
		if sLine.IndentationLevel > 0 {
			break
		}
		lineBreak := ""
		if i > 0 {
			lineBreak = "\n"
		}
		err := r.SetSummary(r.Summary().ToString() + lineBreak + sLine.ToString())
		c = c.Pop()
		if err != nil {
			errs = append(errs, ErrorMalformedSummary(NewError(sLine, 0, sLine.Length())))
		}
	}

	// ========== ENTRIES ==========
entries:
	for _, eLine := range c {
		if eLine.IndentationLevel != 1 {
			errs = append(errs, ErrorIllegalIndentation(NewError(eLine, 0, eLine.Length()), "entry"))
			continue
		}
		durationCandidate, _ := eLine.PeekUntil(func(r rune) bool { return IsWhitespace(r) })
		duration, err := NewDurationFromString(durationCandidate.ToString())
		if err == nil {
			eLine.Advance(durationCandidate.Length())
			eLine.SkipWhitespace()
			summaryText, _ := eLine.PeekUntil(func(r rune) bool { return false })
			r.AddDuration(duration, Summary(summaryText.ToString()))
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
			err := r.StartOpenRange(start, Summary(summaryText.ToString()))
			if err != nil {
				errs = append(errs, ErrorDuplicateOpenRange(NewError(eLine, 0, eLine.PointerPosition)))
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
			}
			eLine.SkipWhitespace()
			summaryText, _ := eLine.PeekUntil(func(r rune) bool { return false })
			r.AddRange(timeRange, Summary(summaryText.ToString()))
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}
	return r, nil
}
