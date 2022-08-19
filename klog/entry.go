package klog

// Entry is a time value and an associated entry summary.
// A time value can be a Range, a Duration, or an OpenRange.
type Entry struct {
	value   any
	summary EntrySummary
}

func NewEntryFromDuration(value Duration, summary EntrySummary) Entry {
	return Entry{value, summary}
}

func NewEntryFromRange(value Range, summary EntrySummary) Entry {
	return Entry{value, summary}
}

func NewEntryFromOpenRange(value OpenRange, summary EntrySummary) Entry {
	return Entry{value, summary}
}

func (e *Entry) Summary() EntrySummary {
	return e.summary
}

// Unbox converts the underlying time value.
func Unbox[TargetT any](e *Entry, r func(Range) TargetT, d func(Duration) TargetT, o func(OpenRange) TargetT) TargetT {
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

// Duration returns the duration value of the underlying time value.
func (e *Entry) Duration() Duration {
	return Unbox[Duration](e,
		func(r Range) Duration { return r.Duration() },
		func(d Duration) Duration { return d },
		func(o OpenRange) Duration { return NewDuration(0, 0) },
	)
}
