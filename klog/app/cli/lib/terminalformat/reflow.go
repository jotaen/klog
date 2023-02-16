package terminalformat

import "strings"

type Reflower struct {
	maxLength int
	newLine   string
}

func NewReflower(maxLineLength int, newLineChar string) Reflower {
	return Reflower{
		maxLength: maxLineLength,
		newLine:   newLineChar,
	}
}

func (b Reflower) Reflow(text string, linePrefixes []string) string {
	SPACE := " "
	var resultParagraphs []string

	for _, paragraph := range strings.Split(text, b.newLine) {
		words := strings.Split(paragraph, SPACE)
		lines := []string{""}
		currentLinePrefix := ""
		for i, word := range words {
			nr := len(lines) - 1
			isLastWordOfText := i == len(words)-1
			if !isLastWordOfText && len(lines[nr])+len(words[i+1]) > b.maxLength {
				lines = append(lines, "")
				nr = len(lines) - 1
			}
			if lines[nr] == "" {
				if len(linePrefixes) > nr {
					currentLinePrefix = linePrefixes[nr]
				}
				lines[nr] += currentLinePrefix
			} else {
				lines[nr] += SPACE
			}
			lines[nr] += word
		}
		resultParagraphs = append(resultParagraphs, strings.Join(lines, b.newLine))
	}
	return strings.Join(resultParagraphs, b.newLine)
}
