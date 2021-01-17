package engine

import (
	"regexp"
	"strings"
)

func SubRune(text []rune, start int, length int) []rune {
	if start >= len(text) {
		return nil
	}
	if start+length > len(text) {
		length = len(text) - start
	}
	return text[start : start+length]
}

func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}

var blankTextPattern = regexp.MustCompile(
	`^[  \t]*$`, // match space, non-breaking space, tab
)

var indentationPattern = regexp.MustCompile(`^(\t| {2,4})`)

func SplitIntoChunksOfLines(text string) []Chunk {
	var chunks []Chunk
	var currentChunk Chunk
	currentIndentation := 0
	text = text + "\n"
	for i, l := range strings.Split(text, "\n") {
		if blankTextPattern.MatchString(l) {
			if currentChunk != nil {
				chunks = append(chunks, currentChunk)
			}
			currentChunk = nil
			continue
		}
		indent := indentationPattern.FindString(l)
		if indent == "" {
			currentIndentation = 0
		} else {
			currentIndentation = 1
		}
		currentChunk = append(currentChunk, Text{
			Value:            []rune(l)[len(indent):],
			IndentationLevel: currentIndentation,
			LineNumber:       i + 1,
		})
	}
	return chunks
}
