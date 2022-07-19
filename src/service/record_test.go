package service

import (
	. "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	gotime "time"
)

func TestReturnsOriginalRecordsIfNoOpenRange(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.AddDuration(NewDuration(1, 0), nil)
	r.AddRange(Ɀ_Range_(Ɀ_Time_(1, 0), Ɀ_Time_(2, 0)), nil)

	result, err := CloseOpenRanges(gotime.Now(), r)
	require.Nil(t, err)
	assert.Equal(t, NewDuration(2, 0), Total(result...))
}

func TestReturnsRecordsWithClosedOpenRange(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.AddDuration(NewDuration(1, 0), nil)
	r.AddRange(Ɀ_Range_(Ɀ_Time_(1, 0), Ɀ_Time_(2, 0)), nil)
	r.StartOpenRange(Ɀ_Time_(3, 0), nil)

	endTime, _ := gotime.Parse("2006-01-02T15:04:05-0700", "2020-01-01T05:30:00-0000")
	result, err := CloseOpenRanges(endTime, r)
	require.Nil(t, err)
	assert.Equal(t, NewDuration(2+2, 30), Total(result...))
}

func TestReturnsRecordsWithClosedOpenRangeAndShiftedTime(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.AddDuration(NewDuration(1, 0), nil)
	r.AddRange(Ɀ_Range_(Ɀ_Time_(1, 0), Ɀ_Time_(2, 0)), nil)
	r.StartOpenRange(Ɀ_Time_(3, 0), nil)

	endTime, _ := gotime.Parse("2006-01-02T15:04:05-0700", "2020-01-02T05:30:00-0000")
	result, err := CloseOpenRanges(endTime, r)
	require.Nil(t, err)
	assert.Equal(t, NewDuration(2+24+2, 30), Total(result...))
}

func TestReturnsErrorIfOpenRangeCannotBeClosedAnymore(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.AddDuration(NewDuration(1, 0), nil)
	r.AddRange(Ɀ_Range_(Ɀ_Time_(1, 0), Ɀ_Time_(2, 0)), nil)
	r.StartOpenRange(Ɀ_Time_(3, 0), nil)

	endTime, _ := gotime.Parse("2006-01-02T15:04:05-0700", "2020-01-03T05:30:00-0000")
	result, err := CloseOpenRanges(endTime, r)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestReturnsErrorIfOpenRangeCannotBeClosedYet(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.AddDuration(NewDuration(1, 0), nil)
	r.AddRange(Ɀ_Range_(Ɀ_Time_(1, 0), Ɀ_Time_(2, 0)), nil)
	r.StartOpenRange(Ɀ_Time_(3, 0), nil)

	endTime, _ := gotime.Parse("2006-01-02T15:04:05-0700", "2020-01-01T01:30:00-0000")
	result, err := CloseOpenRanges(endTime, r)
	require.Error(t, err)
	assert.Nil(t, result)
}
