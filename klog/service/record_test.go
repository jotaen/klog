package service

import (
	"github.com/jotaen/klog/klog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	gotime "time"
)

func TestReturnsOriginalRecordsIfNoOpenRange(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r.AddDuration(klog.NewDuration(1, 0), nil)
	r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(1, 0), klog.Ɀ_Time_(2, 0)), nil)

	result, err := CloseOpenRanges(gotime.Now(), r)
	require.Nil(t, err)
	assert.Equal(t, klog.NewDuration(2, 0), Total(result...))
}

func TestReturnsRecordsWithClosedOpenRange(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r.AddDuration(klog.NewDuration(1, 0), nil)
	r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(1, 0), klog.Ɀ_Time_(2, 0)), nil)
	r.StartOpenRange(klog.Ɀ_Time_(3, 0), nil)

	endTime, _ := gotime.Parse("2006-01-02T15:04:05-0700", "2020-01-01T05:30:00-0000")
	result, err := CloseOpenRanges(endTime, r)
	require.Nil(t, err)
	assert.Equal(t, klog.NewDuration(2+2, 30), Total(result...))
}

func TestReturnsRecordsWithClosedOpenRangeAndShiftedTime(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r.AddDuration(klog.NewDuration(1, 0), nil)
	r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(1, 0), klog.Ɀ_Time_(2, 0)), nil)
	r.StartOpenRange(klog.Ɀ_Time_(3, 0), nil)

	endTime, _ := gotime.Parse("2006-01-02T15:04:05-0700", "2020-01-02T05:30:00-0000")
	result, err := CloseOpenRanges(endTime, r)
	require.Nil(t, err)
	assert.Equal(t, klog.NewDuration(2+24+2, 30), Total(result...))
}

func TestReturnsErrorIfOpenRangeCannotBeClosedAnymore(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r.AddDuration(klog.NewDuration(1, 0), nil)
	r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(1, 0), klog.Ɀ_Time_(2, 0)), nil)
	r.StartOpenRange(klog.Ɀ_Time_(3, 0), nil)

	endTime, _ := gotime.Parse("2006-01-02T15:04:05-0700", "2020-01-03T05:30:00-0000")
	result, err := CloseOpenRanges(endTime, r)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestReturnsErrorIfOpenRangeCannotBeClosedYet(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r.AddDuration(klog.NewDuration(1, 0), nil)
	r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(1, 0), klog.Ɀ_Time_(2, 0)), nil)
	r.StartOpenRange(klog.Ɀ_Time_(3, 0), nil)

	endTime, _ := gotime.Parse("2006-01-02T15:04:05-0700", "2020-01-01T01:30:00-0000")
	result, err := CloseOpenRanges(endTime, r)
	require.Error(t, err)
	assert.Nil(t, result)
}
