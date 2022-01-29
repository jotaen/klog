package service

import (
	. "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"testing"
	gotime "time"
)

func CountWarningsOfKind(c checker, ws []Warning) int {
	count := 0
	for _, w := range ws {
		if w.Warning() == c.Message() {
			count++
		}
	}
	return count
}

func TestNoWarnForOpenRanges(t *testing.T) {
	timestamp := gotime.Date(2000, 3, 5, 12, 00, 0, 0, gotime.Local)
	today := NewDateFromGo(timestamp)
	now := NewTimeFromGo(timestamp)

	rs1 := []Record{
		func() Record {
			// This open range is okay, because there is no record at today’s date
			r := NewRecord(today.PlusDays(-1))
			r.StartOpenRange(now, nil)
			return r
		}(), func() Record {
			r := NewRecord(today.PlusDays(2))
			return r
		}(),
	}
	ws1 := CheckForWarnings(timestamp, rs1)
	assert.Equal(t, 0, CountWarningsOfKind(&unclosedOpenRangeChecker{}, ws1))

	rs2 := []Record{
		func() Record {
			r := NewRecord(today)
			r.StartOpenRange(now, nil)
			return r
		}(),
	}
	ws2 := CheckForWarnings(timestamp, rs2)
	assert.Equal(t, 0, CountWarningsOfKind(&unclosedOpenRangeChecker{}, ws2))
}

func TestOpenRangeWarningWhenUnclosedOpenRangeBeforeTodayRegardlessOfOrder(t *testing.T) {
	timestamp := gotime.Date(2000, 3, 5, 12, 00, 0, 0, gotime.Local)
	today := NewDateFromGo(timestamp)
	now := NewTimeFromGo(timestamp)
	// The warnings must work reliably even when the records are not ordered by date initially
	rs := []Record{
		func() Record {
			// NOT OK: There is a record at today’s date
			r := NewRecord(today.PlusDays(-1))
			r.StartOpenRange(now, nil)
			return r
		}(), func() Record {
			r := NewRecord(today)
			return r
		}(), func() Record {
			// NOT OK: There is a record at today’s date
			r := NewRecord(today.PlusDays(-2))
			r.StartOpenRange(now, nil)
			return r
		}(),
	}
	ws := CheckForWarnings(timestamp, rs)
	assert.Equal(t, 2, CountWarningsOfKind(&unclosedOpenRangeChecker{}, ws))
	assert.Equal(t, today.PlusDays(-1), ws[0].Date())
	assert.Equal(t, today.PlusDays(-2), ws[1].Date())
}

func TestNoWarningForFutureEntries(t *testing.T) {
	timestamp := gotime.Date(2000, 3, 5, 12, 00, 0, 0, gotime.Local)
	today := NewDateFromGo(timestamp)
	rs := []Record{
		func() Record {
			// Future entry okay if it doesn’t contain entries
			r := NewRecord(today.PlusDays(1))
			return r
		}(),
		func() Record {
			r := NewRecord(today.PlusDays(-1))
			r.AddRange(Ɀ_Range_(Ɀ_Time_(0, 0), Ɀ_Time_(12, 30)), nil)
			r.AddDuration(NewDuration(2, 0), nil)
			return r
		}(),
		func() Record {
			r := NewRecord(today.PlusDays(1))
			// Times shifted to yesterday
			r.AddRange(Ɀ_Range_(Ɀ_TimeYesterday_(11, 0), Ɀ_TimeYesterday_(12, 30)), nil)
			return r
		}(),
		func() Record {
			r := NewRecord(today)
			// Has grace period of 30 minutes.
			r.AddRange(Ɀ_Range_(Ɀ_Time_(0, 0), Ɀ_Time_(12, 30)), nil)
			// If the total time exceeds “now”, that’s okay. (0:00-12:30 + 2h would be 14:30)
			r.AddDuration(NewDuration(2, 0), nil)
			return r
		}(),
	}
	ws := CheckForWarnings(timestamp, rs)
	assert.Equal(t, 0, CountWarningsOfKind(&futureEntriesChecker{}, ws))
}

func TestFutureEntriesWarning(t *testing.T) {
	timestamp := gotime.Date(2000, 3, 5, 12, 00, 0, 0, gotime.Local)
	today := NewDateFromGo(timestamp)
	rs := []Record{
		func() Record {
			r := NewRecord(today.PlusDays(1))
			r.AddDuration(NewDuration(2, 0), nil)
			return r
		}(),
		func() Record {
			r := NewRecord(today.PlusDays(4))
			r.AddDuration(NewDuration(2, 0), nil)
			return r
		}(),
		func() Record {
			r := NewRecord(today.PlusDays(1))
			r.AddRange(Ɀ_Range_(Ɀ_Time_(2, 0), Ɀ_Time_(10, 0)), nil)
			return r
		}(),
		func() Record {
			r := NewRecord(today)
			r.AddRange(Ɀ_Range_(Ɀ_Time_(11, 00), Ɀ_Time_(12, 31)), nil)
			return r
		}(),
		func() Record {
			r := NewRecord(today)
			// Times shifted to next day
			r.AddRange(Ɀ_Range_(Ɀ_TimeTomorrow_(1, 00), Ɀ_TimeTomorrow_(3, 0)), nil)
			return r
		}(),
		func() Record {
			r := NewRecord(today.PlusDays(1))
			// Times shifted to yesterday, but there is also a duration
			r.AddRange(Ɀ_Range_(Ɀ_TimeYesterday_(11, 0), Ɀ_TimeYesterday_(12, 30)), nil)
			r.AddDuration(NewDuration(2, 0), nil)
			return r
		}(),
		func() Record {
			r := NewRecord(today)
			r.StartOpenRange(Ɀ_Time_(12, 31), nil)
			return r
		}(),
	}
	ws := CheckForWarnings(timestamp, rs)
	assert.Equal(t, len(rs), CountWarningsOfKind(&futureEntriesChecker{}, ws))
}

func TestNoWarnForMoreThan24HoursPerRecord(t *testing.T) {
	timestamp := gotime.Date(2000, 3, 5, 12, 00, 0, 0, gotime.Local)
	today := NewDateFromGo(timestamp)
	rs := []Record{
		func() Record {
			r := NewRecord(today.PlusDays(-3))
			r.AddDuration(NewDuration(24, 0), nil)
			return r
		}(),
	}
	ws := CheckForWarnings(timestamp, rs)
	assert.Equal(t, 0, CountWarningsOfKind(&moreThan24HoursChecker{}, ws))
}

func TestMoreThan24HoursPerRecord(t *testing.T) {
	timestamp := gotime.Date(2000, 3, 5, 12, 00, 0, 0, gotime.Local)
	today := NewDateFromGo(timestamp)
	rs := []Record{
		func() Record {
			r := NewRecord(today.PlusDays(-1))
			r.AddDuration(NewDuration(24, 1), nil)
			return r
		}(), func() Record {
			r := NewRecord(today)
			r.AddRange(Ɀ_Range_(Ɀ_Time_(0, 0), Ɀ_Time_(12, 0)), nil)
			r.AddDuration(NewDuration(13, 0), nil)
			return r
		}(),
	}
	ws := CheckForWarnings(timestamp, rs)
	assert.Equal(t, len(rs), CountWarningsOfKind(&moreThan24HoursChecker{}, ws))
}

func TestNoWarnForOverlappingTimeRanges(t *testing.T) {
	timestamp := gotime.Date(2000, 3, 5, 12, 00, 0, 0, gotime.Local)
	today := NewDateFromGo(timestamp)
	rs := []Record{
		func() Record {
			// No overlap
			r := NewRecord(today.PlusDays(-9999))
			r.AddDuration(NewDuration(5, 0), nil)
			r.AddRange(Ɀ_Range_(Ɀ_Time_(4, 0), Ɀ_Time_(4, 59)), nil)
			r.AddRange(Ɀ_Range_(Ɀ_Time_(0, 0), Ɀ_Time_(2, 0)), nil)
			r.AddRange(Ɀ_Range_(Ɀ_Time_(2, 0), Ɀ_Time_(4, 0)), nil)
			r.AddRange(Ɀ_Range_(Ɀ_Time_(4, 0), Ɀ_Time_(4, 0)), nil) // point in time range
			r.AddRange(Ɀ_Range_(Ɀ_Time_(5, 0), Ɀ_Time_(6, 0)), nil)
			return r
		}(),
	}
	ws := CheckForWarnings(timestamp, rs)
	assert.Equal(t, 0, CountWarningsOfKind(&overlappingTimeRangesChecker{}, ws))
}

func TestOverlappingTimeRanges(t *testing.T) {
	timestamp := gotime.Date(2000, 3, 5, 12, 00, 0, 0, gotime.Local)
	today := NewDateFromGo(timestamp)
	rs := []Record{
		func() Record {
			// Overlap with started time
			r := NewRecord(today)
			r.AddRange(Ɀ_Range_(Ɀ_Time_(1, 30), Ɀ_Time_(5, 45)), nil)
			r.StartOpenRange(Ɀ_Time_(3, 0), nil)
			return r
		}(), func() Record {
			// Overlap with sorted entries
			r := NewRecord(today.PlusDays(-1))
			r.AddRange(Ɀ_Range_(Ɀ_Time_(0, 30), Ɀ_Time_(1, 0)), nil)
			r.AddRange(Ɀ_Range_(Ɀ_Time_(2, 0), Ɀ_Time_(5, 0)), nil)
			r.AddRange(Ɀ_Range_(Ɀ_Time_(4, 59), Ɀ_Time_(6, 0)), nil)
			r.AddRange(Ɀ_Range_(Ɀ_Time_(18, 30), Ɀ_Time_(19, 0)), nil)
			return r
		}(), func() Record {
			// Overlap with unsorted entries
			r := NewRecord(today.PlusDays(-2))
			r.AddRange(Ɀ_Range_(Ɀ_Time_(0, 30), Ɀ_Time_(0, 45)), nil)
			r.AddRange(Ɀ_Range_(Ɀ_Time_(2, 45), Ɀ_Time_(3, 45)), nil)
			r.AddRange(Ɀ_Range_(Ɀ_TimeYesterday_(23, 0), Ɀ_Time_(1, 0)), nil)
			return r
		}(),
	}
	ws := CheckForWarnings(timestamp, rs)
	assert.Equal(t, len(rs), CountWarningsOfKind(&overlappingTimeRangesChecker{}, ws))
}
