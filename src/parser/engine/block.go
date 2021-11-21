package engine

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
	var blocks []Block
	var currentBlock Block
	significantMode := false
	isFirstException := true
	hasSeenSignificantContent := false
	for _, l := range lines {
		shallCommit := false
		if significantMode || isBlank(l) {
			shallCommit = true
			if isBlank(l) {
				significantMode = false
			}
		} else if isFirstException {
			shallCommit = true
			significantMode = true
			isFirstException = false
		}
		if shallCommit {
			if !hasSeenSignificantContent && !isBlank(l) {
				hasSeenSignificantContent = true
			}
			currentBlock = append(currentBlock, l)
			continue
		}
		blocks = append(blocks, currentBlock)
		significantMode = true
		currentBlock = Block{l}
	}
	if currentBlock != nil {
		blocks = append(blocks, currentBlock)
	}
	if len(blocks) == 1 && !hasSeenSignificantContent {
		return nil
	}
	return blocks
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
