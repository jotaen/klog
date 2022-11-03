package engine

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

// Split breaks up text into a list of Line’s. The text must use `\n` as
// line delimiters.
func Split(text string) []Line {
	var result []Line
	var currentLine string
	lineNumber := 1
	for i, char := range text {
		currentLine += string(char)
		if char == '\n' || i+1 == len(text) {
			result = append(result, NewLineFromString(currentLine, lineNumber))
			currentLine = ""
			lineNumber++
		}
	}
	return result
}

var lineEndingPatterns = []string{"\r\n", "\n"}

func splitOffLineEnding(text string) (string, string) {
	for _, e := range lineEndingPatterns {
		if strings.HasSuffix(text, e) {
			return text[:len(text)-len(e)], e
		}
	}
	return text, ""
}
