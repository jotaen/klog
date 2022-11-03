package txt

// Block is multiple consecutive lines with text, with no blank lines
// in between, but possibly one or more blank lines before or after.
// It’s basically like a paragraph of text, with surrounding whitespace.
type Block []Line

// SignificantLines returns the lines that are not blank.
func (b Block) SignificantLines() []Line {
	var linesWithContent []Line
	for _, l := range b {
		if !isBlank(l) {
			linesWithContent = append(linesWithContent, l)
		}
	}
	return linesWithContent
}

// GroupIntoBlocks splits up lines into Block’s.
func GroupIntoBlocks(lines []Line) []Block {
	const (
		MODE_PRECEDING_BLANK_LINES = iota
		MODE_SIGNIFICANT_LINES
		MODE_TRAILING_BLANK_LINES
	)
	var blocks []Block
	var currentBlock Block
	currentMode := MODE_PRECEDING_BLANK_LINES
	for _, l := range lines {
		switch currentMode {
		case MODE_PRECEDING_BLANK_LINES:
			if !isBlank(l) {
				currentMode = MODE_SIGNIFICANT_LINES
			}
		case MODE_SIGNIFICANT_LINES:
			if isBlank(l) {
				currentMode = MODE_TRAILING_BLANK_LINES
			}
		case MODE_TRAILING_BLANK_LINES:
			if !isBlank(l) {
				blocks = append(blocks, currentBlock)
				currentBlock = nil
				currentMode = MODE_SIGNIFICANT_LINES
			}
		}
		currentBlock = append(currentBlock, l)
	}
	if len(blocks) == 0 && currentMode == MODE_PRECEDING_BLANK_LINES {
		// If the file only contained blank lines, act as if the file was empty altogether.
		return nil
	}
	// Commit the latest ongoing currentBlock.
	return append(blocks, currentBlock)
}

// isBlank checks whether a line is all spaces or tabs.
func isBlank(l Line) bool {
	if len(l.Text) == 0 {
		return true
	}
	for _, c := range l.Text {
		if c != ' ' && c != '\t' {
			return false
		}
	}
	return true
}
