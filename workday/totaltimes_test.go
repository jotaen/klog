package workday

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	datetime2 "klog/testutil/datetime"
	"testing"
)

func TestSumUpTimes(t *testing.T) {
	w := NewWorkDay(datetime2.Date_(2020, 1, 1))
	w.AddDuration(datetime.NewDuration(1, 0))
	w.AddDuration(datetime.NewDuration(2, 0))
	assert.Equal(t, datetime.NewDuration(3, 0), w.TotalWorkTime())
}

func TestSumUpZeroIfNoTimesAvailable(t *testing.T) {
	w := NewWorkDay(datetime2.Date_(2020, 1, 1))
	assert.Equal(t, datetime.NewDuration(0, 0), w.TotalWorkTime())
}

func TestSumUpRanges(t *testing.T) {
	range1 := datetime2.Range_(datetime2.Time_(9, 7), datetime2.Time_(12, 59))
	range2 := datetime2.Range_(datetime2.Time_(13, 49), datetime2.Time_(17, 12))
	w := NewWorkDay(datetime2.Date_(2020, 1, 1))
	w.AddRange(range1)
	w.AddRange(range2)
	assert.Equal(t, datetime.NewDuration(7, 15), w.TotalWorkTime())
}

func TestSumUpTimesAndRanges(t *testing.T) {
	range1 := datetime2.Range_(datetime2.Time_(8, 0), datetime2.Time_(12, 0))
	w := NewWorkDay(datetime2.Date_(2020, 1, 1))
	w.AddDuration(datetime.NewDuration(1, 33))
	w.AddRange(range1)
	assert.Equal(t, datetime.NewDuration(5, 33), w.TotalWorkTime())
}
