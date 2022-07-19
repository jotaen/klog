package service

import (
	. "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTotalSumUpZeroIfNoTimesSpecified(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	assert.Equal(t, NewDuration(0, 0), Total(r))
}

func TestTotalSumsUpTimesAndRangesButNotOpenRanges(t *testing.T) {
	r1 := NewRecord(Ɀ_Date_(2020, 1, 1))
	r1.AddDuration(NewDuration(3, 0), nil)
	r1.AddDuration(NewDuration(1, 33), nil)
	r1.AddRange(Ɀ_Range_(Ɀ_TimeYesterday_(8, 0), Ɀ_TimeTomorrow_(12, 0)), nil)
	r1.AddRange(Ɀ_Range_(Ɀ_Time_(13, 49), Ɀ_Time_(17, 12)), nil)
	_ = r1.StartOpenRange(Ɀ_Time_(1, 2), nil)
	r2 := NewRecord(Ɀ_Date_(2020, 1, 2))
	r2.AddDuration(NewDuration(7, 55), nil)
	assert.Equal(t, NewDuration(3+1+(16+24+12)+3+7, 33+11+12+55), Total(r1, r2))
}
