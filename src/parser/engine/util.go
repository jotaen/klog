package engine

// SubRune returns a subset of a rune list. It might be shorter than
// the requested length, if the text doesn’t contain enough characters.
// It returns empty, if the start position is bigger than the length.
func SubRune(text []rune, start int, length int) []rune {
	if start >= len(text) {
		return nil
	}
	if start+length > len(text) {
		length = len(text) - start
	}
	return text[start : start+length]
}

// IsWhitespace checks whether a rune is a space or a tab.
func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}
