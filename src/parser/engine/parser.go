package engine

import (
	"errors"
	"regexp"
	"strings"
)

type (
	Chunk   []Line
	ParseFn func(Chunk) error
)

var blankLinePattern = regexp.MustCompile(
	`^[  \t]*$`, // match space, non-breaking space, tab
)

func Parse(text string, fn ParseFn) error {
	p := splitIntoChunksOfLines(text)
	var errs []error
	for _, chunk := range p {
		err := fn(chunk)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.New("PARSER_ERROR")
}

func splitIntoChunksOfLines(text string) []Chunk {
	var chunks []Chunk
	var currentChunk Chunk
	for nr, l := range strings.Split(text, "\n") {
		if blankLinePattern.MatchString(l) {
			if currentChunk != nil {
				chunks = append(chunks, currentChunk)
			}
			currentChunk = nil
			continue
		}
		currentChunk = append(currentChunk, Line{
			Text: []rune(l),
			Nr:   nr,
		})
	}
	return chunks
}
