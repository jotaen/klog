package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "klog"
	"testing"
	gotime "time"
)

func TestNoWarningWhenAllGood(t *testing.T) {
	timestamp := gotime.Now()
	today := NewDateFromTime(timestamp)
	now := NewTimeFromTime(timestamp)
	rs := []Record{
		func() Record {
			// OK: Record in the future but without entries
			r := NewRecord(today.PlusDays(1))
			return r
		}(), func() Record {
			// OK: Open range today
			r := NewRecord(today)
			r.StartOpenRange(now, "")
			return r
		}(), func() Record {
			// OK: Just a regular record in the past
			r := NewRecord(today.PlusDays(-1))
			r.AddDuration(NewDuration(1, 2), "")
			return r
		}(),
	}
	ws := SanityCheck(timestamp, rs)
	require.Nil(t, ws)
}

func TestNoOpenRangeWarningWhenYesterdayAndNoRecordToday(t *testing.T) {
	timestamp := gotime.Now()
	today := NewDateFromTime(timestamp)
	now := NewTimeFromTime(timestamp)
	rs := []Record{
		func() Record {
			// This open range is okay, because there is no record at today’s date
			r := NewRecord(today.PlusDays(-1))
			r.StartOpenRange(now, "")
			return r
		}(), func() Record {
			r := NewRecord(today.PlusDays(2))
			return r
		}(),
	}
	ws := SanityCheck(timestamp, rs)
	require.Nil(t, ws)
}

func TestOpenRangeWarningWhenUnclosedOpenRangeBeforeTodayRegardlessOfOrder(t *testing.T) {
	timestamp := gotime.Now()
	today := NewDateFromTime(timestamp)
	now := NewTimeFromTime(timestamp)
	// The warnings must work reliably even when the records are not ordered by date initially
	rs := []Record{
		func() Record {
			// NOT OK: There is a record at today’s date
			r := NewRecord(today.PlusDays(-1))
			r.StartOpenRange(now, "")
			return r
		}(), func() Record {
			r := NewRecord(today)
			return r
		}(), func() Record {
			// NOT OK: There is a record at today’s date
			r := NewRecord(today.PlusDays(-2))
			r.StartOpenRange(now, "")
			return r
		}(),
	}
	ws := SanityCheck(timestamp, rs)
	require.NotNil(t, ws)
	require.Len(t, ws, 2)
	assert.Equal(t, today.PlusDays(-1), ws[0].Date)
	assert.Equal(t, today.PlusDays(-2), ws[1].Date)
	for _, w := range ws {
		assert.Equal(t, "Unclosed open range", w.Message)
	}
}

func TestFutureEntriesWarning(t *testing.T) {
	timestamp := gotime.Now()
	today := NewDateFromTime(timestamp)
	rs := []Record{
		func() Record {
			r := NewRecord(today.PlusDays(1))
			r.AddDuration(NewDuration(2, 0), "")
			return r
		}(), func() Record {
			r := NewRecord(today)
			return r
		}(),
	}
	ws := SanityCheck(timestamp, rs)
	require.NotNil(t, ws)
	require.Len(t, ws, 1)
	assert.Equal(t, today.PlusDays(1), ws[0].Date)
	for _, w := range ws {
		assert.Equal(t, "Entry in future record", w.Message)
	}
}

func TestMoreThan24HoursPerRecord(t *testing.T) {
	timestamp := gotime.Now()
	today := NewDateFromTime(timestamp)
	rs := []Record{
		func() Record {
			r := NewRecord(today.PlusDays(-1))
			r.AddDuration(NewDuration(24, 1), "")
			return r
		}(), func() Record {
			r := NewRecord(today)
			r.AddRange(Ɀ_Range_(Ɀ_Time_(0, 0), Ɀ_Time_(23, 0)), "")
			r.AddDuration(NewDuration(2, 0), "")
			return r
		}(),
		func() Record {
			r := NewRecord(today.PlusDays(-3))
			r.AddDuration(NewDuration(24, 0), "")
			return r
		}(),
	}
	ws := SanityCheck(timestamp, rs)
	require.NotNil(t, ws)
	require.Len(t, ws, 2)
	assert.Equal(t, today, ws[0].Date)
	assert.Equal(t, today.PlusDays(-1), ws[1].Date)
	for _, w := range ws {
		assert.Equal(t, "Total time exceeds 24 hours", w.Message)
	}
}

func TestOverlappingTimeRanges(t *testing.T) {
	timestamp := gotime.Now()
	today := NewDateFromTime(timestamp)
	rs := []Record{
		func() Record {
			// No overlap
			r := NewRecord(today)
			r.AddDuration(NewDuration(5, 0), "")
			r.AddRange(Ɀ_Range_(Ɀ_Time_(4, 0), Ɀ_Time_(4, 59)), "")
			r.AddRange(Ɀ_Range_(Ɀ_Time_(0, 0), Ɀ_Time_(2, 0)), "")
			r.AddRange(Ɀ_Range_(Ɀ_Time_(2, 0), Ɀ_Time_(4, 0)), "")
			r.AddRange(Ɀ_Range_(Ɀ_Time_(4, 0), Ɀ_Time_(4, 0)), "") // point in time range
			r.AddRange(Ɀ_Range_(Ɀ_Time_(5, 0), Ɀ_Time_(6, 0)), "")
			r.StartOpenRange(Ɀ_Time_(0, 44), "")
			return r
		}(), func() Record {
			// Overlap with sorted entries
			r := NewRecord(today.PlusDays(-1))
			r.AddRange(Ɀ_Range_(Ɀ_Time_(0, 30), Ɀ_Time_(1, 0)), "")
			r.AddRange(Ɀ_Range_(Ɀ_Time_(2, 0), Ɀ_Time_(5, 0)), "")
			r.AddRange(Ɀ_Range_(Ɀ_Time_(4, 59), Ɀ_Time_(6, 0)), "")
			r.AddRange(Ɀ_Range_(Ɀ_Time_(18, 30), Ɀ_Time_(19, 0)), "")
			return r
		}(), func() Record {
			// Overlap with unsorted entries
			r := NewRecord(today.PlusDays(-2))
			r.AddRange(Ɀ_Range_(Ɀ_Time_(0, 30), Ɀ_Time_(0, 45)), "")
			r.AddRange(Ɀ_Range_(Ɀ_Time_(2, 45), Ɀ_Time_(3, 45)), "")
			r.AddRange(Ɀ_Range_(Ɀ_TimeYesterday_(23, 0), Ɀ_Time_(1, 0)), "")
			return r
		}(),
	}
	ws := SanityCheck(timestamp, rs)
	require.NotNil(t, ws)
	require.Len(t, ws, 2)
	assert.Equal(t, today.PlusDays(-1), ws[0].Date)
	assert.Equal(t, today.PlusDays(-2), ws[1].Date)
	for _, w := range ws {
		assert.Equal(t, "Overlapping time ranges", w.Message)
	}
}
