package parsing

import (
	"regexp"
	"strings"
)

// Line is a data structure that represent one line of the source file.
type Line struct {
	// Text contains the copy of the line.
	Text string

	// LineNumber is the line number, starting with 1.
	LineNumber int

	// originalLineEnding is the encountered line ending sequence `\n` or `\r\n`.
	originalLineEnding string

	// originalIndentation is the exact whitespace sequence used for indentation.
	originalIndentation string
}

var lineDelimiterPattern = regexp.MustCompile(`^.*\n?`)

// NewLineFromString turns data into a Line object.
func NewLineFromString(rawLineText string, lineNumber int) Line {
	text, indentation := splitOffPrecedingWhitespace(rawLineText)
	text, lineEnding := splitOffLineEnding(text)
	return Line{
		Text:                text,
		LineNumber:          lineNumber,
		originalLineEnding:  lineEnding,
		originalIndentation: indentation,
	}
}

// Original returns the (byte-wise) identical line of text as it appeared in the file.
func (l *Line) Original() string {
	return l.originalIndentation + l.Text + l.originalLineEnding
}

// IndentationLevel returns `0` for top level, `1` for first level, and `-1` for illegal indentation styles.
func (l *Line) IndentationLevel() int {
	normalised := strings.ReplaceAll(l.originalIndentation, "\t", "    ")
	if normalised == "" {
		return 0
	}
	if len(normalised) == 1 || len(normalised) > 4 {
		return -1
	}
	return 1
}

// Split breaks up text into a list of Line’s. The text must use `\n` as
// line delimiters.
func Split(text string) []Line {
	var result []Line
	remainder := text
	lineNumber := 0
	for len(remainder) > 0 {
		lineNumber += 1
		original := lineDelimiterPattern.FindString(remainder)
		result = append(result, NewLineFromString(original, lineNumber))
		remainder = remainder[len(original):]
	}
	return result
}

// Join restores a blob of text from Line’s. The result is (byte-wise) identical
// to the original copy.
func Join(ls []Line) string {
	result := ""
	for _, l := range ls {
		result += l.Original()
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

func splitOffPrecedingWhitespace(line string) (string, string) {
	text := strings.TrimLeftFunc(line, func(r rune) bool {
		return r == '\t' || r == ' '
	})
	return text, line[:len(line)-len(text)]
}
