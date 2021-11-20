package reconciler

import "github.com/jotaen/klog/src/parser/parsing"

type Text struct {
	Text        string
	Indentation int
}

// Insert inserts some new lines into a text at a specific line number (position).
func Insert(ls []parsing.Line, position int, texts []Text, prefs parsing.Preferences) []parsing.Line {
	if position > len(ls)+1 {
		panic("Out of bounds")
	}
	result := make([]parsing.Line, len(ls)+len(texts))
	offset := 0
	for i := range result {
		if i >= position && offset < len(texts) {
			line := ""
			if texts[offset].Indentation > 0 {
				line += prefs.IndentationStyle
			}
			line += texts[offset].Text + prefs.LineEnding
			result[i] = parsing.NewLineFromString(line, -999)
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
