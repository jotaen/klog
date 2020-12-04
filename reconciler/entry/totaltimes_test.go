package entry

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSumUpTimes(t *testing.T) {
	date := Date{Year: 2020, Month: time.January, Day: 1}
	entry, _ := Create(date)
	entry.AddTime(Minutes(60))
	entry.AddTime(Minutes(120))
	assert.Equal(t, Minutes(180), entry.TotalTime())
}

func TestSumUpRanges(t *testing.T) {
	date := Date{Year: 2020, Month: time.January, Day: 1}
	entry, _ := Create(date)
	entry.AddRange(Time{Hour: 9, Minute: 07}, Time{Hour: 12, Minute: 59})
	entry.AddRange(Time{Hour: 13, Minute: 49}, Time{Hour: 17, Minute: 12})
	assert.Equal(t, Minutes(435), entry.TotalTime())
}

func TestSumUpTimesAndRanges(t *testing.T) {
	date := Date{Year: 2020, Month: time.January, Day: 1}
	entry, _ := Create(date)
	entry.AddTime(Minutes(93))
	entry.AddRange(Time{Hour: 8, Minute: 00}, Time{Hour: 12, Minute: 00})
	assert.Equal(t, Minutes(333), entry.TotalTime())
}
