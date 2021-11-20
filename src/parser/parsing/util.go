package parsing

// SubRune returns a subset of a rune list. It might be shorter than
// the requested length, if the text doesn’t contain enough characters.
// It returns empty, if the start position is bigger than the length.
func SubRune(text []rune, start int, length int) []rune {
	if start >= len(text) {
		return nil
	}
	if start+length > len(text) {
		length = len(text) - start
	}
	return text[start : start+length]
}

// IsWhitespace checks whether a rune is a space or a tab.
func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}

// GroupIntoBlocks splits a list of Line’s on the blank lines in between them.
// It’s basically like paragraphs in a text.
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

// IsBlank checks whether a line is all spaces or tabs.
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

type Text struct {
	Text        string
	Indentation int
}

// Insert inserts some new lines into a text at a specific line number (position).
func Insert(ls []Line, position int, texts []Text, prefs Preferences) []Line {
	if position > len(ls)+1 {
		panic("Out of bounds")
	}
	result := make([]Line, len(ls)+len(texts))
	offset := 0
	for i := range result {
		if i >= position && offset < len(texts) {
			line := ""
			if texts[offset].Indentation > 0 {
				line += prefs.IndentationStyle
			}
			line += texts[offset].Text + prefs.LineEnding
			result[i] = NewLineFromString(line, -999)
			offset++
		} else {
			result[i] = ls[i-offset]
		}
		result[i].LineNumber = i + 1
	}
	if position > 0 && result[position-1].LineEnding == "" {
		result[position-1].LineEnding = prefs.LineEnding
	}
	return result
}
