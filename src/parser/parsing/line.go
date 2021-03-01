package parsing

import (
	"regexp"
	"strings"
)

type Line struct {
	Text        string
	LineNumber  int
	lineEnding  string
	indentation string
}

var lineDelimiterPattern = regexp.MustCompile(`^.*\n?`)

func NewLineFromString(lineText string, lineNumber int) Line {
	text, indentation := splitOffPrecedingWhitespace(lineText)
	text, lineEnding := splitOffLineEnding(text)
	return Line{
		Text:        text,
		LineNumber:  lineNumber,
		lineEnding:  lineEnding,
		indentation: indentation,
	}
}

func (l *Line) ToString() string {
	return l.indentation + l.Text + l.lineEnding
}

func (l *Line) IndentationLevel() int {
	normalised := strings.ReplaceAll(l.indentation, "\t", "    ")
	if normalised == "" {
		return 0
	}
	if len(normalised) == 1 {
		return -1
	}
	if len(normalised) > 4 {
		return 2
	}
	return 1
}

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

func Join(ls []Line) string {
	result := ""
	preferredLineEnding := "\n"
	for _, l := range ls {
		result += l.ToString()
		if l.lineEnding == "" {
			result += preferredLineEnding
		} else {
			preferredLineEnding = l.lineEnding
		}
	}
	return result
}

func IsBlank(l Line) bool {
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

func Insert(ls []Line, position int, lineText string) []Line {
	if position > len(ls)+1 {
		panic("Out of bounds")
	}
	result := make([]Line, len(ls)+1)
	offset := 0
	for i := range result {
		if i == position {
			result[i] = NewLineFromString(lineText, i+1)
			offset = 1
		} else {
			result[i] = ls[i-offset]
		}
		result[i].LineNumber = i + 1
	}
	return result
}
