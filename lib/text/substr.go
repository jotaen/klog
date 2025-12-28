package text

import "unicode/utf8"

// TextSubstrWithContext returns a fragment of a string like a regular `substr`
// method would do. However, it returns a bit of surrounding text for context.
// The surrounding text is between `minSurroundingRunes` and `maxSurroundingRunes`
// long, and it tries to find a word boundary (space character) as natural cut-off.
// Only if it cannot find one, it makes a hard cut.
func TextSubstrWithContext(text string, start int, length int, minSurroundingRunes int, maxSurroundingRunes int) (string, int) {
	if start < 0 || length < 0 || start >= len(text) {
		return "", 0
	}

	end := start + length
	if end > len(text) {
		end = len(text)
	}

	fuzzyStart := start
	charCount := 0

	for fuzzyStart > 0 && charCount < maxSurroundingRunes {
		_, size := utf8.DecodeLastRuneInString(text[:fuzzyStart])
		if size == 0 {
			break
		}
		fuzzyStart -= size
		charCount++

		if charCount >= minSurroundingRunes && text[fuzzyStart] == ' ' {
			break
		}
	}

	if fuzzyStart < len(text) && text[fuzzyStart] == ' ' {
		fuzzyStart++
	}

	fuzzyEnd := end
	charCount = 0

	for fuzzyEnd < len(text) && charCount < maxSurroundingRunes {
		r, size := utf8.DecodeRuneInString(text[fuzzyEnd:])
		if r == utf8.RuneError && size == 1 {
			break
		}

		if charCount >= minSurroundingRunes && r == ' ' {
			break
		}

		fuzzyEnd += size
		charCount++
	}

	translatedPos := start - fuzzyStart

	return text[fuzzyStart:fuzzyEnd], translatedPos
}
