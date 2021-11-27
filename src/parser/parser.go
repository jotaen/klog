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

		// There is no leading whitespace allowed in the headline.
		if IsSpaceOrTab(headline.Peek()) {
			errs = append(errs, ErrorIllegalIndentation().New(headline.Line, 0, headline.Length()))
			return nil
		}

		// Parse the date.
		dateText, _ := headline.PeekUntil(IsSpaceOrTab)
		date, dErr := NewDateFromString(dateText.ToString())
		if dErr != nil {
			errs = append(errs, ErrorInvalidDate().New(headline.Line, headline.PointerPosition, dateText.Length()))
			return nil
		}
		headline.Advance(dateText.Length())
		headline.SkipWhile(IsSpaceOrTab)
		r := NewRecord(date)

		// Check if there is a should-total set, and if so, parse it.
		if headline.Peek() == '(' {
			headline.Advance(1) // '('
			headline.SkipWhile(IsSpaceOrTab)
			allPropsText, hasClosingParenthesis := headline.PeekUntil(Is(')'))
			if !hasClosingParenthesis {
				errs = append(errs, ErrorMalformedPropertiesSyntax().New(headline.Line, headline.Length(), 1))
				return r
			}
			if allPropsText.Length() == 0 {
				errs = append(errs, ErrorMalformedPropertiesSyntax().New(headline.Line, headline.PointerPosition, 1))
				return r
			}
			shouldTotalText, hasExclamationMark := headline.PeekUntil(Is('!'))
			if !hasExclamationMark {
				errs = append(errs, ErrorUnrecognisedProperty().New(headline.Line, headline.PointerPosition, shouldTotalText.Length()-1))
				return r
			}
			shouldTotal, sErr := NewDurationFromString(shouldTotalText.ToString())
			if sErr != nil {
				errs = append(errs, ErrorMalformedShouldTotal().New(headline.Line, headline.PointerPosition, shouldTotalText.Length()))
				return r
			}
			r.SetShouldTotal(shouldTotal)
			headline.Advance(shouldTotalText.Length())
			headline.Advance(1) // '!'
			headline.SkipWhile(IsSpaceOrTab)

			// Make sure there is no other text between the braces.
			if headline.Peek() != ')' {
				errs = append(errs, ErrorUnrecognisedProperty().New(headline.Line, headline.PointerPosition, headline.RemainingLength()-1))
				return r
			}
			headline.Advance(1) // ')'
		}

		// Make sure there is no other text left in the headline.
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
		newSummary, sErr := NewRecordSummary(append(record.Summary(), summary.ToString())...)
		if sErr != nil {
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
		eErr := func() Error {
			// Check for correct indentation.
			entry := indentator.NewIndentedParseable(l, 1)
			if entry == nil || IsSpaceOrTab(entry.Peek()) {
				return ErrorIllegalIndentation().New(l, 0, len(l.Text))
			}

			// Try to interpret the entry value as duration.
			durationCandidate, _ := entry.PeekUntil(IsSpaceOrTab)
			duration, dErr := NewDurationFromString(durationCandidate.ToString())
			if dErr == nil {
				entry.Advance(durationCandidate.Length())
				if IsSpaceOrTab(entry.Peek()) {
					entry.Advance(1)
				}
				summaryText, _ := entry.PeekUntil(Is(END_OF_TEXT))
				record.AddDuration(duration, NewEntrySummary(summaryText.ToString()))
				return nil
			}

			// If the entry value isn’t a duration, it must be the start time of a range.
			startCandidate, _ := entry.PeekUntil(Is('-', ' '))
			if startCandidate.Length() == 0 {
				// Handle case where `-` appears right at the beginning of the line.
				firstToken, _ := entry.PeekUntil(IsSpaceOrTab)
				return ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, firstToken.Length())
			}
			start, t1Err := NewTimeFromString(startCandidate.ToString())
			if t1Err != nil {
				return ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, startCandidate.Length())
			}
			entryStartPosition := startCandidate.PointerPosition
			entry.Advance(startCandidate.Length())
			entry.SkipWhile(Is(' '))
			if entry.Peek() != '-' {
				return ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, 1)
			}
			entry.Advance(1) // '-'
			entry.SkipWhile(Is(' '))

			// Check whether the range is a regular or an open-ended one.
			if entry.Peek() == '?' {
				entry.Advance(1)
				placeholder, _ := entry.PeekUntil(IsSpaceOrTab)

				// The placeholder can appear multiple times.
				for _, p := range placeholder.Chars {
					if p != '?' {
						return ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, placeholder.Length())
					}
				}
				entry.Advance(placeholder.Length())
				if IsSpaceOrTab(entry.Peek()) {
					entry.Advance(1)
				}
				summaryText, _ := entry.PeekUntil(Is(END_OF_TEXT))
				sErr := record.StartOpenRange(start, NewEntrySummary(summaryText.ToString()))
				if sErr != nil {
					return ErrorDuplicateOpenRange().New(entry.Line, entryStartPosition, entry.PointerPosition-entryStartPosition)
				}
			} else {
				endCandidate, _ := entry.PeekUntil(IsSpaceOrTab)
				if endCandidate.Length() == 0 {
					return ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, 1)
				}
				end, t2Err := NewTimeFromString(endCandidate.ToString())
				if t2Err != nil {
					return ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, endCandidate.Length())
				}
				entry.Advance(endCandidate.Length())
				timeRange, rErr := NewRange(start, end)
				if rErr != nil {
					return ErrorIllegalRange().New(entry.Line, entryStartPosition, entry.PointerPosition-entryStartPosition)
				}
				if IsSpaceOrTab(entry.Peek()) {
					entry.Advance(1)
				}
				summaryText, _ := entry.PeekUntil(Is(END_OF_TEXT))
				record.AddRange(timeRange, NewEntrySummary(summaryText.ToString()))
			}
			return nil
		}()
		if eErr != nil {
			errs = append(errs, eErr)
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}
	return record, nil
}
