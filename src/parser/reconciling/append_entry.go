package reconciling

// AppendEntry adds a new entry to the end of the record.
func (r *Reconciler) AppendEntry(newEntry string) (*Result, error) {
	r.insert(r.lastLinePointer, []insertableText{{newEntry, 1}})
	return r.MakeResult()
}
