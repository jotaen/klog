package serialiser

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	"klog/workday"
	"testing"
)

func TestSerialiseDate(t *testing.T) {
	date, _ := datetime.CreateDate(1859, 6, 2)
	workDay := workday.Create(date)
	text := Serialise(workDay)
	assert.Equal(t, `date: 1859-06-02
`, text)
}

func TestSerialiseSummaryIfPresent(t *testing.T) {
	date, _ := datetime.CreateDate(1859, 6, 2)
	workDay := workday.Create(date)
	workDay.SetSummary("Hello World!")
	text := Serialise(workDay)
	assert.Equal(t, `date: 1859-06-02
summary: Hello World!
`, text)
}

func TestSerialiseRanges(t *testing.T) {
	date, _ := datetime.CreateDate(1859, 6, 2)
	workDay := workday.Create(date)
	time1, _ := datetime.CreateTime(8, 31)
	time2, _ := datetime.CreateTime(14, 2)
	workDay.AddRange(time1, time2)
	time3, _ := datetime.CreateTime(15, 0)
	workDay.AddOpenRange(time3)
	text := Serialise(workDay)
	assert.Equal(t, `date: 1859-06-02
hours:
- start: 08:31
  end: 14:02
- start: 15:00
`, text)
}

func TestSerialiseTimes(t *testing.T) {
	date, _ := datetime.CreateDate(1859, 6, 2)
	workDay := workday.Create(date)
	workDay.AddTime(datetime.Duration(3))
	workDay.AddTime(datetime.Duration(819))
	text := Serialise(workDay)
	assert.Equal(t, `date: 1859-06-02
hours:
- time: 00:03
- time: 13:39
`, text)
}
