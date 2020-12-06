package workday

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	"testing"
)

func TestSumUpTimes(t *testing.T) {
	w, _ := Create(datetime.Date{Year: 2020, Month: 1, Day: 1})
	w.AddTime(datetime.Minutes(60))
	w.AddTime(datetime.Minutes(120))
	assert.Equal(t, datetime.Minutes(180), w.TotalTime())
}

func TestSumUpRanges(t *testing.T) {
	w, _ := Create(datetime.Date{Year: 2020, Month: 1, Day: 1})
	w.AddRange(datetime.Time{Hour: 9, Minute: 07}, datetime.Time{Hour: 12, Minute: 59})
	w.AddRange(datetime.Time{Hour: 13, Minute: 49}, datetime.Time{Hour: 17, Minute: 12})
	assert.Equal(t, datetime.Minutes(435), w.TotalTime())
}

func TestSumUpTimesAndRanges(t *testing.T) {
	w, _ := Create(datetime.Date{Year: 2020, Month: 1, Day: 1})
	w.AddTime(datetime.Minutes(93))
	w.AddRange(datetime.Time{Hour: 8, Minute: 00}, datetime.Time{Hour: 12, Minute: 00})
	assert.Equal(t, datetime.Minutes(333), w.TotalTime())
}
