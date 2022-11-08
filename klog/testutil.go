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
func Ɀ_Slashes_(d Date) Date {
	df, canCast := d.(*date)
	if !canCast {
		panic("Operation failed!")
	}
	df.format.UseDashes = false
	return df
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
func Ɀ_EntrySummary_(line ...string) EntrySummary {
	summary, err := NewEntrySummary(line...)
	if err != nil {
		panic("Operation failed!")
	}
	return summary
}

// Deprecated
func Ɀ_ForceSign_(d Duration) Duration {
	do, canCast := d.(*duration)
	if !canCast {
		panic("Operation failed!")
	}
	do.format.ForceSign = true
	return do
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
	tm.format.Use24HourClock = false
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

// Deprecated
func Ɀ_NoSpaces_(r Range) Range {
	tr, canCast := r.(*timeRange)
	if !canCast {
		panic("Operation failed!")
	}
	tr.format.UseSpacesAroundDash = false
	return tr
}

// Deprecated
func Ɀ_NoSpacesO_(r OpenRange) OpenRange {
	or, canCast := r.(*openRange)
	if !canCast {
		panic("Operation failed!")
	}
	or.format.UseSpacesAroundDash = false
	return or
}

// Deprecated
func Ɀ_QuestionMarks_(r OpenRange, additionalQuestionMarks int) OpenRange {
	or, canCast := r.(*openRange)
	if !canCast {
		panic("Operation failed!")
	}
	or.format.AdditionalPlaceholderChars = additionalQuestionMarks
	return or
}
