package reconciler

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
)

// BlockReconciler is for inserting a new record into a list of records.
type BlockReconciler struct {
	records            []parser.ParsedRecord
	maybeRecordPointer int
}

func NewBlockReconciler(rs []parser.ParsedRecord, newDate Date) *BlockReconciler {
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
	return &BlockReconciler{
		records:            rs,
		maybeRecordPointer: index,
	}
}

var blankLine = InsertableText{"", 0}

func (r *BlockReconciler) InsertBlock(texts []InsertableText) (*ReconcileResult, error) {
	lineIndex, newRecordIndex, insertable := func() (int, uint, []InsertableText) {
		if r.maybeRecordPointer == -1 {
			if len(r.records) == 0 {
				return 0, 0, texts
			}
			return 0, 0, append(texts, blankLine)
		}
		return r.pr.LastLineOfRecord[r.maybeRecordPointer],
			uint(r.maybeRecordPointer + 1),
			append([]InsertableText{blankLine}, texts...)
	}()
	lines := insert(
		decompose(r.records),
		lineIndex,
		insertable,
		r.pr.Preferences,
	)
	return makeResult(lines, newRecordIndex)
}
