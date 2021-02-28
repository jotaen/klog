package engine

import (
	"regexp"
	"strings"
)

type Line struct {
	Value            []rune
	Original         string
	PointerPosition  int
	LineNumber       int
	IndentationLevel int
}

var lineDelimiterPattern = regexp.MustCompile(`^.*\n?`)

func Split(text string) []Line {
	var result []Line
	remainder := text
	lineNumber := 0
	for len(remainder) > 0 {
		lineNumber += 1
		original := lineDelimiterPattern.FindString(remainder)
		value, indentation := trimIndentation(original, 0)
		if strings.HasPrefix(value, " ") {
			// A single “dangling” space is not allowed
			value = value[1:]
			indentation = -99999
		}
		result = append(result, Line{
			Value:            []rune(trimLineEnding(value)),
			Original:         original,
			PointerPosition:  0,
			LineNumber:       lineNumber,
			IndentationLevel: indentation,
		})
		remainder = remainder[len(original):]
	}
	return result
}

func Join(ls []Line) string {
	result := ""
	for _, l := range ls {
		result += l.Original
	}
	return result
}

var END_OF_TEXT int32 = -1

func (l *Line) Peek() rune {
	char := SubRune(l.Value, l.PointerPosition, 1)
	if char == nil {
		return END_OF_TEXT
	}
	return char[0]
}

func (l *Line) PeekUntil(isMatch func(rune) bool) (Line, bool) {
	result := Line{
		PointerPosition: l.PointerPosition,
		Value:           nil,
		LineNumber:      l.LineNumber,
	}
	for i := l.PointerPosition; i < len(l.Value); i++ {
		next := SubRune(l.Value, i, 1)
		if isMatch(next[0]) {
			return result, true
		}
		result.Value = append(result.Value, next[0])
	}
	return result, false
}

func (l *Line) Advance(increment int) {
	l.PointerPosition += increment
}

func (l *Line) SkipWhitespace() {
	for IsWhitespace(l.Peek()) {
		l.Advance(1)
	}
	return
}

func (l *Line) Length() int {
	return len(l.Value)
}

func (l *Line) RemainingLength() int {
	return l.Length() - l.PointerPosition
}

func (l *Line) ToString() string {
	return string(l.Value)
}
