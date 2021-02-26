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

// SplitIntoChunksOfLines splits up a text into paragraphs (at blank lines)
// and the paragraphs into its individual lines.
func SplitIntoChunksOfLines(text string) []Chunk {
	var chunks []Chunk
	var currentChunk Chunk
	text = text + "\n"
	for i, l := range strings.Split(text, "\n") {
		if blankTextPattern.MatchString(l) {
			if currentChunk != nil {
				chunks = append(chunks, currentChunk)
			}
			currentChunk = nil
			continue
		}
		l = replaceTabsWithSpaces(l)
		l, spacesCount := trimLeftCount(l, ' ')
		currentIndentation := -1
		if spacesCount == 0 {
			currentIndentation = 0
		} else if spacesCount >= 2 && spacesCount <= 4 {
			currentIndentation = 1
		}
		currentChunk = append(currentChunk, Text{
			Value:            []rune(l),
			IndentationLevel: currentIndentation,
			LineNumber:       i + 1,
		})
	}
	return chunks
}

func replaceTabsWithSpaces(text string) string {
	text, tabsCount := trimLeftCount(text, '\t')
	return strings.Repeat("    ", tabsCount) + text
}

func trimLeftCount(text string, char rune) (string, int) {
	count := 0
	text = strings.TrimLeftFunc(text, func(c rune) bool {
		if c == char {
			count++
			return true
		}
		return false
	})
	return text, count
}
