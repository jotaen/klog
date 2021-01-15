package engine

func substr(text []rune, start int, length int) string {
	if start >= len(text) {
		return ""
	}
	if start+length > len(text) {
		length = len(text) - start
	}
	return string(text[start : start+length])
}
