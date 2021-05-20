package terminalformat

import "regexp"

var ansiSequencePattern = regexp.MustCompile(`\x1b\[[\d;]+m`)

func StripAllAnsiSequences(text string) string {
	return ansiSequencePattern.ReplaceAllString(text, "")
}
