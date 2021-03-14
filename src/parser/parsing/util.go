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

func Join(ls []Line) string {
	result := ""
	for _, l := range ls {
		result += l.Original()
	}
	return result
}

func IsBlank(l Line) bool {
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

func Insert(ls []Line, position int, text string, isIndented bool, prefs Preferences) []Line {
	if position > len(ls)+1 {
		panic("Out of bounds")
	}
	result := make([]Line, len(ls)+1)
	offset := 0
	for i := range result {
		if i == position {
			line := ""
			if isIndented {
				line += prefs.Indentation
			}
			line += text + prefs.LineEnding
			result[i] = NewLineFromString(line, i+1)
			offset = 1
		} else {
			result[i] = ls[i-offset]
		}
		result[i].LineNumber = i + 1
	}
	if position > 0 && result[position-1].originalLineEnding == "" {
		result[position-1].originalLineEnding = prefs.LineEnding
	}
	return result
}
