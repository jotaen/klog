/*
Package reconciler contains logic to manipulate klog source text.
The idea of reconcilers in general is to add or modify serialised records
in a minimally invasive manner. Instead or re-serialising the record itself,
it tries to find the location in the original text and modify that directly.
While this approach might feel a little hacky, it avoids lots of other
complications and sources of bugs, that could potentially mess up user data.
*/
package reconciler

import (
	"errors"
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/lineparsing"
)

// ReconcileResult is the return value of reconcilers and contains
// the modified record as Record and serialised.
type ReconcileResult struct {
	NewRecord Record
	NewText   string
}

// InsertableText is for inserting lines of text into a list of Line’s,
// without needing to know anything about indentation or line ending style.
type InsertableText struct {
	Text        string
	Indentation int
}

type stylePreferences struct {
	indentationStyle string
	lineEndingStyle  string
}

func newDefaultStylePreferences() stylePreferences {
	return stylePreferences{
		indentationStyle: "    ",
		lineEndingStyle:  "\n",
	}
}

func makeResult(ls []lineparsing.Line, recordIndex uint) (*ReconcileResult, error) {
	newText := join(ls)
	newRecords, _, pErr := parser.Parse(newText)
	if pErr != nil {
		err := pErr.Get()[0]
		return nil, errors.New(err.Message())
	}
	return &ReconcileResult{
		newRecords[recordIndex],
		newText,
	}, nil
}

// insert inserts some new lines into a text at a specific line number (position).
func insert(ls []lineparsing.Line, position int, texts []InsertableText, stylePrefs parser.StylePreferences) []lineparsing.Line {
	if position > len(ls)+1 {
		panic("Out of bounds")
	}
	result := make([]lineparsing.Line, len(ls)+len(texts))
	offset := 0
	for i := range result {
		if i >= position && offset < len(texts) {
			line := ""
			if texts[offset].Indentation > 0 {
				line += stylePrefs.IndentationStyle()
			}
			line += texts[offset].Text + stylePrefs.LineEndingStyle()
			result[i] = lineparsing.NewLineFromString(line, -999)
			offset++
		} else {
			result[i] = ls[i-offset]
		}
		result[i].LineNumber = i + 1
	}
	if position > 0 && result[position-1].LineEnding == "" {
		result[position-1].LineEnding = stylePrefs.LineEndingStyle()
	}
	return result
}

func join(ls []lineparsing.Line) string {
	result := ""
	for _, l := range ls {
		result += l.Original()
	}
	return result
}

func decompose(records []parser.ParsedRecord) []lineparsing.Line {
	var result []lineparsing.Line
	for _, r := range records {
		result = append(result, r.Block()...)
	}
	return result
}
