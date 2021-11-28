package klog

/**
Only use these functions in test code.
(They cannot live in a `_test.go` file
because they need to be imported elsewhere.
They cannot live in a separate package
neither due to circular imports.)
The `Deprecated` markers and the funny naming
are supposed to act as a reminder for this.
*/

// Deprecated
func Ɀ_Date_(year int, month int, day int) Date {
	date, err := NewDate(year, month, day)
	if err != nil {
		panic("Operation failed!")
	}
	return date
}

// Deprecated
func Ɀ_RecordSummary_(line ...string) RecordSummary {
	summary, err := NewRecordSummary(line...)
	if err != nil {
		panic("Operation failed!")
	}
	return summary
}

// Deprecated
func Ɀ_Time_(hour int, minute int) Time {
	time, err := NewTime(hour, minute)
	if err != nil {
		panic("Operation failed!")
	}
	return time
}

// Deprecated
func Ɀ_IsAmPm_(t Time) Time {
	tm, canCast := t.(*time)
	if !canCast {
		panic("Operation failed!")
	}
	tm.is24HourClock = false
	return tm
}

// Deprecated
func Ɀ_TimeYesterday_(hour int, minute int) Time {
	time, err := NewTimeYesterday(hour, minute)
	if err != nil {
		panic("Operation failed!")
	}
	return time
}

// Deprecated
func Ɀ_TimeTomorrow_(hour int, minute int) Time {
	time, err := NewTimeTomorrow(hour, minute)
	if err != nil {
		panic("Operation failed!")
	}
	return time
}

// Deprecated
func Ɀ_Range_(start Time, end Time) Range {
	r, err := NewRange(start, end)
	if err != nil {
		panic("Operation failed!")
	}
	return r
}
