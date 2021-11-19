package klog

// Entry is a time value and a summary.
// A time value can be a Range, a Duration, or an OpenRange.
type Entry struct {
	value   interface{}
	summary EntrySummary
}

func NewEntry(value interface{}, summary EntrySummary) Entry {
	return Entry{value, summary}
}

func (e *Entry) Summary() EntrySummary {
	return e.summary
}

func (e *Entry) Unbox(r func(Range) interface{}, d func(Duration) interface{}, o func(OpenRange) interface{}) interface{} {
	switch x := e.value.(type) {
	case Range:
		return r(x)
	case Duration:
		return d(x)
	case OpenRange:
		return o(x)
	}
	panic("Incomplete switch statement")
}

func (e *Entry) Duration() Duration {
	return (e.Unbox(
		func(r Range) interface{} { return r.Duration() },
		func(d Duration) interface{} { return d },
		func(o OpenRange) interface{} { return NewDuration(0, 0) },
	)).(Duration)
}
