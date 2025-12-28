package kfl

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

type ParseError interface {
	error
	Original() error
}

type parseError struct {
	err      error
	position int
	length   int
	query    string
}

func (e parseError) Error() string {
	errorLength := int(math.Max(float64(e.length), 1))
	relevantQueryFragment, newStart := fuzzySubstr(e.query, e.position, errorLength)
	return fmt.Sprintf(
		// TODO remove   once reflower fix has been rebased in.
		"%s\n\n    %s\n    %s%s%s\n    (Char %d in query.)",
		e.err,
		relevantQueryFragment,
		strings.Repeat("—", newStart),
		strings.Repeat("^", errorLength),
		strings.Repeat("—", len(relevantQueryFragment)-(newStart+errorLength)),
		e.position,
	)
}

func (e parseError) Original() error {
	return e.err
}

func fuzzySubstr(text string, start int, length int) (string, int) {
	if start < 0 || length < 0 || start >= len(text) {
		return "", 0
	}

	// Clamp the end position to the text length
	end := start + length
	if end > len(text) {
		end = len(text)
	}

	// Find fuzzy start: go back at least 10 chars, up to 20, stop at first space after 10
	fuzzyStart := start
	charCount := 0

	for fuzzyStart > 0 && charCount < 20 {
		// Move back one rune
		_, size := utf8.DecodeLastRuneInString(text[:fuzzyStart])
		if size == 0 {
			break
		}
		fuzzyStart -= size
		charCount++

		// If we've gone at least 10 chars and hit a space, stop here
		if charCount >= 10 && text[fuzzyStart] == ' ' {
			break
		}
	}

	// Find fuzzy end: go forward at least 10 chars, up to 20, stop at first space after 10
	fuzzyEnd := end
	charCount = 0

	for fuzzyEnd < len(text) && charCount < 20 {
		r, size := utf8.DecodeRuneInString(text[fuzzyEnd:])
		if r == utf8.RuneError && size == 1 {
			break
		}

		// If we've gone at least 10 chars and hit a space, stop here
		if charCount >= 10 && r == ' ' {
			break
		}

		fuzzyEnd += size
		charCount++
	}

	// Calculate the translated position (where 'start' is in the returned substring)
	translatedPos := start - fuzzyStart

	return text[fuzzyStart:fuzzyEnd], translatedPos
}
