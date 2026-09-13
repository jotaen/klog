package service

import (
	"testing"
	gotime "time"

	"github.com/jotaen/klog/klog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDoesNotTouchRecordsIfNoOpenRange(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r.AddDuration(klog.NewDuration(1, 0), nil)
	r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(1, 0), klog.Ɀ_Time_(2, 0)), nil)

	closedEntries, err := CloseOpenRanges(gotime.Now(), r)
	require.Nil(t, err)
	assert.Len(t, closedEntries, 0)
	assert.Equal(t, klog.NewDuration(2, 0), Total(r))
}

func TestClosesOpenRange(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r.AddDuration(klog.NewDuration(1, 0), nil)
	r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(1, 0), klog.Ɀ_Time_(2, 0)), nil)
	r.Start(klog.NewOpenRange(klog.Ɀ_Time_(3, 0)), nil)

	endTime, _ := gotime.Parse("2006-01-02T15:04:05-0700", "2020-01-01T05:30:00-0000")
	closedEntries, err := CloseOpenRanges(endTime, r)
	require.Nil(t, err)
	require.Len(t, closedEntries, 1)
	assert.Equal(t, 2, closedEntries[r])
	assert.Equal(t, klog.NewDuration(2+2, 30), Total(r))
}

func TestClosesOpenRangeAndShiftsTime(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r.AddDuration(klog.NewDuration(1, 0), nil)
	r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(1, 0), klog.Ɀ_Time_(2, 0)), nil)
	r.Start(klog.NewOpenRange(klog.Ɀ_Time_(3, 0)), nil)

	endTime, _ := gotime.Parse("2006-01-02T15:04:05-0700", "2020-01-02T05:30:00-0000")
	closedEntries, err := CloseOpenRanges(endTime, r)
	require.Nil(t, err)
	require.Len(t, closedEntries, 1)
	assert.Equal(t, 2, closedEntries[r])
	assert.Equal(t, klog.NewDuration(2+24+2, 30), Total(r))
}

func TestReturnsErrorIfOpenRangeCannotBeClosedAnymore(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r.AddDuration(klog.NewDuration(1, 0), nil)
	r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(1, 0), klog.Ɀ_Time_(2, 0)), nil)
	r.Start(klog.NewOpenRange(klog.Ɀ_Time_(3, 0)), nil)

	endTime, _ := gotime.Parse("2006-01-02T15:04:05-0700", "2020-01-03T05:30:00-0000")
	_, err := CloseOpenRanges(endTime, r)
	require.Error(t, err)
}

func TestReturnsErrorIfOpenRangeCannotBeClosedYet(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r.AddDuration(klog.NewDuration(1, 0), nil)
	r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(1, 0), klog.Ɀ_Time_(2, 0)), nil)
	r.Start(klog.NewOpenRange(klog.Ɀ_Time_(3, 0)), nil)

	endTime, _ := gotime.Parse("2006-01-02T15:04:05-0700", "2020-01-01T01:30:00-0000")
	_, err := CloseOpenRanges(endTime, r)
	require.Error(t, err)
}
