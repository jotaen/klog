package workday

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	"klog/testutil"
	"testing"
)

func TestSumUpTimes(t *testing.T) {
	w := Create(testutil.Date_(2020, 1, 1))
	w.AddDuration(datetime.Duration(60))
	w.AddDuration(datetime.Duration(120))
	assert.Equal(t, datetime.Duration(180), w.TotalWorkTime())
}

func TestSumUpZeroIfNoTimesAvailable(t *testing.T) {
	w := Create(testutil.Date_(2020, 1, 1))
	assert.Equal(t, datetime.Duration(0), w.TotalWorkTime())
}

func TestSumUpRanges(t *testing.T) {
	range1 := testutil.Range_(testutil.Time_(9, 7), testutil.Time_(12, 59))
	range2 := testutil.Range_(testutil.Time_(13, 49), testutil.Time_(17, 12))
	w := Create(testutil.Date_(2020, 1, 1))
	w.AddRange(range1)
	w.AddRange(range2)
	assert.Equal(t, datetime.Duration(435), w.TotalWorkTime())
}

func TestSumUpTimesAndRanges(t *testing.T) {
	range1 := testutil.Range_(testutil.Time_(8, 0), testutil.Time_(12, 0))
	w := Create(testutil.Date_(2020, 1, 1))
	w.AddDuration(datetime.Duration(93))
	w.AddRange(range1)
	assert.Equal(t, datetime.Duration(333), w.TotalWorkTime())
}
