package reconciling

import "github.com/jotaen/klog/klog"

// AppendEntry adds a new entry to the end of the record.
// `newEntry` must include the entry value at the beginning of its first line.
func (r *Reconciler) AppendEntry(newEntry klog.EntrySummary) error {
	r.insert(r.lastLinePointer, toMultilineEntryTexts("", newEntry))
	return nil
}
