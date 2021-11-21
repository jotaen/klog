package reconciling

import (
	"errors"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/lineparsing"
)

type Handler func(Reconciler) (*Result, error)

// Chain tries to apply multiple reconcilers one after the other. It returns the result
// of the first successful one.
func Chain(base Reconciler, handler ...Handler) (*Result, error) {
	for i, reconcile := range handler {
		result, err := reconcile(base)
		if err == nil && result != nil {
			return result, nil
		}
		_, isNotEligibleError := err.(NotEligibleError)
		if isNotEligibleError && i < len(handler)-1 {
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

func makeResult(ls []lineparsing.Line, recordIndex uint) (*Result, error) {
	newText := join(ls)
	newRecords, _, pErr := parser.Parse(newText)
	if pErr != nil {
		err := pErr.Get()[0]
		return nil, errors.New(err.Message())
	}
	return &Result{
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
