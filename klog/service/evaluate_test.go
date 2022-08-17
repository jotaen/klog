package service

import (
	"github.com/jotaen/klog/klog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTotalSumUpZeroIfNoTimesSpecified(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	assert.Equal(t, klog.NewDuration(0, 0), Total(r))
}

func TestTotalSumsUpTimesAndRangesButNotOpenRanges(t *testing.T) {
	r1 := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r1.AddDuration(klog.NewDuration(3, 0), nil)
	r1.AddDuration(klog.NewDuration(1, 33), nil)
	r1.AddRange(klog.Ɀ_Range_(klog.Ɀ_TimeYesterday_(8, 0), klog.Ɀ_TimeTomorrow_(12, 0)), nil)
	r1.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(13, 49), klog.Ɀ_Time_(17, 12)), nil)
	_ = r1.StartOpenRange(klog.Ɀ_Time_(1, 2), nil)
	r2 := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 2))
	r2.AddDuration(klog.NewDuration(7, 55), nil)
	assert.Equal(t, klog.NewDuration(3+1+(16+24+12)+3+7, 33+11+12+55), Total(r1, r2))
}
