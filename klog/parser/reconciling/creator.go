package reconciling

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/parser"
	"github.com/jotaen/klog/klog/parser/txt"
)

// Creator is a function interface for creating a new reconciler.
type Creator func(parsedRecords []parser.ParsedRecord) *Reconciler

type RecordParams struct {
	Date        klog.Date
	ShouldTotal klog.ShouldTotal
	Summary     klog.RecordSummary
}

// NewReconcilerForNewRecord creates a reconciler for a new record at a given date and
// with the given parameters.
func NewReconcilerForNewRecord(parsedRecords []parser.ParsedRecord, params RecordParams) *Reconciler {
	record := klog.NewRecord(params.Date)
	if params.ShouldTotal != nil {
		record.SetShouldTotal(params.ShouldTotal)
	}
	if params.Summary != nil {
		record.SetSummary(params.Summary)
	}
	reconciler := &Reconciler{
		record:          record,
		recordPointer:   -1,
		lastLinePointer: -1,
		style:           parser.Elect(*parser.DefaultStyle(), parsedRecords),
		lines:           flatten(parsedRecords),
	}
	recordText := func() []insertableText {
		result := params.Date.ToStringWithFormat(reconciler.style.DateFormat.Get())
		if params.ShouldTotal != nil {
			result += " (" + params.ShouldTotal.ToString() + ")"
		}
		return []insertableText{{result, 0}}
	}()
	for _, s := range params.Summary {
		recordText = append(recordText, insertableText{s, 0})
	}
	newRecordLines, insertPointer, lastLineOffset, newRecordIndex := func() ([]insertableText, int, int, int) {
		if len(parsedRecords) == 0 {
			return recordText, 0, 1, 0
		}
		i := 0
		for _, r := range parsedRecords {
			if i == 0 && !params.Date.IsAfterOrEqual(r.Date()) {
				// The new record is dated prior to the first one, so we have to append a blank line.
				recordText = append(recordText, blankLine)
				return recordText, 0, 1, 0
			}
			if len(parsedRecords)-1 == i || (params.Date.IsAfterOrEqual(r.Date()) && !params.Date.IsAfterOrEqual(parsedRecords[i+1].Date())) {
				// The record is in between.
				break
			}
			i++
		}
		// The new record is dated after the last one, so we have to prepend a blank line.
		recordText = append([]insertableText{blankLine}, recordText...)
		return recordText, lastLine(parsedRecords[i].Block.SignificantLines()).LineNumber, 2, i + 1
	}()

	// Insert record and adjust pointers accordingly.
	reconciler.insert(insertPointer, newRecordLines)
	reconciler.lastLinePointer = insertPointer + lastLineOffset
	reconciler.recordPointer = newRecordIndex
	return reconciler
}

// NewReconcilerAtRecord creates a reconciler for an existing record at a given date.
func NewReconcilerAtRecord(parsedRecords []parser.ParsedRecord, atDate klog.Date) *Reconciler {
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

func flatten(parsedRecords []parser.ParsedRecord) []txt.Line {
	var result []txt.Line
	for _, r := range parsedRecords {
		result = append(result, r.Block...)
	}
	return result
}

func lastLine(block txt.Block) txt.Line {
	return block[len(block)-1]
}
