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

var allowedIndentationStyles = []string{"    ", "   ", "  ", "\t"}

func parseRecord(lines []Line) (Record, []Error) {
	var errs []Error

	// ========== HEADLINE ==========
	record := func() Record {
		headline := NewParseable(lines[0], 0)
		headline.SkipWhile(IsSpaceOrTab)
		if headline.PointerPosition > 0 {
			errs = append(errs, ErrorIllegalIndentation().New(lines[0], 0, len(lines[0].Text)))
			return nil
		}
		dateText, _ := headline.PeekUntil(IsSpaceOrTab)
		date, err := NewDateFromString(dateText.ToString())
		if err != nil {
			errs = append(errs, ErrorInvalidDate().New(headline.Line, headline.PointerPosition, dateText.Length()))
			return nil
		}
		headline.Advance(dateText.Length())
		headline.SkipWhile(IsSpaceOrTab)
		r := NewRecord(date)
		if headline.Peek() == '(' {
			headline.Advance(1) // '('
			headline.SkipWhile(IsSpaceOrTab)
			allPropsText, hasClosingParenthesis := headline.PeekUntil(func(r rune) bool { return r == ')' })
			if !hasClosingParenthesis {
				errs = append(errs, ErrorMalformedPropertiesSyntax().New(headline.Line, headline.Length(), 1))
				return r
			}
			if allPropsText.Length() == 0 {
				errs = append(errs, ErrorMalformedPropertiesSyntax().New(headline.Line, headline.PointerPosition, 1))
				return r
			}
			shouldTotalText, hasExclamationMark := headline.PeekUntil(func(r rune) bool { return r == '!' })
			if !hasExclamationMark {
				errs = append(errs, ErrorUnrecognisedProperty().New(headline.Line, headline.PointerPosition, shouldTotalText.Length()-1))
				return r
			}
			shouldTotal, err := NewDurationFromString(shouldTotalText.ToString())
			if err != nil {
				errs = append(errs, ErrorMalformedShouldTotal().New(headline.Line, headline.PointerPosition, shouldTotalText.Length()))
				return r
			}
			r.SetShouldTotal(shouldTotal)
			headline.Advance(shouldTotalText.Length())
			headline.Advance(1) // '!'
			headline.SkipWhile(IsSpaceOrTab)
			if headline.Peek() != ')' {
				errs = append(errs, ErrorUnrecognisedProperty().New(headline.Line, headline.PointerPosition, headline.RemainingLength()-1))
				return r
			}
			headline.Advance(1) // ')'
		}
		headline.SkipWhile(IsSpaceOrTab)
		if headline.RemainingLength() > 0 {
			errs = append(errs, ErrorUnrecognisedTextInHeadline().New(headline.Line, headline.PointerPosition, headline.RemainingLength()))
		}
		return r
	}()
	lines = lines[1:]

	if record == nil {
		// In case there was an error, generate dummy record to ensure that we have something to
		// work with during parsing. That allows us to continue even if there are errors early on.
		dummyDate, _ := NewDate(0, 0, 0)
		record = NewRecord(dummyDate)
	}

	var indentator *Indentator

	// ========== SUMMARY LINES ==========
	for _, l := range lines {
		indentator = NewIndentator(allowedIndentationStyles, lines[0])
		if indentator != nil {
			break
		}
		summary := NewParseable(l, 0)
		newSummary, err := NewRecordSummary(append(record.Summary(), summary.ToString())...)
		if err != nil {
			errs = append(errs, ErrorMalformedSummary().New(summary.Line, 0, summary.Length()))
		}
		lines = lines[1:]
		record.SetSummary(newSummary)
	}

	// ========== ENTRIES ==========
	for _, l := range lines {
		if indentator == nil {
			// We should never make it here if the indentation could not be determined.
			panic("Could not detect indentation")
		}
		err := func() Error {
			entry := indentator.NewIndentedParseable(l, 1)
			if entry == nil || IsSpaceOrTab(entry.Peek()) {
				return ErrorIllegalIndentation().New(l, 0, len(l.Text))
			}
			durationCandidate, _ := entry.PeekUntil(IsSpaceOrTab)
			duration, err := NewDurationFromString(durationCandidate.ToString())
			if err == nil {
				entry.Advance(durationCandidate.Length())
				entry.SkipWhile(IsSpaceOrTab)
				summaryText, _ := entry.PeekUntil(func(r rune) bool { return false })
				record.AddDuration(duration, NewEntrySummary(summaryText.ToString()))
				return nil
			}
			startCandidate, _ := entry.PeekUntil(func(r rune) bool { return r == '-' || IsSpace(r) })
			if startCandidate.Length() == 0 {
				// Handle case where `-` appears right at the beginning of the line
				firstToken, _ := entry.PeekUntil(IsSpaceOrTab)
				return ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, firstToken.Length())
			}
			start, err := NewTimeFromString(startCandidate.ToString())
			if err != nil {
				return ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, startCandidate.Length())
			}
			entryStartPosition := startCandidate.PointerPosition
			entry.Advance(startCandidate.Length())
			entry.SkipWhile(IsSpace)
			if entry.Peek() != '-' {
				return ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, 1)
			}
			entry.Advance(1) // '-'
			entry.SkipWhile(IsSpace)
			if entry.Peek() == '?' {
				entry.Advance(1)
				placeholder, _ := entry.PeekUntil(IsSpaceOrTab)
				for _, p := range placeholder.Chars {
					if p != '?' {
						return ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, placeholder.Length())
					}
				}
				entry.Advance(placeholder.Length())
				entry.SkipWhile(IsSpaceOrTab)
				summaryText, _ := entry.PeekUntil(func(r rune) bool { return false })
				err := record.StartOpenRange(start, NewEntrySummary(summaryText.ToString()))
				if err != nil {
					return ErrorDuplicateOpenRange().New(entry.Line, entryStartPosition, entry.PointerPosition-entryStartPosition)
				}
			} else {
				endCandidate, _ := entry.PeekUntil(IsSpaceOrTab)
				if endCandidate.Length() == 0 {
					return ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, 1)
				}
				end, err := NewTimeFromString(endCandidate.ToString())
				if err != nil {
					return ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, endCandidate.Length())
				}
				entry.Advance(endCandidate.Length())
				timeRange, err := NewRange(start, end)
				if err != nil {
					return ErrorIllegalRange().New(entry.Line, entryStartPosition, entry.PointerPosition-entryStartPosition)
				}
				entry.SkipWhile(IsSpaceOrTab)
				summaryText, _ := entry.PeekUntil(func(r rune) bool { return false })
				record.AddRange(timeRange, NewEntrySummary(summaryText.ToString()))
			}
			return nil
		}()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}
	return record, nil
}
