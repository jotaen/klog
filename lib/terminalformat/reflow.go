package terminalformat

import (
	"strings"
	"unicode"
)

type Reflower struct {
	maxWidth int
}

// NewReflower creates a new Reflower with the given maximum line length.
func NewReflower(maxWidth int) *Reflower {
	if maxWidth <= 0 {
		maxWidth = 80
	}
	return &Reflower{maxWidth: maxWidth}
}

// Reflow reflows text to meet a maximum line length. It preserves and respects
// existing newlines and existing line prefixes. A line prefix can also be
// added optionally.
func (r *Reflower) Reflow(text string, linePrefix string) string {
	if text == "" {
		return ""
	}

	prefixLen := len([]rune(linePrefix))
	lines := strings.Split(text, "\n")
	var result []string

	for _, line := range lines {
		// Preserve blank lines.
		if len(strings.TrimSpace(line)) == 0 {
			result = append(result, linePrefix+line)
			continue
		}

		// Extract leading whitespace.
		leadingWS := getLeadingWhitespace(line)
		leadingWSLen := len([]rune(leadingWS))
		content := strings.TrimLeft(line, " \t")

		// Calculate available width: maxWidth - linePrefix - leadingWhitespace.
		availableWidth := r.maxWidth - prefixLen - leadingWSLen

		// Edge case: prefix + indent exceeds or equals maxWidth.
		// Ensure at least 1 character of space for content.
		if availableWidth <= 0 {
			availableWidth = 1
		}

		// Reflow the line content.
		reflowed := reflowLine(content, availableWidth, linePrefix, leadingWS)
		result = append(result, reflowed...)
	}

	return strings.Join(result, "\n")
}

func getLeadingWhitespace(line string) string {
	for i, ch := range line {
		if !unicode.IsSpace(ch) || ch == '\n' {
			return line[:i]
		}
	}
	return line
}

func reflowLine(content string, maxWidth int, linePrefix string, indent string) []string {
	if maxWidth <= 0 {
		maxWidth = 1
	}
	if content == "" {
		return []string{linePrefix + indent}
	}
	words := strings.Fields(content)
	if len(words) == 0 {
		return []string{linePrefix + indent}
	}

	var result []string
	currentLine := linePrefix + indent
	currentLen := 0

	for _, word := range words {
		wordRunes := []rune(word)
		wordLen := len(wordRunes)

		if wordLen > maxWidth {
			if currentLen > 0 {
				result = append(result, currentLine)
				currentLine = linePrefix + indent
				currentLen = 0
			}

			for len(wordRunes) > maxWidth {
				result = append(result, linePrefix+indent+string(wordRunes[:maxWidth]))
				wordRunes = wordRunes[maxWidth:]
			}

			currentLine = linePrefix + indent + string(wordRunes)
			currentLen = len(wordRunes)
			continue
		}

		spaceNeeded := wordLen
		if currentLen > 0 {
			spaceNeeded++
		}

		if currentLen+spaceNeeded > maxWidth {
			result = append(result, currentLine)
			currentLine = linePrefix + indent + word
			currentLen = wordLen
		} else {
			if currentLen > 0 {
				currentLine += " "
				currentLen++
			}
			currentLine += word
			currentLen += wordLen
		}
	}

	if currentLen > 0 || len(result) == 0 {
		result = append(result, currentLine)
	}

	return result
}
