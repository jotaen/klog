package service

import (
	. "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"testing"
	gotime "time"
)

func TestTotalSumUpZeroIfNoTimesSpecified(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	assert.Equal(t, NewDuration(0, 0), Total(r))
}

func TestTotalSumsUpTimesAndRangesButNotOpenRanges(t *testing.T) {
	r1 := NewRecord(Ɀ_Date_(2020, 1, 1))
	r1.AddDuration(NewDuration(3, 0), NewSummary())
	r1.AddDuration(NewDuration(1, 33), NewSummary())
	r1.AddRange(Ɀ_Range_(Ɀ_TimeYesterday_(8, 0), Ɀ_TimeTomorrow_(12, 0)), NewSummary())
	r1.AddRange(Ɀ_Range_(Ɀ_Time_(13, 49), Ɀ_Time_(17, 12)), NewSummary())
	_ = r1.StartOpenRange(Ɀ_Time_(1, 2), NewSummary())
	r2 := NewRecord(Ɀ_Date_(2020, 1, 2))
	r2.AddDuration(NewDuration(7, 55), NewSummary())
	assert.Equal(t, NewDuration(3+1+(16+24+12)+3+7, 33+11+12+55), Total(r1, r2))
}

func TestSumUpHypotheticalTotalAtGivenTime(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.AddDuration(NewDuration(2, 14), NewSummary())
	r.AddRange(Ɀ_Range_(Ɀ_TimeYesterday_(23, 0), Ɀ_Time_(4, 0)), NewSummary())
	_ = r.StartOpenRange(Ɀ_Time_(5, 7), NewSummary())

	time1, _ := gotime.Parse("2006-01-02T15:04:05-0700", "2020-01-01T05:06:59-0000")
	ht1, isOngoing1 := HypotheticalTotal(time1, r)
	assert.False(t, isOngoing1)
	assert.Equal(t, NewDuration(2+(1+4), 14), ht1)

	time2, _ := gotime.Parse("2006-01-02T15:04:05-0700", "2020-01-01T10:48:13-0000")
	ht2, isOngoing2 := HypotheticalTotal(time2, r)
	assert.True(t, isOngoing2)
	assert.Equal(t, NewDuration(2+(1+4)+4, 14+53+48), ht2)

	time3, _ := gotime.Parse("2006-01-02T15:04:05-0700", "2020-01-02T03:01:29-0000")
	ht3, isOngoing3 := HypotheticalTotal(time3, r)
	assert.True(t, isOngoing3)
	assert.Equal(t, NewDuration(2+(1+4)+18+3, 14+53+1), ht3)
}
