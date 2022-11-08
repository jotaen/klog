/*
Package parser contains the logic how to convert Record objects from and to plain text.
*/
package parser

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/parser/txt"
)

func parse(block txt.Block) (klog.Record, []txt.Error) {
	lines, initialLineOffset, _ := block.SignificantLines()
	initialLineCount := len(lines) // Capture current value
	nr := func(lines []txt.Line) int {
		return initialLineOffset + initialLineCount - len(lines)
	}
	var errs []txt.Error

	// ========== HEADLINE ==========
	record := func() klog.Record {
		headline := txt.NewParseable(lines[0], 0)

		// There is no leading whitespace allowed in the headline.
		if txt.IsSpaceOrTab(headline.Peek()) {
			errs = append(errs, ErrorIllegalIndentation().New(block, nr(lines), 0, headline.Length()))
			return nil
		}

		// Parse the date.
		dateText, _ := headline.PeekUntil(txt.IsSpaceOrTab)
		rDate, dErr := klog.NewDateFromString(dateText.ToString())
		if dErr != nil {
			errs = append(errs, ErrorInvalidDate().New(block, nr(lines), headline.PointerPosition, dateText.Length()))
			return nil
		}
		headline.Advance(dateText.Length())
		headline.SkipWhile(txt.IsSpaceOrTab)
		r := klog.NewRecord(rDate)

		// Check if there is a should-total set, and if so, parse it.
		if headline.Peek() == '(' {
			headline.Advance(1) // '('
			headline.SkipWhile(txt.IsSpaceOrTab)
			allPropsText, hasClosingParenthesis := headline.PeekUntil(txt.Is(')'))
			if !hasClosingParenthesis {
				errs = append(errs, ErrorMalformedPropertiesSyntax().New(block, nr(lines), headline.Length(), 1))
				return r
			}
			if allPropsText.Length() == 0 {
				errs = append(errs, ErrorMalformedPropertiesSyntax().New(block, nr(lines), headline.PointerPosition, 1))
				return r
			}
			shouldTotalText, hasExclamationMark := headline.PeekUntil(txt.Is('!'))
			if !hasExclamationMark {
				errs = append(errs, ErrorUnrecognisedProperty().New(block, nr(lines), headline.PointerPosition, shouldTotalText.Length()-1))
				return r
			}
			shouldTotal, sErr := klog.NewDurationFromString(shouldTotalText.ToString())
			if sErr != nil {
				errs = append(errs, ErrorMalformedShouldTotal().New(block, nr(lines), headline.PointerPosition, shouldTotalText.Length()))
				return r
			}
			r.SetShouldTotal(shouldTotal)
			headline.Advance(shouldTotalText.Length())
			headline.Advance(1) // '!'
			headline.SkipWhile(txt.IsSpaceOrTab)

			// Make sure there is no other text between the braces.
			if headline.Peek() != ')' {
				errs = append(errs, ErrorUnrecognisedProperty().New(block, nr(lines), headline.PointerPosition, headline.RemainingLength()-1))
				return r
			}
			headline.Advance(1) // ')'
		}

		// Make sure there is no other text left in the headline.
		headline.SkipWhile(txt.IsSpaceOrTab)
		if headline.RemainingLength() > 0 {
			errs = append(errs, ErrorUnrecognisedTextInHeadline().New(block, nr(lines), headline.PointerPosition, headline.RemainingLength()))
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

	var indentator *txt.Indentator

	// ========== SUMMARY LINES ==========
	for _, l := range lines {
		indentator = txt.NewIndentator(txt.Indentations, lines[0])
		if indentator != nil {
			break
		}
		summary := txt.NewParseable(l, 0)
		newSummary, sErr := klog.NewRecordSummary(append(record.Summary(), summary.ToString())...)
		if sErr != nil {
			errs = append(errs, ErrorMalformedSummary().New(block, nr(lines), 0, summary.Length()))
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

		// Check for correct indentation.
		entry := indentator.NewIndentedParseable(l, 1)
		if entry == nil || txt.IsSpaceOrTab(entry.Peek()) {
			errs = append(errs, ErrorIllegalIndentation().New(block, nr(lines), 0, len(l.Text)))
			break
		}

		// Parse entry value.
		createEntry, evErr := func() (func(klog.EntrySummary) txt.Error, txt.Error) {
			// Try to interpret the entry value as duration.
			durationCandidate, _ := entry.PeekUntil(txt.IsSpaceOrTab)
			duration, dErr := klog.NewDurationFromString(durationCandidate.ToString())
			if dErr == nil {
				entry.Advance(durationCandidate.Length())
				return func(s klog.EntrySummary) txt.Error {
					record.AddDuration(duration, s)
					return nil
				}, nil
			}

			// If the entry value isn’t a duration, it must be the start time of a range.
			startCandidate, _ := entry.PeekUntil(txt.Is('-', ' '))
			if startCandidate.Length() == 0 {
				// Handle case where `-` appears right at the beginning of the line.
				firstToken, _ := entry.PeekUntil(txt.IsSpaceOrTab)
				return nil, ErrorMalformedEntry().New(block, nr(lines), entry.PointerPosition, firstToken.Length())
			}
			start, t1Err := klog.NewTimeFromString(startCandidate.ToString())
			if t1Err != nil {
				return nil, ErrorMalformedEntry().New(block, nr(lines), entry.PointerPosition, startCandidate.Length())
			}
			entryStartPosition := startCandidate.PointerPosition
			entry.Advance(startCandidate.Length())

			entryStartPositionEnds := entry.PointerPosition
			entry.SkipWhile(txt.Is(' '))
			hasRangeSpacesAroundDash := true
			if entryStartPositionEnds == entry.PointerPosition {
				hasRangeSpacesAroundDash = false
			}

			if entry.Peek() != '-' {
				return nil, ErrorMalformedEntry().New(block, nr(lines), entry.PointerPosition, 1)
			}
			entry.Advance(1) // '-'
			entry.SkipWhile(txt.Is(' '))

			// Check whether the range is open-ended.
			if entry.Peek() == '?' {
				entry.Advance(1)
				placeholderRepetition, _ := entry.PeekUntil(txt.IsSpaceOrTab)

				// The placeholder can appear multiple times.
				for _, p := range placeholderRepetition.Chars {
					if p != '?' {
						return nil, ErrorMalformedEntry().New(block, nr(lines), entry.PointerPosition, placeholderRepetition.Length())
					}
				}
				entry.Advance(placeholderRepetition.Length())
				lineNr := nr(lines) // Capture state of `line` at time of function creation.
				return func(s klog.EntrySummary) txt.Error {
					or := klog.NewOpenRangeWithFormat(start, klog.OpenRangeFormat{
						UseSpacesAroundDash:        hasRangeSpacesAroundDash,
						AdditionalPlaceholderChars: placeholderRepetition.Length(),
					})
					sErr := record.Start(or, s)
					if sErr != nil {
						return ErrorDuplicateOpenRange().New(block, lineNr, entryStartPosition, entry.PointerPosition-entryStartPosition)
					}
					return nil
				}, nil
			}

			// Ultimately, the entry can only be a regular range.
			endCandidate, _ := entry.PeekUntil(txt.IsSpaceOrTab)
			if endCandidate.Length() == 0 {
				return nil, ErrorMalformedEntry().New(block, nr(lines), entry.PointerPosition, 1)
			}
			end, t2Err := klog.NewTimeFromString(endCandidate.ToString())
			if t2Err != nil {
				return nil, ErrorMalformedEntry().New(block, nr(lines), entry.PointerPosition, endCandidate.Length())
			}
			entry.Advance(endCandidate.Length())
			timeRange, rErr := klog.NewRangeWithFormat(start, end, klog.RangeFormat{UseSpacesAroundDash: hasRangeSpacesAroundDash})
			if rErr != nil {
				return nil, ErrorIllegalRange().New(block, nr(lines), entryStartPosition, entry.PointerPosition-entryStartPosition)
			}
			return func(s klog.EntrySummary) txt.Error {
				record.AddRange(timeRange, s)
				return nil
			}, nil
		}()
		lines = lines[1:]

		// Check for error while parsing the entry value.
		if evErr != nil {
			errs = append(errs, evErr)
			continue
		}

		// Parse entry summary.
		entrySummary, esErr := func() (klog.EntrySummary, txt.Error) {
			var result klog.EntrySummary

			// Parse first line of entry summary.
			if txt.IsSpaceOrTab(entry.Peek()) {
				entry.Advance(1)
				summaryText := entry.Remainder()
				firstLine, sErr := klog.NewEntrySummary(summaryText.ToString())
				if sErr != nil {
					return nil, ErrorMalformedSummary().New(block, nr(lines), 0, summaryText.Length())
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
					return nil, ErrorMalformedSummary().New(block, nr(lines), 0, nextEntrySummaryLine.Length())
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
		return nil, errs
	}
	return record, nil
}
