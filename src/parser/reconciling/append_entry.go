package reconciling

// AppendEntry adds a new entry to the end of the record.
func (r *Reconciler) AppendEntry(newEntry string) (*Result, error) {
	r.insert(r.lastLinePointer, toMultilineEntryTexts("", newEntry))
	return r.MakeResult()
}
