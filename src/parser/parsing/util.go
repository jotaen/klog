package parsing

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

func IsNewline(r rune) bool {
	return r == '\n' || r == '\r'
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
