package txt

// Block is multiple consecutive lines with text, with no blank lines
// in between, but possibly one or more blank lines before or after.
// It’s basically like a paragraph of text, with surrounding whitespace.
// The Block is guaranteed to contain exactly a single sequence of
// significant lines, i.e. lines that contain text.
type Block []Line

// ParseBlock parses a block from the beginning of a text. It returns
// the parsed block, along with the number of bytes consumed from the
// string. If the text doesn’t contain significant lines, it returns nil.
func ParseBlock(text string, initialLineNumber int) (Block, int) {
	const (
		MODE_PRECEDING_BLANK_LINES = iota
		MODE_SIGNIFICANT_LINES
		MODE_TRAILING_BLANK_LINES
	)

	var block Block
	bytesConsumed := 0
	currentLineStart := 0
	currentMode := MODE_PRECEDING_BLANK_LINES

	// Parse text line-wise.
parsingLoop:
	for i, char := range text { // Note: char is a UTF-8 rune
		if char != '\n' && i+1 != len(text) {
			continue
		}

		// Process line.
		nextChar := i + len(string(char))
		currentLine := text[currentLineStart:nextChar]
		line := NewLineFromString(currentLine, initialLineNumber)

		switch currentMode {
		case MODE_PRECEDING_BLANK_LINES:
			if !line.IsBlank() {
				currentMode = MODE_SIGNIFICANT_LINES
			}
		case MODE_SIGNIFICANT_LINES:
			if line.IsBlank() {
				currentMode = MODE_TRAILING_BLANK_LINES
			}
		case MODE_TRAILING_BLANK_LINES:
			if !line.IsBlank() {
				break parsingLoop
			}
		}
		block = append(block, line)
		bytesConsumed += len(currentLine)
		currentLineStart = nextChar
		initialLineNumber++
	}

	hasSignificantLines := currentMode != MODE_PRECEDING_BLANK_LINES
	if !hasSignificantLines {
		block = nil
	}
	return block, bytesConsumed
}

// SignificantLines returns the lines that are not blank.
func (b Block) SignificantLines() []Line {
	first, last := 0, len(b)
	hasSeenSignificant := false
	for i, l := range b {
		if !hasSeenSignificant && !l.IsBlank() {
			first = i
			hasSeenSignificant = true
			continue
		}
		if hasSeenSignificant && l.IsBlank() {
			last = i
			break
		}
	}
	return b[first:last]
}
