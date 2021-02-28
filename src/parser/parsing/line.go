package parsing

import (
	"regexp"
	"strings"
)

type Line struct {
	Value            []rune
	Original         string
	LineNumber       int
	IndentationLevel int
}

var lineDelimiterPattern = regexp.MustCompile(`^.*\n?`)

func NewLineFromString(lineText string, initialIndentation int, lineNumber int) Line {
	value, indentation := trimIndentation(lineText, initialIndentation)
	if strings.HasPrefix(value, " ") {
		// A single “dangling” space is not allowed
		value = value[1:]
		indentation = -99999
	}
	return Line{
		Value:            []rune(trimLineEnding(value)),
		Original:         lineText,
		LineNumber:       lineNumber,
		IndentationLevel: indentation,
	}
}

func IsNewLineTerminated(l Line) bool {
	return strings.LastIndexFunc(l.Original, IsNewline) != -1
}

func Split(text string) []Line {
	var result []Line
	remainder := text
	lineNumber := 0
	for len(remainder) > 0 {
		lineNumber += 1
		original := lineDelimiterPattern.FindString(remainder)
		result = append(result, NewLineFromString(original, 0, lineNumber))
		remainder = remainder[len(original):]
	}
	return result
}

func trimLineEnding(line string) string {
	return strings.TrimFunc(line, IsNewline)
}

var indentationPatterns = []string{"    ", "   ", "  ", "\t"}

func trimIndentation(line string, initialIndentation int) (string, int) {
	for _, indent := range indentationPatterns {
		if strings.HasPrefix(line, indent) {
			return trimIndentation(line[len(indent):], initialIndentation+1)
		}
	}
	return line, initialIndentation
}

func Join(ls []Line) string {
	result := ""
	for _, l := range ls {
		result += l.Original
	}
	return result
}

func IsBlank(l Line) bool {
	if len(l.Value) == 0 {
		return true
	}
	for _, c := range l.Value {
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
			result[i] = NewLineFromString(lineText, 0, i+1)
			offset = 1
			if i > 0 && !IsNewLineTerminated(result[i-1]) {
				result[i].Original = "\n" + result[i].Original
			}
		} else {
			result[i] = ls[i-offset]
		}
		result[i].LineNumber = i + 1
	}
	return result
}
