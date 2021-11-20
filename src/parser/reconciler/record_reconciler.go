package reconciler

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
	"regexp"
)

// RecordReconciler is for inserting a new entry into a record.
type RecordReconciler struct {
	records       []parser.ParsedRecord
	recordPointer uint // `-1` indicates to prepend
}

func NewRecordReconciler(rs []parser.ParsedRecord, matchRecord func(Record) bool) *RecordReconciler {
	index := -1
	for i, r := range rs {
		if matchRecord(r) {
			index = i
			break
		}
	}
	if index == -1 {
		return nil
	}
	return &RecordReconciler{
		records:       rs,
		recordPointer: uint(index),
	}
}

func (r *RecordReconciler) AppendEntry(handler func(Record) string) (*ReconcileResult, error) {
	newEntry := handler(r.records[r.recordPointer])
	result := insert(
		decompose(r.records),
		r.pr.LastLineOfRecord[r.recordPointer],
		[]InsertableText{{newEntry, 1}},
		r.records[r.recordPointer],
	)
	return makeResult(result, r.recordPointer)
}

func (r *RecordReconciler) CloseOpenRange(handler func(Record) (Time, EntrySummary)) (*ReconcileResult, error) {
	record := r.records[r.recordPointer]
	if record.OpenRange() == nil {
		return nil, errors.New("No open time range found")
	}
	entryIndex := 0
	for i, e := range record.Entries() {
		e.Unbox(
			func(Range) interface{} { return nil },
			func(Duration) interface{} { return nil },
			func(OpenRange) interface{} {
				entryIndex = i
				return nil
			},
		)
	}
	time, summary := handler(record)
	openRangeLineIndex := r.pr.LastLineOfRecord[r.recordPointer] - len(record.Entries()) + entryIndex
	originalText := r.pr.Lines[openRangeLineIndex].Text
	summaryText := func() string {
		if summary.IsEmpty() {
			return ""
		}
		return " " + summary[0]
	}()
	r.pr.Lines[openRangeLineIndex].Text = regexp.MustCompile(`^(.*?)\?+(.*)$`).
		ReplaceAllString(originalText, "${1}"+time.ToString()+"${2}"+summaryText)
	return makeResult(decompose(r.records), r.recordPointer)
}
