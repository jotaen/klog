package engine

import "strings"

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

func GroupIntoBlocks(lines []Line) [][]Line {
	var blocks [][]Line
	var currentBlock []Line
	for _, l := range lines {
		if IsBlank(l) {
			if currentBlock != nil {
				blocks = append(blocks, currentBlock)
				currentBlock = nil
			}
			continue
		}
		currentBlock = append(currentBlock, l)
	}
	if currentBlock != nil {
		blocks = append(blocks, currentBlock)
	}
	return blocks
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

func trimLineEnding(line string) string {
	return strings.TrimFunc(line, func(r rune) bool {
		return r == '\n' || r == '\r'
	})
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
