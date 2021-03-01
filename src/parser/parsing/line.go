package parsing

import (
	"regexp"
	"strings"
)

type Line struct {
	Text                string
	LineNumber          int
	originalLineEnding  string
	originalIndentation string
}

var lineDelimiterPattern = regexp.MustCompile(`^.*\n?`)

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

func (l *Line) Original() string {
	return l.originalIndentation + l.Text + l.originalLineEnding
}

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
	for _, l := range ls {
		result += l.Original()
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

func Insert(ls []Line, position int, text string, isIndented bool, prefs Preferences) []Line {
	if position > len(ls)+1 {
		panic("Out of bounds")
	}
	result := make([]Line, len(ls)+1)
	offset := 0
	for i := range result {
		if i == position {
			line := ""
			if isIndented {
				line += prefs.Indentation
			}
			line += text + prefs.LineEnding
			result[i] = NewLineFromString(line, i+1)
			offset = 1
		} else {
			result[i] = ls[i-offset]
		}
		result[i].LineNumber = i + 1
	}
	if position > 0 && result[position-1].originalLineEnding == "" {
		result[position-1].originalLineEnding = prefs.LineEnding
	}
	return result
}

type Preferences struct {
	LineEnding  string
	Indentation string
}

func DefaultPreferences() Preferences {
	return Preferences{
		LineEnding:  "\n",
		Indentation: "    ",
	}
}

func (p *Preferences) Adapt(l *Line) {
	if l.IndentationLevel() == 1 {
		p.Indentation = l.originalIndentation
	}
	if l.originalLineEnding != "" {
		p.LineEnding = l.originalLineEnding
	}
}
