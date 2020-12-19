package serialiser

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	datetime2 "klog/testutil/datetime"
	"klog/workday"
	"testing"
)

func TestSerialiseDate(t *testing.T) {
	workDay := workday.NewWorkDay(datetime2.Date_(1859, 6, 2))
	text := Serialise(workDay)
	assert.Equal(t, `date: 1859-06-02
`, text)
}

func TestSerialiseSummaryIfPresent(t *testing.T) {
	workDay := workday.NewWorkDay(datetime2.Date_(1859, 6, 2))
	workDay.SetSummary("Hello World!")
	text := Serialise(workDay)
	assert.Equal(t, `date: 1859-06-02
summary: Hello World!
`, text)
}

func TestSerialiseRanges(t *testing.T) {
	workDay := workday.NewWorkDay(datetime2.Date_(1859, 6, 2))
	workDay.AddRange(datetime2.Range_(datetime2.Time_(8, 31), datetime2.Time_(14, 2)))
	workDay.AddRange(datetime2.Range_(datetime2.Time_(15, 12), datetime2.Time_(18, 7)))
	text := Serialise(workDay)
	assert.Equal(t, `date: 1859-06-02
hours:
- start: 8:31
  end: 14:02
- start: 15:12
  end: 18:07
`, text)
}

func TestSerialiseOpenRange(t *testing.T) {
	workDay := workday.NewWorkDay(datetime2.Date_(1859, 6, 2))
	workDay.StartOpenRange(datetime2.Time_(15, 0))
	text := Serialise(workDay)
	assert.Equal(t, `date: 1859-06-02
hours:
- start: 15:00
`, text)
}

func TestSerialiseTimes(t *testing.T) {
	workDay := workday.NewWorkDay(datetime2.Date_(1859, 6, 2))
	workDay.AddDuration(datetime.NewDuration(0, 3))
	workDay.AddDuration(datetime.NewDuration(6, 39))
	text := Serialise(workDay)
	assert.Equal(t, `date: 1859-06-02
hours:
- time: 3m
- time: 6h 39m
`, text)
}
