package engine

import (
	"regexp"
	"strings"
)

type (
	Chunk   []Text
	ParseFn func(Chunk) error
)

var blankTextPattern = regexp.MustCompile(
	`^[  \t]*$`, // match space, non-breaking space, tab
)

func Parse(text string, fn ParseFn) []error {
	p := splitIntoChunksOfLines(text)
	var errs []error
	for _, chunk := range p {
		err := fn(chunk)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func splitIntoChunksOfLines(text string) []Chunk {
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
