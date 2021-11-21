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
// Deprecated. (Use multiple return value instead!)
type ReconcileResult struct {
	NewRecord Record
	NewText   string
}

type Reconcile func(records []Record, blocks []lineparsing.Block) (*ReconcileResult, error)

// InsertableText is for inserting lines of text into a list of Line’s,
// without needing to know anything about indentation or line ending style.
type InsertableText struct {
	Text        string
	Indentation int
}

// NotEligibleError is for Chain to indicate that it should proceed with the next reconciler.
type NotEligibleError struct{}

func (e NotEligibleError) Error() string { return "Boom" } // TODO

// Chain tries to apply multiple reconcilers one after the other. It returns the result
// of the first successful one.
func Chain(records []Record, blocks []lineparsing.Block, reconcilers ...Reconcile) (*ReconcileResult, error) {
	for i, reconcile := range reconcilers {
		result, err := reconcile(records, blocks)
		if err == nil && result != nil {
			return result, nil
		}
		_, isNotEligibleError := err.(NotEligibleError)
		if isNotEligibleError && i < len(reconcilers)-1 {
			// Try next reconcile function
			continue
		}
		return nil, err
	}
	return nil, NotEligibleError{}
}

type stylePreferences struct {
	indentationStyle string
	lineEndingStyle  string
}

func stylePreferencesOrDefault(b lineparsing.Block) stylePreferences {
	defaultPrefs := stylePreferences{
		indentationStyle: "    ",
		lineEndingStyle:  "\n",
	}
	if b == nil {
		return defaultPrefs
	}
	for _, l := range b.SignificantLines() {
		if len(l.LineEnding) > 0 {
			defaultPrefs.lineEndingStyle = l.LineEnding
		}
		if len(l.PrecedingWhitespace) > 0 {
			defaultPrefs.indentationStyle = l.PrecedingWhitespace
		}
	}
	return defaultPrefs
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
func insert(ls []lineparsing.Line, position int, texts []InsertableText, stylePrefs stylePreferences) []lineparsing.Line {
	if position > len(ls)+1 {
		panic("Out of bounds")
	}
	result := make([]lineparsing.Line, len(ls)+len(texts))
	offset := 0
	for i := range result {
		if i >= position && offset < len(texts) {
			line := ""
			if texts[offset].Indentation > 0 {
				line += stylePrefs.indentationStyle
			}
			line += texts[offset].Text + stylePrefs.lineEndingStyle
			result[i] = lineparsing.NewLineFromString(line, -999)
			offset++
		} else {
			result[i] = ls[i-offset]
		}
		result[i].LineNumber = i + 1
	}
	if position > 0 && result[position-1].LineEnding == "" {
		result[position-1].LineEnding = stylePrefs.lineEndingStyle
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

func flatten(blocks []lineparsing.Block) []lineparsing.Line {
	var result []lineparsing.Line
	for _, bs := range blocks {
		result = append(result, bs...)
	}
	return result
}

func lastLine(ls []lineparsing.Line) lineparsing.Line {
	return ls[len(ls)-1]
}
