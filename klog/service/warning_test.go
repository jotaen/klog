package service

import (
	"testing"
	gotime "time"

	"github.com/jotaen/klog/klog"
	"github.com/stretchr/testify/assert"
)

func countWarningsOfKind(c checker, ws []Warning) int {
	count := 0
	for _, w := range ws {
		if w.Warning() == c.Message() {
			count++
		}
	}
	return count
}

func checkForWarningsWithCollect(reference gotime.Time, rs []klog.Record) []Warning {
	var ws []Warning
	CheckForWarnings(func(w Warning) {
		ws = append(ws, w)
	}, reference, rs, []string{})
	return ws
}

func TestNoWarnForOpenRanges(t *testing.T) {
	timestamp := gotime.Date(2000, 3, 5, 12, 00, 0, 0, gotime.Local)
	today := klog.NewDateFromGo(timestamp)
	now := klog.NewTimeFromGo(timestamp)

	rs1 := []klog.Record{
		func() klog.Record {
			// This open range is okay, because there is no record at today’s date
			r := klog.NewRecord(today.PlusDays(-1))
			r.Start(klog.NewOpenRange(now), nil)
			return r
		}(), func() klog.Record {
			r := klog.NewRecord(today.PlusDays(2))
			return r
		}(),
	}
	ws1 := checkForWarningsWithCollect(timestamp, rs1)
	assert.Equal(t, 0, countWarningsOfKind(&unclosedOpenRangeChecker{}, ws1))

	rs2 := []klog.Record{
		func() klog.Record {
			r := klog.NewRecord(today)
			r.Start(klog.NewOpenRange(now), nil)
			return r
		}(),
	}
	ws2 := checkForWarningsWithCollect(timestamp, rs2)
	assert.Equal(t, 0, countWarningsOfKind(&unclosedOpenRangeChecker{}, ws2))
}

func TestOpenRangeWarningWhenUnclosedOpenRangeBeforeTodayRegardlessOfOrder(t *testing.T) {
	timestamp := gotime.Date(2000, 3, 5, 12, 00, 0, 0, gotime.Local)
	today := klog.NewDateFromGo(timestamp)
	now := klog.NewTimeFromGo(timestamp)
	// The warnings must work reliably even when the records are not ordered by date initially
	rs := []klog.Record{
		func() klog.Record {
			// NOT OK: There is a record at today’s date
			r := klog.NewRecord(today.PlusDays(-1))
			r.Start(klog.NewOpenRange(now), nil)
			return r
		}(), func() klog.Record {
			r := klog.NewRecord(today)
			return r
		}(), func() klog.Record {
			// NOT OK: There is a record at today’s date
			r := klog.NewRecord(today.PlusDays(-2))
			r.Start(klog.NewOpenRange(now), nil)
			return r
		}(),
	}
	ws := checkForWarningsWithCollect(timestamp, rs)
	assert.Equal(t, 2, countWarningsOfKind(&unclosedOpenRangeChecker{}, ws))
	assert.Equal(t, today.PlusDays(-1), ws[0].Date())
	assert.Equal(t, today.PlusDays(-2), ws[1].Date())
}

func TestNoWarningForFutureEntries(t *testing.T) {
	timestamp := gotime.Date(2000, 3, 5, 12, 00, 0, 0, gotime.Local)
	today := klog.NewDateFromGo(timestamp)
	rs := []klog.Record{
		func() klog.Record {
			// Future entry okay if it doesn’t contain entries
			r := klog.NewRecord(today.PlusDays(1))
			return r
		}(),
		func() klog.Record {
			r := klog.NewRecord(today.PlusDays(-1))
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(0, 0), klog.Ɀ_Time_(12, 30)), nil)
			r.AddDuration(klog.NewDuration(2, 0), nil)
			return r
		}(),
		func() klog.Record {
			r := klog.NewRecord(today.PlusDays(1))
			// Times shifted to yesterday
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_TimeYesterday_(11, 0), klog.Ɀ_TimeYesterday_(12, 30)), nil)
			return r
		}(),
		func() klog.Record {
			r := klog.NewRecord(today)
			// Has grace period of 30 minutes.
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(0, 0), klog.Ɀ_Time_(12, 30)), nil)
			// If the total time exceeds “now”, that’s okay. (0:00-12:30 + 2h would be 14:30)
			r.AddDuration(klog.NewDuration(2, 0), nil)
			return r
		}(),
	}
	ws := checkForWarningsWithCollect(timestamp, rs)
	assert.Equal(t, 0, countWarningsOfKind(&futureEntriesChecker{}, ws))
}

func TestFutureEntriesWarning(t *testing.T) {
	timestamp := gotime.Date(2000, 3, 5, 12, 00, 0, 0, gotime.Local)
	today := klog.NewDateFromGo(timestamp)
	rs := []klog.Record{
		func() klog.Record {
			r := klog.NewRecord(today.PlusDays(1))
			r.AddDuration(klog.NewDuration(2, 0), nil)
			return r
		}(),
		func() klog.Record {
			r := klog.NewRecord(today.PlusDays(4))
			r.AddDuration(klog.NewDuration(2, 0), nil)
			return r
		}(),
		func() klog.Record {
			r := klog.NewRecord(today.PlusDays(1))
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(2, 0), klog.Ɀ_Time_(10, 0)), nil)
			return r
		}(),
		func() klog.Record {
			r := klog.NewRecord(today)
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(11, 00), klog.Ɀ_Time_(12, 31)), nil)
			return r
		}(),
		func() klog.Record {
			r := klog.NewRecord(today)
			// Times shifted to next day
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_TimeTomorrow_(1, 00), klog.Ɀ_TimeTomorrow_(3, 0)), nil)
			return r
		}(),
		func() klog.Record {
			r := klog.NewRecord(today.PlusDays(1))
			// Times shifted to yesterday, but there is also a duration
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_TimeYesterday_(11, 0), klog.Ɀ_TimeYesterday_(12, 30)), nil)
			r.AddDuration(klog.NewDuration(2, 0), nil)
			return r
		}(),
		func() klog.Record {
			r := klog.NewRecord(today)
			r.Start(klog.NewOpenRange(klog.Ɀ_Time_(12, 31)), nil)
			return r
		}(),
	}
	ws := checkForWarningsWithCollect(timestamp, rs)
	assert.Equal(t, len(rs), countWarningsOfKind(&futureEntriesChecker{}, ws))
}

func TestNoWarnForMoreThan24HoursPerRecord(t *testing.T) {
	timestamp := gotime.Date(2000, 3, 5, 12, 00, 0, 0, gotime.Local)
	today := klog.NewDateFromGo(timestamp)
	rs := []klog.Record{
		func() klog.Record {
			r := klog.NewRecord(today.PlusDays(-3))
			r.AddDuration(klog.NewDuration(24, 0), nil)
			return r
		}(),
	}
	ws := checkForWarningsWithCollect(timestamp, rs)
	assert.Equal(t, 0, countWarningsOfKind(&moreThan24HoursChecker{}, ws))
}

func TestMoreThan24HoursPerRecord(t *testing.T) {
	timestamp := gotime.Date(2000, 3, 5, 12, 00, 0, 0, gotime.Local)
	today := klog.NewDateFromGo(timestamp)
	rs := []klog.Record{
		func() klog.Record {
			r := klog.NewRecord(today.PlusDays(-1))
			r.AddDuration(klog.NewDuration(24, 1), nil)
			return r
		}(), func() klog.Record {
			r := klog.NewRecord(today)
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(0, 0), klog.Ɀ_Time_(12, 0)), nil)
			r.AddDuration(klog.NewDuration(13, 0), nil)
			return r
		}(),
	}
	ws := checkForWarningsWithCollect(timestamp, rs)
	assert.Equal(t, len(rs), countWarningsOfKind(&moreThan24HoursChecker{}, ws))
}

func TestNoWarnForOverlappingTimeRanges(t *testing.T) {
	timestamp := gotime.Date(2000, 3, 5, 12, 00, 0, 0, gotime.Local)
	today := klog.NewDateFromGo(timestamp)
	rs := []klog.Record{
		func() klog.Record {
			// No overlap
			r := klog.NewRecord(today.PlusDays(-9999))
			r.AddDuration(klog.NewDuration(5, 0), nil)
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(4, 0), klog.Ɀ_Time_(4, 59)), nil)
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(0, 0), klog.Ɀ_Time_(2, 0)), nil)
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(2, 0), klog.Ɀ_Time_(4, 0)), nil)
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(4, 0), klog.Ɀ_Time_(4, 0)), nil) // point in time range
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(5, 0), klog.Ɀ_Time_(6, 0)), nil)
			return r
		}(),
	}
	ws := checkForWarningsWithCollect(timestamp, rs)
	assert.Equal(t, 0, countWarningsOfKind(&overlappingTimeRangesChecker{}, ws))
}

func TestOverlappingTimeRanges(t *testing.T) {
	timestamp := gotime.Date(2000, 3, 5, 12, 00, 0, 0, gotime.Local)
	today := klog.NewDateFromGo(timestamp)
	rs := []klog.Record{
		func() klog.Record {
			// Overlap with started time
			r := klog.NewRecord(today)
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(1, 30), klog.Ɀ_Time_(5, 45)), nil)
			r.Start(klog.NewOpenRange(klog.Ɀ_Time_(3, 0)), nil)
			return r
		}(), func() klog.Record {
			// Overlap with sorted entries
			r := klog.NewRecord(today.PlusDays(-1))
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(0, 30), klog.Ɀ_Time_(1, 0)), nil)
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(2, 0), klog.Ɀ_Time_(5, 0)), nil)
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(4, 59), klog.Ɀ_Time_(6, 0)), nil)
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(18, 30), klog.Ɀ_Time_(19, 0)), nil)
			return r
		}(), func() klog.Record {
			// Overlap with unsorted entries
			r := klog.NewRecord(today.PlusDays(-2))
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(0, 30), klog.Ɀ_Time_(0, 45)), nil)
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(2, 45), klog.Ɀ_Time_(3, 45)), nil)
			r.AddRange(klog.Ɀ_Range_(klog.Ɀ_TimeYesterday_(23, 0), klog.Ɀ_Time_(1, 0)), nil)
			return r
		}(),
	}
	ws := checkForWarningsWithCollect(timestamp, rs)
	assert.Equal(t, len(rs), countWarningsOfKind(&overlappingTimeRangesChecker{}, ws))
}
