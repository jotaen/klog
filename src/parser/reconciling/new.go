package reconciling

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/engine"
)

type Creator func(parsedRecords []parser.ParsedRecord) *Reconciler

func NewReconcilerAtNewRecord(parsedRecords []parser.ParsedRecord, newDate Date, shouldTotal ShouldTotal) *Reconciler {
	insertPointer, newRecordIndex, shallPrepend := func() (int, int, bool) {
		if len(parsedRecords) == 0 {
			return 0, 0, true
		}
		for i, r := range parsedRecords {
			if i == 0 && !newDate.IsAfterOrEqual(r.Date()) {
				// The new record is dated prior to the first one.
				return 0, 0, true
			}
			if len(parsedRecords)-1 == i || (newDate.IsAfterOrEqual(r.Date()) && !newDate.IsAfterOrEqual(parsedRecords[i+1].Date())) {
				return lastLine(parsedRecords[i].Block.SignificantLines()).LineNumber, i, false
			}
		}
		// The new record is dated after the last one.
		return lastLine(parsedRecords[len(parsedRecords)-1].Block.SignificantLines()).LineNumber, len(parsedRecords), false
	}()
	style := parser.DefaultStyle()
	if len(parsedRecords) > 0 {
		style = parsedRecords[len(parsedRecords)-1].Style
	}
	reconciler := &Reconciler{
		record:          NewRecord(newDate),
		recordPointer:   newRecordIndex,
		lastLinePointer: -1,
		style:           style,
		lines:           flatten(parsedRecords),
	}
	headline := func() insertableText {
		result := newDate.ToStringWithFormat(DateFormat{UseDashes: reconciler.style.UsesDashesInDate})
		if shouldTotal != nil {
			result += " (" + shouldTotal.ToString() + ")"
		}
		return insertableText{result, 0}
	}()
	newRecordLines, lastLineOffset := func() ([]insertableText, int) {
		if len(parsedRecords) == 0 {
			return []insertableText{headline}, 1
		}
		if shallPrepend {
			return []insertableText{headline, blankLine}, 1
		}
		return []insertableText{blankLine, headline}, 2
	}()
	reconciler.lastLinePointer = insertPointer + lastLineOffset
	reconciler.insert(insertPointer, newRecordLines)
	return reconciler
}

func NewReconcilerAtRecord(parsedRecords []parser.ParsedRecord, atDate Date) *Reconciler {
	index := -1
	for i, r := range parsedRecords {
		if r.Date().IsEqualTo(atDate) {
			index = i
			break
		}
	}
	if index == -1 {
		return nil
	}
	return &Reconciler{
		record:          parsedRecords[index],
		style:           parsedRecords[index].Style,
		lastLinePointer: lastLine(parsedRecords[index].Block.SignificantLines()).LineNumber,
		recordPointer:   index,
		lines:           flatten(parsedRecords),
	}
}

func flatten(parsedRecords []parser.ParsedRecord) []engine.Line {
	var result []engine.Line
	for _, r := range parsedRecords {
		result = append(result, r.Block...)
	}
	return result
}

func lastLine(block engine.Block) engine.Line {
	return block[len(block)-1]
}
