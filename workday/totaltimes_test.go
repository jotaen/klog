package workday

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSumUpTimes(t *testing.T) {
	date := Date{Year: 2020, Month: time.January, Day: 1}
	w, _ := Create(date)
	w.AddTime(Minutes(60))
	w.AddTime(Minutes(120))
	assert.Equal(t, Minutes(180), w.TotalTime())
}

func TestSumUpRanges(t *testing.T) {
	date := Date{Year: 2020, Month: time.January, Day: 1}
	w, _ := Create(date)
	w.AddRange(Time{Hour: 9, Minute: 07}, Time{Hour: 12, Minute: 59})
	w.AddRange(Time{Hour: 13, Minute: 49}, Time{Hour: 17, Minute: 12})
	assert.Equal(t, Minutes(435), w.TotalTime())
}

func TestSumUpTimesAndRanges(t *testing.T) {
	date := Date{Year: 2020, Month: time.January, Day: 1}
	w, _ := Create(date)
	w.AddTime(Minutes(93))
	w.AddRange(Time{Hour: 8, Minute: 00}, Time{Hour: 12, Minute: 00})
	assert.Equal(t, Minutes(333), w.TotalTime())
}
