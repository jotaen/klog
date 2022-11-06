package txt

// Block is multiple consecutive lines with text, with no blank lines
// in between, but possibly one or more blank lines before or after.
// It’s basically like a paragraph of text, with surrounding whitespace.
// The Block is guaranteed to contain exactly a single sequence of
// significant lines, i.e. lines that contain text.
type Block interface {
	// Lines returns all lines.
	Lines() []Line

	// SignificantLines returns the lines that are not blank. The two integers
	// are the number of insignificant lines at the beginning and the end.
	SignificantLines() ([]Line, int, int)

	// OverallLineIndex returns the overall line index, taking into
	// account the context of all preceding blocks.
	OverallLineIndex(int) int

	// SetPrecedingLineCount adjusts the overall line count.
	SetPrecedingLineCount(int)
}

type block struct {
	precedingLineCount int
	lines              []Line
}

// ParseBlock parses a block from the beginning of a text. It returns
// the parsed block, along with the number of bytes consumed from the
// string. If the text doesn’t contain significant lines, it returns nil.
func ParseBlock(text string, precedingLineCount int) (Block, int) {
	const (
		MODE_PRECEDING_BLANK_LINES = iota
		MODE_SIGNIFICANT_LINES
		MODE_TRAILING_BLANK_LINES
	)

	var lines []Line
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
		line := NewLineFromString(currentLine)

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
		lines = append(lines, line)
		bytesConsumed += len(currentLine)
		currentLineStart = nextChar
	}

	hasSignificantLines := currentMode != MODE_PRECEDING_BLANK_LINES
	if !hasSignificantLines {
		return nil, bytesConsumed
	}
	return &block{precedingLineCount, lines}, bytesConsumed
}

func (b *block) OverallLineIndex(lineIndex int) int {
	return b.precedingLineCount + lineIndex
}

func (b *block) SetPrecedingLineCount(count int) {
	b.precedingLineCount = count
}

func (b *block) Lines() []Line {
	return b.lines
}

func (b *block) SignificantLines() ([]Line, int, int) {
	first, last := 0, len(b.lines)
	hasSeenSignificant := false
	for i, l := range b.lines {
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
	return b.lines[first:last], first, last
}
