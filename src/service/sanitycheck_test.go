package service

import (
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
}
