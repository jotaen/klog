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
	for len(lines) > 0 {
		l := lines[0]
		if indentator == nil {
			// We should never make it here if the indentation could not be determined.
			panic("Could not detect indentation")
		}
		style.SetIndentation(indentator.Style())
		// Check for correct indentation.
		entry := indentator.NewIndentedParseable(l, 1)
		if entry == nil || IsSpaceOrTab(entry.Peek()) {
			errs = append(errs, ErrorIllegalIndentation().New(l, 0, len(l.Text)))
			break
		}
		lines = lines[1:]

		// Parse entry value.
		createEntry, evErr := func() (func(EntrySummary) Error, Error) {
			// Try to interpret the entry value as duration.
			durationCandidate, _ := entry.PeekUntil(IsSpaceOrTab)
			duration, dErr := NewDurationFromString(durationCandidate.ToString())
			if dErr == nil {
				entry.Advance(durationCandidate.Length())
				return func(s EntrySummary) Error {
					record.AddDuration(duration, s)
					return nil
				}, nil
			}

			// If the entry value isn’t a duration, it must be the start time of a range.
			startCandidate, _ := entry.PeekUntil(Is('-', ' '))
			if startCandidate.Length() == 0 {
				// Handle case where `-` appears right at the beginning of the line.
				firstToken, _ := entry.PeekUntil(IsSpaceOrTab)
				return nil, ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, firstToken.Length())
			}
			start, t1Err := NewTimeFromString(startCandidate.ToString())
			if t1Err != nil {
				return nil, ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, startCandidate.Length())
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
				return nil, ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, 1)
			}
			entry.Advance(1) // '-'
			entry.SkipWhile(Is(' '))

			// Check whether the range is open-ended.
			if entry.Peek() == '?' {
				entry.Advance(1)
				placeholder, _ := entry.PeekUntil(IsSpaceOrTab)

				// The placeholder can appear multiple times.
				for _, p := range placeholder.Chars {
					if p != '?' {
						return nil, ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, placeholder.Length())
					}
				}
				entry.Advance(placeholder.Length())
				return func(s EntrySummary) Error {
					sErr := record.StartOpenRange(start, s)
					if sErr != nil {
						return ErrorDuplicateOpenRange().New(entry.Line, entryStartPosition, entry.PointerPosition-entryStartPosition)
					}
					return nil
				}, nil
			}

			// Ultimately, the entry can only be a regular range.
			endCandidate, _ := entry.PeekUntil(IsSpaceOrTab)
			if endCandidate.Length() == 0 {
				return nil, ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, 1)
			}
			end, t2Err := NewTimeFromString(endCandidate.ToString())
			if t2Err != nil {
				return nil, ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, endCandidate.Length())
			}
			entry.Advance(endCandidate.Length())
			timeRange, rErr := NewRange(start, end)
			if rErr != nil {
				return nil, ErrorIllegalRange().New(entry.Line, entryStartPosition, entry.PointerPosition-entryStartPosition)
			}
			return func(s EntrySummary) Error {
				record.AddRange(timeRange, s)
				return nil
			}, nil
		}()

		// Check for error while parsing the entry value.
		if evErr != nil {
			errs = append(errs, evErr)
			continue
		}

		// Parse entry summary.
		entrySummary, esErr := func() (EntrySummary, Error) {
			var result EntrySummary

			// Parse first line of entry summary.
			if IsSpaceOrTab(entry.Peek()) {
				entry.Advance(1)
				summaryText, _ := entry.PeekUntil(Is(END_OF_TEXT))
				firstLine, sErr := NewEntrySummary(summaryText.ToString())
				if sErr != nil {
					return nil, ErrorMalformedSummary().New(summaryText.Line, 0, summaryText.Length())
				}
				result = firstLine
			} else {
				result, _ = NewEntrySummary("")
			}

			// Parse subsequent lines of multiline entry summary.
			for len(lines) > 0 {
				nextEntrySummaryLine := indentator.NewIndentedParseable(lines[0], 2)
				if nextEntrySummaryLine == nil {
					break
				}
				lines = lines[1:]
				additionalText, _ := nextEntrySummaryLine.PeekUntil(func(_ rune) bool {
					return false // Move forward until end of line
				})
				newEntrySummary, sErr := NewEntrySummary(append(result, additionalText.ToString())...)
				if sErr != nil {
					return nil, ErrorMalformedSummary().New(nextEntrySummaryLine.Line, 0, nextEntrySummaryLine.Length())
				}
				result = newEntrySummary
			}

			return result, nil
		}()

		// Check for error while parsing the entry summary.
		if esErr != nil {
			errs = append(errs, esErr)
			continue
		}

		// Check for error when eventually applying the entry.
		eErr := createEntry(entrySummary)
		if eErr != nil {
			errs = append(errs, eErr)
		}
	}

	if len(errs) > 0 {
		return nil, nil, errs
	}
	return record, style, nil
}
