package record

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	. "klog/testutil/datetime"
	"testing"
)

func TestSumUpTimes(t *testing.T) {
	r := NewRecord(Date_(2020, 1, 1))
	r.AddDuration(datetime.NewDuration(1, 0))
	r.AddDuration(datetime.NewDuration(2, 0))
	assert.Equal(t, datetime.NewDuration(3, 0), TotalWorkTime(r))
}

func TestSumUpZeroIfNoTimesAvailable(t *testing.T) {
	r := NewRecord(Date_(2020, 1, 1))
	assert.Equal(t, datetime.NewDuration(0, 0), TotalWorkTime(r))
}

func TestSumUpRanges(t *testing.T) {
	range1 := Range_(Time_(9, 7), Time_(12, 59))
	range2 := Range_(Time_(13, 49), Time_(17, 12))
	r := NewRecord(Date_(2020, 1, 1))
	r.AddRange(range1)
	r.AddRange(range2)
	assert.Equal(t, datetime.NewDuration(7, 15), TotalWorkTime(r))
}

func TestSumUpTimesAndRanges(t *testing.T) {
	range1 := Range_(Time_(8, 0), Time_(12, 0))
	r := NewRecord(Date_(2020, 1, 1))
	r.AddDuration(datetime.NewDuration(1, 33))
	r.AddRange(range1)
	assert.Equal(t, datetime.NewDuration(5, 33), TotalWorkTime(r))
}
