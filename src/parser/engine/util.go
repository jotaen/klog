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

var blankTextPattern = regexp.MustCompile(
	`^[  \t]*$`, // match space, non-breaking space, tab
)

func SplitIntoChunksOfLines(text string) []Chunk {
	var chunks []Chunk
	var currentChunk Chunk
	currentIndentation := 0
	for nr, l := range strings.Split(text, "\n") {
		if blankTextPattern.MatchString(l) {
			if currentChunk != nil {
				chunks = append(chunks, currentChunk)
			}
			currentChunk = nil
			currentIndentation = 0
			continue
		}
		if regexp.MustCompile(`^\t`).MatchString(l) {
			currentIndentation = 1
		}
		currentChunk = append(currentChunk, Text{
			Value:            []rune(l)[currentIndentation:],
			IndentationLevel: currentIndentation,
			LineNumber:       nr,
		})
	}
	return chunks
}
