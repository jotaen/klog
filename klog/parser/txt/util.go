package txt

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

// IsSpaceOrTab checks whether a rune is a space or a tab character.
func IsSpaceOrTab(r rune) bool {
	return r == ' ' || r == '\t'
}

func Is(matchingCharacter ...rune) func(rune) bool {
	return func(r rune) bool {
		for _, m := range matchingCharacter {
			if m == r {
				return true
			}
		}
		return false
	}
}
