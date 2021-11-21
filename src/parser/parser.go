/*
Package parser contains the logic how to convert Record objects from and to plain text.
*/
package parser

import (
	. "github.com/jotaen/klog/src"
	. "github.com/jotaen/klog/src/parser/engine"
)

// Parse parses a text into a list of Record datastructures.
func Parse(recordsAsText string) ([]Record, []Block, Errors) {
	var records []Record
	var allErrs []Error
	blocks := GroupIntoBlocks(Split(recordsAsText))
	for _, block := range blocks {
		record, errs := parseRecord(block.SignificantLines())
		if len(errs) > 0 {
			allErrs = append(allErrs, errs...)
			continue
		}
		records = append(records, record)
	}
	if len(allErrs) > 0 {
		return nil, nil, NewErrors(allErrs)
	}
	return records, blocks, nil
}

func parseRecord(lines []Line) (Record, []Error) {
	var errs []Error

	// ========== HEADLINE ==========
	record := func(headline Parseable) Record {
		if len(headline.PrecedingWhitespace) > 0 {
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
	}(NewParseable(lines[0]))
	lines = lines[1:]

	if record == nil {
		// In case there was an error, generate dummy record to ensure that we have something to
		// work with during parsing. That allows us to continue even if there are errors early on.
		dummyDate, _ := NewDate(0, 0, 0)
		record = NewRecord(dummyDate)
	}
	indentationDetector := newIndentationDetector()

	// ========== SUMMARY LINES ==========
	for _, s := range lines {
		summary := NewParseable(s)
		isIndented, iErr := indentationDetector.isIndented(summary)
		if iErr != nil {
			errs = append(errs, iErr)
		}
		if isIndented {
			break
		}
		newSummary, err := NewRecordSummary(append(record.Summary(), summary.ToString())...)
		if err != nil {
			errs = append(errs, ErrorMalformedSummary(NewError(summary.Line, 0, summary.Length())))
		}
		lines = lines[1:]
		record.SetSummary(newSummary)
	}

	// ========== ENTRIES ==========
entries:
	for _, e := range lines {
		entry := NewParseable(e)
		iErr := indentationDetector.collate(entry)
		if iErr != nil {
			errs = append(errs, iErr)
			continue
		}
		durationCandidate, _ := entry.PeekUntil(func(r rune) bool { return IsWhitespace(r) })
		duration, err := NewDurationFromString(durationCandidate.ToString())
		if err == nil {
			entry.Advance(durationCandidate.Length())
			entry.SkipWhitespace()
			summaryText, _ := entry.PeekUntil(func(r rune) bool { return false })
			record.AddDuration(duration, NewEntrySummary(summaryText.ToString()))
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
			err := record.StartOpenRange(start, NewEntrySummary(summaryText.ToString()))
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
			record.AddRange(timeRange, NewEntrySummary(summaryText.ToString()))
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}
	return record, nil
}

type indentationDetector struct {
	allowedIndentationStyles []string
	actualIndentationStyle   string
}

func newIndentationDetector() *indentationDetector {
	return &indentationDetector{
		allowedIndentationStyles: []string{"  ", "   ", "    ", "\t"},
	}
}

// collate requires a line to be indented and enforces the indentation style to
// be consistent across subsequent calls.
func (i *indentationDetector) collate(p Parseable) Error {
	if i.actualIndentationStyle == "" {
		i.actualIndentationStyle = p.PrecedingWhitespace
	} else if p.PrecedingWhitespace != i.actualIndentationStyle {
		return ErrorIllegalIndentation(NewError(p.Line, 0, p.Length()))
	}
	return nil
}

// isIndented checks whether a line is indented according to one of the allowed styles.
func (i *indentationDetector) isIndented(l Parseable) (bool, Error) {
	if len(l.PrecedingWhitespace) == 0 {
		return false, nil
	}
	for _, s := range i.allowedIndentationStyles {
		if l.PrecedingWhitespace == s {
			return true, nil
		}
	}
	return false, ErrorIllegalIndentation(NewError(l.Line, 0, l.Length()))
}
