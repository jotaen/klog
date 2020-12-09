package workday

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	"testing"
)

func TestSumUpTimes(t *testing.T) {
	date, _ := datetime.CreateDate(2020, 1, 1)
	w := Create(date)
	w.AddTime(datetime.Duration(60))
	w.AddTime(datetime.Duration(120))
	assert.Equal(t, datetime.Duration(180), w.TotalTime())
}

func TestSumUpZeroIfNoTimesAvailable(t *testing.T) {
	date, _ := datetime.CreateDate(2020, 1, 1)
	w := Create(date)
	assert.Equal(t, datetime.Duration(0), w.TotalTime())
}

func TestSumUpRanges(t *testing.T) {
	date, _ := datetime.CreateDate(2020, 1, 1)
	time1, _ := datetime.CreateTime(9, 7)
	time2, _ := datetime.CreateTime(12, 59)
	time3, _ := datetime.CreateTime(13, 49)
	time4, _ := datetime.CreateTime(17, 12)
	w := Create(date)
	w.AddRange(time1, time2)
	w.AddRange(time3, time4)
	assert.Equal(t, datetime.Duration(435), w.TotalTime())
}

func TestSumUpTimesAndRanges(t *testing.T) {
	date, _ := datetime.CreateDate(2020, 1, 1)
	time1, _ := datetime.CreateTime(8, 0)
	time2, _ := datetime.CreateTime(12, 0)
	w := Create(date)
	w.AddTime(datetime.Duration(93))
	w.AddRange(time1, time2)
	assert.Equal(t, datetime.Duration(333), w.TotalTime())
}

func TestDisregardsOpenRanges(t *testing.T) {
	date, _ := datetime.CreateDate(2020, 1, 1)
	time1, _ := datetime.CreateTime(9, 7)
	w := Create(date)
	w.AddOpenRange(time1)
	assert.Equal(t, datetime.Duration(0), w.TotalTime())
}
