package txt

import (
	"strings"
)

// Line is a data structure that represent one line of the source file.
type Line struct {
	// Text contains the copy of the line.
	Text string

	// LineNumber is the line number, starting with 1.
	LineNumber int

	// LineEnding is the encountered line ending sequence `\n` or `\r\n`.
	// Note that for the last line in a file, there might be no line ending.
	LineEnding string
}

var LineEndings = []string{"\r\n", "\n"}

var Indentations = []string{"    ", "   ", "  ", "\t"}

// NewLineFromString turns data into a Line object.
func NewLineFromString(rawLineText string, lineNumber int) Line {
	text, lineEnding := splitOffLineEnding(rawLineText)
	return Line{
		Text:       text,
		LineNumber: lineNumber,
		LineEnding: lineEnding,
	}
}

// Original returns the (byte-wise) identical line of text as it appeared in the file.
func (l *Line) Original() string {
	return l.Text + l.LineEnding
}

// IsBlank checks whether a line is all spaces or tabs.
func (l *Line) IsBlank() bool {
	if len(l.Text) == 0 {
		return true
	}
	for _, c := range l.Text {
		if c != ' ' && c != '\t' {
			return false
		}
	}
	return true
}

// Indentation returns the indentation sequence of a line that is indented at least once.
// Note: it cannot determine the level of indentation. If the line is not indented,
// it returns empty string.
func (l *Line) Indentation() string {
	for _, i := range Indentations {
		if strings.HasPrefix(l.Original(), i) {
			return i
		}
	}
	return ""
}

func splitOffLineEnding(text string) (string, string) {
	for _, e := range LineEndings {
		if strings.HasSuffix(text, e) {
			return text[:len(text)-len(e)], e
		}
	}
	return text, ""
}
