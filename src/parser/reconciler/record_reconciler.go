package reconciler

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser/lineparsing"
)

// RecordReconciler is for inserting a new record into a list of records.
type RecordReconciler struct {
	records            []Record
	blocks             []lineparsing.Block
	maybeRecordPointer int
}

func NewRecordReconciler(rs []Record, bs []lineparsing.Block, newDate Date) *RecordReconciler {
	index := -1
	for i, r := range rs {
		if i == 0 && !newDate.IsAfterOrEqual(r.Date()) {
			break
		}
		if i == len(rs)-1 {
			index = len(rs) - 1
			break
		}
		if newDate.IsAfterOrEqual(r.Date()) && !newDate.IsAfterOrEqual(rs[i+1].Date()) {
			index = i
			break
		}
	}
	return &RecordReconciler{
		records:            rs,
		blocks:             bs,
		maybeRecordPointer: index,
	}
}

var blankLine = InsertableText{"", 0}

func (r *RecordReconciler) InsertBlock(texts []InsertableText) (*ReconcileResult, error) {
	lineNumber, newRecordIndex, insertable := func() (int, uint, []InsertableText) {
		if r.maybeRecordPointer == -1 {
			if len(r.records) == 0 {
				return 0, 0, texts
			}
			return 0, 0, append(texts, blankLine)
		}
		lastEntry := lastLine(r.blocks[r.maybeRecordPointer].SignificantLines())
		return lastEntry.LineNumber,
			uint(r.maybeRecordPointer + 1),
			append([]InsertableText{blankLine}, texts...)
	}()
	var styleReferenceBlock lineparsing.Block
	if len(r.blocks) > 0 {
		styleReferenceBlock = r.blocks[0]
	}
	lines := insert(
		flatten(r.blocks),
		lineNumber,
		insertable,
		stylePreferencesOrDefault(styleReferenceBlock),
	)
	return makeResult(lines, newRecordIndex)
}
