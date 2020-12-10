package workday

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	"testing"
)

func TestSumUpTimes(t *testing.T) {
	date, _ := datetime.CreateDate(2020, 1, 1)
	w := Create(date)
	w.AddDuration(datetime.Duration(60))
	w.AddDuration(datetime.Duration(120))
	assert.Equal(t, datetime.Duration(180), w.TotalWorkTime())
}

func TestSumUpZeroIfNoTimesAvailable(t *testing.T) {
	date, _ := datetime.CreateDate(2020, 1, 1)
	w := Create(date)
	assert.Equal(t, datetime.Duration(0), w.TotalWorkTime())
}

func TestSumUpRanges(t *testing.T) {
	date, _ := datetime.CreateDate(2020, 1, 1)
	time1, _ := datetime.CreateTime(9, 7)
	time2, _ := datetime.CreateTime(12, 59)
	time3, _ := datetime.CreateTime(13, 49)
	time4, _ := datetime.CreateTime(17, 12)
	range1, _ := datetime.CreateTimeRange(time1, time2)
	range2, _ := datetime.CreateTimeRange(time3, time4)
	w := Create(date)
	w.AddRange(range1)
	w.AddRange(range2)
	assert.Equal(t, datetime.Duration(435), w.TotalWorkTime())
}

func TestSumUpTimesAndRanges(t *testing.T) {
	date, _ := datetime.CreateDate(2020, 1, 1)
	time1, _ := datetime.CreateTime(8, 0)
	time2, _ := datetime.CreateTime(12, 0)
	range1, _ := datetime.CreateTimeRange(time1, time2)
	w := Create(date)
	w.AddDuration(datetime.Duration(93))
	w.AddRange(range1)
	assert.Equal(t, datetime.Duration(333), w.TotalWorkTime())
}

func TestDisregardsOpenRanges(t *testing.T) {
	date, _ := datetime.CreateDate(2020, 1, 1)
	time1, _ := datetime.CreateTime(9, 7)
	range1, _ := datetime.CreateTimeRange(time1, nil)
	w := Create(date)
	w.AddRange(range1)
	assert.Equal(t, datetime.Duration(0), w.TotalWorkTime())
}
