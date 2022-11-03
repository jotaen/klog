/*
Package parser contains the logic how to convert Record objects from and to plain text.
*/
package parser

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/parser/engine"
	"strings"
)

// ParsedRecord is a record along with some meta information which is
// obtained throughout the parsing process.
type ParsedRecord struct {
	klog.Record

	// Block contains the original lines of text.
	Block engine.Block

	// Style contains the original styling preferences.
	Style *Style
}

// Parse parses a text into a list of Record datastructures. On success, it returns
// the parsed records. Otherwise, it returns all encountered parser errors.
func Parse(recordsAsText string) ([]ParsedRecord, []engine.Error) {
	blocks := engine.GroupIntoBlocks(engine.Split(recordsAsText))
	records := make([]ParsedRecord, len(blocks))
	var allErrs []engine.Error
	for i, block := range blocks {
		record, style, errs := parseRecord(block.SignificantLines())
		if len(errs) > 0 {
			allErrs = append(allErrs, errs...)
			continue
		}
		if block[0].LineEnding != "" {
			style.LineEnding.Set(block[0].LineEnding)
		}
		records[i] = ParsedRecord{
			Record: record,
			Block:  block,
			Style:  style,
		}
	}
	if len(allErrs) > 0 {
		return nil, allErrs
	}
	return records, nil
}

var allowedIndentationStyles = []string{"    ", "   ", "  ", "\t"}

func parseRecord(lines []engine.Line) (klog.Record, *Style, []engine.Error) {
	var errs []engine.Error
	style := DefaultStyle()

	// ========== HEADLINE ==========
	record := func() klog.Record {
		headline := engine.NewParseable(lines[0], 0)

		// There is no leading whitespace allowed in the headline.
		if engine.IsSpaceOrTab(headline.Peek()) {
			errs = append(errs, ErrorIllegalIndentation().New(headline.Line, 0, headline.Length()))
			return nil
		}

		// Parse the date.
		dateText, _ := headline.PeekUntil(engine.IsSpaceOrTab)
		rDate, dErr := klog.NewDateFromString(dateText.ToString())
		if dErr != nil {
			errs = append(errs, ErrorInvalidDate().New(headline.Line, headline.PointerPosition, dateText.Length()))
			return nil
		}
		headline.Advance(dateText.Length())
		headline.SkipWhile(engine.IsSpaceOrTab)
		r := klog.NewRecord(rDate)
		style.DateFormat.Set(rDate.Format())

		// Check if there is a should-total set, and if so, parse it.
		if headline.Peek() == '(' {
			headline.Advance(1) // '('
			headline.SkipWhile(engine.IsSpaceOrTab)
			allPropsText, hasClosingParenthesis := headline.PeekUntil(engine.Is(')'))
			if !hasClosingParenthesis {
				errs = append(errs, ErrorMalformedPropertiesSyntax().New(headline.Line, headline.Length(), 1))
				return r
			}
			if allPropsText.Length() == 0 {
				errs = append(errs, ErrorMalformedPropertiesSyntax().New(headline.Line, headline.PointerPosition, 1))
				return r
			}
			shouldTotalText, hasExclamationMark := headline.PeekUntil(engine.Is('!'))
			if !hasExclamationMark {
				errs = append(errs, ErrorUnrecognisedProperty().New(headline.Line, headline.PointerPosition, shouldTotalText.Length()-1))
				return r
			}
			shouldTotal, sErr := klog.NewDurationFromString(shouldTotalText.ToString())
			if sErr != nil {
				errs = append(errs, ErrorMalformedShouldTotal().New(headline.Line, headline.PointerPosition, shouldTotalText.Length()))
				return r
			}
			r.SetShouldTotal(shouldTotal)
			headline.Advance(shouldTotalText.Length())
			headline.Advance(1) // '!'
			headline.SkipWhile(engine.IsSpaceOrTab)

			// Make sure there is no other text between the braces.
			if headline.Peek() != ')' {
				errs = append(errs, ErrorUnrecognisedProperty().New(headline.Line, headline.PointerPosition, headline.RemainingLength()-1))
				return r
			}
			headline.Advance(1) // ')'
		}

		// Make sure there is no other text left in the headline.
		headline.SkipWhile(engine.IsSpaceOrTab)
		if headline.RemainingLength() > 0 {
			errs = append(errs, ErrorUnrecognisedTextInHeadline().New(headline.Line, headline.PointerPosition, headline.RemainingLength()))
		}
		return r
	}()
	lines = lines[1:]

	if record == nil {
		// In case there was an error, generate dummy record to ensure that we have something to
		// work with during parsing. That allows us to continue even if there are errors early on.
		dummyDate, _ := klog.NewDate(0, 0, 0)
		record = klog.NewRecord(dummyDate)
	}

	var indentator *engine.Indentator

	// ========== SUMMARY LINES ==========
	for _, l := range lines {
		indentator = engine.NewIndentator(allowedIndentationStyles, lines[0])
		if indentator != nil {
			break
		}
		summary := engine.NewParseable(l, 0)
		newSummary, sErr := klog.NewRecordSummary(append(record.Summary(), summary.ToString())...)
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
		style.Indentation.Set(indentator.Style())
		// Check for correct indentation.
		entry := indentator.NewIndentedParseable(l, 1)
		if entry == nil || engine.IsSpaceOrTab(entry.Peek()) {
			errs = append(errs, ErrorIllegalIndentation().New(l, 0, len(l.Text)))
			break
		}
		lines = lines[1:]

		// Parse entry value.
		createEntry, evErr := func() (func(klog.EntrySummary) engine.Error, engine.Error) {
			// Try to interpret the entry value as duration.
			durationCandidate, _ := entry.PeekUntil(engine.IsSpaceOrTab)
			duration, dErr := klog.NewDurationFromString(durationCandidate.ToString())
			if dErr == nil {
				entry.Advance(durationCandidate.Length())
				return func(s klog.EntrySummary) engine.Error {
					record.AddDuration(duration, s)
					return nil
				}, nil
			}

			// If the entry value isn’t a duration, it must be the start time of a range.
			startCandidate, _ := entry.PeekUntil(engine.Is('-', ' '))
			if startCandidate.Length() == 0 {
				// Handle case where `-` appears right at the beginning of the line.
				firstToken, _ := entry.PeekUntil(engine.IsSpaceOrTab)
				return nil, ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, firstToken.Length())
			}
			start, t1Err := klog.NewTimeFromString(startCandidate.ToString())
			if t1Err != nil {
				return nil, ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, startCandidate.Length())
			}
			entryStartPosition := startCandidate.PointerPosition
			entry.Advance(startCandidate.Length())
			style.TimeFormat.Set(start.Format())

			entryStartPositionEnds := entry.PointerPosition
			entry.SkipWhile(engine.Is(' '))
			if entryStartPositionEnds != entry.PointerPosition {
				style.SpacingInRange.Set(strings.Repeat(" ", entry.PointerPosition-entryStartPositionEnds))
			} else {
				style.SpacingInRange.Set("")
			}

			if entry.Peek() != '-' {
				return nil, ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, 1)
			}
			entry.Advance(1) // '-'
			entry.SkipWhile(engine.Is(' '))

			// Check whether the range is open-ended.
			if entry.Peek() == '?' {
				entry.Advance(1)
				placeholder, _ := entry.PeekUntil(engine.IsSpaceOrTab)

				// The placeholder can appear multiple times.
				for _, p := range placeholder.Chars {
					if p != '?' {
						return nil, ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, placeholder.Length())
					}
				}
				entry.Advance(placeholder.Length())
				return func(s klog.EntrySummary) engine.Error {
					sErr := record.StartOpenRange(start, s)
					if sErr != nil {
						return ErrorDuplicateOpenRange().New(entry.Line, entryStartPosition, entry.PointerPosition-entryStartPosition)
					}
					return nil
				}, nil
			}

			// Ultimately, the entry can only be a regular range.
			endCandidate, _ := entry.PeekUntil(engine.IsSpaceOrTab)
			if endCandidate.Length() == 0 {
				return nil, ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, 1)
			}
			end, t2Err := klog.NewTimeFromString(endCandidate.ToString())
			if t2Err != nil {
				return nil, ErrorMalformedEntry().New(entry.Line, entry.PointerPosition, endCandidate.Length())
			}
			entry.Advance(endCandidate.Length())
			timeRange, rErr := klog.NewRange(start, end)
			if rErr != nil {
				return nil, ErrorIllegalRange().New(entry.Line, entryStartPosition, entry.PointerPosition-entryStartPosition)
			}
			return func(s klog.EntrySummary) engine.Error {
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
		entrySummary, esErr := func() (klog.EntrySummary, engine.Error) {
			var result klog.EntrySummary

			// Parse first line of entry summary.
			if engine.IsSpaceOrTab(entry.Peek()) {
				entry.Advance(1)
				summaryText, _ := entry.PeekUntil(engine.Is(engine.END_OF_TEXT))
				firstLine, sErr := klog.NewEntrySummary(summaryText.ToString())
				if sErr != nil {
					return nil, ErrorMalformedSummary().New(summaryText.Line, 0, summaryText.Length())
				}
				result = firstLine
			} else {
				result, _ = klog.NewEntrySummary("")
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
				newEntrySummary, sErr := klog.NewEntrySummary(append(result, additionalText.ToString())...)
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
