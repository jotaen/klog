/*
Package parser contains the logic how to convert Record objects from and to plain text.
*/
package parser

import (
	. "github.com/jotaen/klog/src"
	. "github.com/jotaen/klog/src/parser/engine"
	"strings"
)

// ParsedRecord is a record along with some meta information which is
// obtained throughout the parsing process.
type ParsedRecord struct {
	Record

	// Block contains the original lines of text.
	Block Block

	// Style contains the original styling preferences.
	Style *Style
}

// Parse parses a text into a list of Record datastructures. On success, it returns
// the parsed records. Otherwise, it returns all encountered parser errors.
func Parse(recordsAsText string) ([]ParsedRecord, []Error) {
	var results []ParsedRecord
	var allErrs []Error
	blocks := GroupIntoBlocks(Split(recordsAsText))
	for _, block := range blocks {
		record, style, errs := parseRecord(block.SignificantLines())
		if len(errs) > 0 {
			allErrs = append(allErrs, errs...)
			continue
		}
		if block[0].LineEnding != "" {
			style.SetLineEnding(block[0].LineEnding)
		}
		results = append(results, ParsedRecord{
			Record: record,
			Block:  block,
			Style:  style,
		})
	}
	if len(allErrs) > 0 {
		return nil, allErrs
	}
	return results, nil
}

var allowedIndentationStyles = []string{"    ", "   ", "  ", "\t"}

func parseRecord(lines []Line) (Record, *Style, []Error) {
	var errs []Error
	style := DefaultStyle()

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
		rDate, dErr := NewDateFromString(dateText.ToString())
		if dErr != nil {
			errs = append(errs, ErrorInvalidDate().New(headline.Line, headline.PointerPosition, dateText.Length()))
			return nil
		}
		headline.Advance(dateText.Length())
		headline.SkipWhile(IsSpaceOrTab)
		r := NewRecord(rDate)
		style.SetDateFormat(rDate.Format())

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
		style.SetIndentation(indentator.Style())
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
				var entrySummary EntrySummary
				if IsSpaceOrTab(entry.Peek()) {
					entry.Advance(1)
					summaryText, _ := entry.PeekUntil(Is(END_OF_TEXT))
					entrySummary = newEntrySummaryOrNil(summaryText)
				}
				record.AddDuration(duration, entrySummary)
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
			style.SetTimeFormat(start.Format())

			entryStartPositionEnds := entry.PointerPosition
			entry.SkipWhile(Is(' '))
			if entryStartPositionEnds != entry.PointerPosition {
				style.SetSpacingInRange(strings.Repeat(" ", entry.PointerPosition-entryStartPositionEnds))
			} else {
				style.SetSpacingInRange("")
			}

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
				var entrySummary EntrySummary
				if IsSpaceOrTab(entry.Peek()) {
					entry.Advance(1)
					summaryText, _ := entry.PeekUntil(Is(END_OF_TEXT))
					entrySummary = newEntrySummaryOrNil(summaryText)
				}
				sErr := record.StartOpenRange(start, entrySummary)
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
				var entrySummary EntrySummary
				if IsSpaceOrTab(entry.Peek()) {
					entry.Advance(1)
					summaryText, _ := entry.PeekUntil(Is(END_OF_TEXT))
					entrySummary = newEntrySummaryOrNil(summaryText)
				}
				record.AddRange(timeRange, entrySummary)
			}
			return nil
		}()
		if eErr != nil {
			errs = append(errs, eErr)
		}
	}

	if len(errs) > 0 {
		return nil, nil, errs
	}
	return record, style, nil
}

func newEntrySummaryOrNil(singleLineText Parseable) EntrySummary {
	s, err := NewEntrySummary(singleLineText.ToString())
	if err != nil {
		// This can’t happen yet with single-line entry summaries.
		panic("Illegal entry summary")
	}
	return s
}
