package reconciling

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/engine"
)

// Creator is a function interface for creating a new reconciler.
type Creator func(parsedRecords []parser.ParsedRecord) *Reconciler

// NewReconcilerForNewRecord creates a reconciler for a new record at a given date.
func NewReconcilerForNewRecord(parsedRecords []parser.ParsedRecord, newDate Date, shouldTotal ShouldTotal) *Reconciler {
	record := NewRecord(newDate)
	if shouldTotal != nil {
		record.SetShouldTotal(shouldTotal)
	}
	reconciler := &Reconciler{
		record:          record,
		recordPointer:   -1,
		lastLinePointer: -1,
		style:           parser.Elect(*parser.DefaultStyle(), parsedRecords),
		lines:           flatten(parsedRecords),
	}
	headline := func() insertableText {
		result := newDate.ToStringWithFormat(reconciler.style.DateFormat())
		if shouldTotal != nil {
			result += " (" + shouldTotal.ToString() + ")"
		}
		return insertableText{result, 0}
	}()
	newRecordLines, insertPointer, lastLineOffset, newRecordIndex := func() ([]insertableText, int, int, int) {
		if len(parsedRecords) == 0 {
			return []insertableText{headline}, 0, 1, 0
		}
		i := 0
		for _, r := range parsedRecords {
			if i == 0 && !newDate.IsAfterOrEqual(r.Date()) {
				// The new record is dated prior to the first one.
				return []insertableText{headline, blankLine}, 0, 1, 0
			}
			if len(parsedRecords)-1 == i || (newDate.IsAfterOrEqual(r.Date()) && !newDate.IsAfterOrEqual(parsedRecords[i+1].Date())) {
				// The record is in between.
				break
			}
			i++
		}
		// The new record is dated after the last one.
		return []insertableText{blankLine, headline}, lastLine(parsedRecords[i].Block.SignificantLines()).LineNumber, 2, i + 1
	}()

	// Insert record and adjust pointers accordingly.
	reconciler.insert(insertPointer, newRecordLines)
	reconciler.lastLinePointer = insertPointer + lastLineOffset
	reconciler.recordPointer = newRecordIndex
	return reconciler
}

// NewReconcilerAtRecord creates a reconciler for an existing record at a given date.
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
		style:           parser.Elect(*parsedRecords[index].Style, parsedRecords),
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
