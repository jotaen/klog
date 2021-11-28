package reconciling

import (
	"errors"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/engine"
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

func makeResult(ls []engine.Line, recordIndex uint) (*Result, error) {
	newText := join(ls)
	newRecords, _, pErrs := parser.Parse(newText)
	if pErrs != nil {
		// This is just a safe guard mechanism. If it happens, then there is a bug
		// in the calling reconciler method.
		return nil, errors.New("This operation wouldn’t result in a valid record")
	}
	return &Result{
		newRecords[recordIndex],
		newText,
	}, nil
}

// insert inserts some new lines into a text at a specific line number (position).
func insert(ls []engine.Line, position int, texts []InsertableText, stylePrefs stylePreferences) []engine.Line {
	if position > len(ls)+1 {
		panic("Out of bounds")
	}
	result := make([]engine.Line, len(ls)+len(texts))
	offset := 0
	for i := range result {
		if i >= position && offset < len(texts) {
			line := ""
			if texts[offset].Indentation > 0 {
				line += stylePrefs.indentationStyle
			}
			line += texts[offset].Text + stylePrefs.lineEndingStyle
			result[i] = engine.NewLineFromString(line, -999)
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

func join(ls []engine.Line) string {
	result := ""
	for _, l := range ls {
		result += l.Original()
	}
	return result
}

func flatten(blocks []engine.Block) []engine.Line {
	var result []engine.Line
	for _, bs := range blocks {
		result = append(result, bs...)
	}
	return result
}

func lastLine(ls []engine.Line) engine.Line {
	return ls[len(ls)-1]
}
