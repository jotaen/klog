package reconciling

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/parser/txt"
)

// Creator is a function interface for creating a new reconciler.
type Creator func([]klog.Record, []txt.Block) *Reconciler

type AdditionalData struct {
	ShouldTotal klog.ShouldTotal
	Summary     klog.RecordSummary
}

// NewReconcilerForNewRecord is a reconciler creator for a new record at a given date and
// with the given parameters.
func NewReconcilerForNewRecord(atDate Styled[klog.Date], ad AdditionalData) Creator {
	return func(rs []klog.Record, bs []txt.Block) *Reconciler {
		record := klog.NewRecord(atDate.Value)
		if ad.ShouldTotal != nil {
			record.SetShouldTotal(ad.ShouldTotal)
		}
		if ad.Summary != nil {
			record.SetSummary(ad.Summary)
		}
		reconciler := &Reconciler{
			Record:          record,
			recordPointer:   -1,
			lastLinePointer: -1,
			style:           elect(*defaultStyle(), rs, bs),
			lines:           flatten(bs),
		}
		dateFormat := reconciler.style.dateFormat()
		if !atDate.AutoStyle {
			dateFormat = atDate.Value.Format()
		}
		recordText := func() []insertableText {
			result := atDate.Value.ToStringWithFormat(dateFormat)
			if ad.ShouldTotal != nil {
				result += " (" + ad.ShouldTotal.ToString() + ")"
			}
			return []insertableText{{result, 0}}
		}()
		for _, s := range ad.Summary {
			recordText = append(recordText, insertableText{s, 0})
		}
		newRecordLines, insertPointer, lastLineOffset, newRecordIndex := func() ([]insertableText, int, int, int) {
			if len(rs) == 0 {
				return recordText, 0, 1, 0
			}
			i := 0
			for _, r := range rs {
				if i == 0 && !atDate.Value.IsAfterOrEqual(r.Date()) {
					// The new record is dated prior to the first one, so we have to append a blank line.
					recordText = append(recordText, blankLine)
					return recordText, 0, 1, 0
				}
				if len(rs)-1 == i || (atDate.Value.IsAfterOrEqual(r.Date()) && !atDate.Value.IsAfterOrEqual(rs[i+1].Date())) {
					// The record is in between.
					break
				}
				i++
			}
			// The new record is dated after the last one, so we have to prepend a blank line.
			recordText = append([]insertableText{blankLine}, recordText...)
			return recordText, indexOfLastSignificantLine(bs[i]), 2, i + 1
		}()

		// Insert record and adjust pointers accordingly.
		reconciler.insert(insertPointer, newRecordLines)
		reconciler.lastLinePointer = insertPointer + lastLineOffset
		reconciler.recordPointer = newRecordIndex
		return reconciler
	}
}

// NewReconcilerAtRecord is a reconciler creator for an existing record at a given date.
func NewReconcilerAtRecord(atDate klog.Date) Creator {
	return func(rs []klog.Record, bs []txt.Block) *Reconciler {
		index := -1
		for i, r := range rs {
			if r.Date().IsEqualTo(atDate) {
				index = i
				break
			}
		}
		if index == -1 {
			return nil
		}
		style := determine(rs[index], bs[index])
		return &Reconciler{
			Record:          rs[index],
			style:           elect(*style, rs, bs),
			lastLinePointer: indexOfLastSignificantLine(bs[index]),
			recordPointer:   index,
			lines:           flatten(bs),
		}
	}
}

func flatten(bs []txt.Block) []txt.Line {
	var result []txt.Line
	for _, b := range bs {
		result = append(result, b.Lines()...)
	}
	return result
}

func indexOfLastSignificantLine(block txt.Block) int {
	significantLines, precedingInsignificantLineCount, _ := block.SignificantLines()
	return block.OverallLineIndex(precedingInsignificantLineCount + len(significantLines))
}
