package reconciler

import (
	"errors"
	. "github.com/jotaen/klog/src"
	"regexp"
)

// EntryReconciler is for an existing or inserting a new entry into a record.
type EntryReconciler struct {
	Reconciler
	recordPointer uint // `-1` indicates to prepend
}

func NewEntryReconciler(base Reconciler, matchRecord func(Record) bool) *EntryReconciler {
	index := -1
	for i, r := range base.records {
		if matchRecord(r) {
			index = i
			break
		}
	}
	if index == -1 {
		return nil
	}
	return &EntryReconciler{
		Reconciler:    base,
		recordPointer: uint(index),
	}
}

func (r *EntryReconciler) AppendEntry(handler func(Record) string) (*ReconcileResult, error) {
	newEntry := handler(r.records[r.recordPointer])
	lastEntry := lastLine(r.blocks[r.recordPointer].SignificantLines())
	result := insert(
		flatten(r.blocks),
		lastEntry.LineNumber,
		[]InsertableText{{newEntry, 1}},
		stylePreferencesOrDefault(r.blocks[r.recordPointer]),
	)
	return makeResult(result, r.recordPointer)
}

func (r *EntryReconciler) CloseOpenRange(handler func(Record) (Time, EntrySummary)) (*ReconcileResult, error) {
	record := r.records[r.recordPointer]
	openRangeEntryIndex := -1
	for i, e := range record.Entries() {
		e.Unbox(
			func(Range) interface{} { return nil },
			func(Duration) interface{} { return nil },
			func(OpenRange) interface{} {
				openRangeEntryIndex = i
				return nil
			},
		)
	}
	if openRangeEntryIndex == -1 {
		return nil, errors.New("No open time range found")
	}
	time, summary := handler(record)
	lastEntry := lastLine(r.blocks[r.recordPointer].SignificantLines())
	openRangeLineIndex := lastEntry.LineNumber - len(record.Entries()) + openRangeEntryIndex
	allLines := flatten(r.blocks)
	originalText := allLines[openRangeLineIndex].Text
	summaryText := func() string {
		if summary.IsEmpty() {
			return ""
		}
		return " " + summary[0]
	}()
	allLines[openRangeLineIndex].Text = regexp.MustCompile(`^(.*?)\?+(.*)$`).
		ReplaceAllString(originalText, "${1}"+time.ToString()+"${2}"+summaryText)
	return makeResult(allLines, r.recordPointer)
}
