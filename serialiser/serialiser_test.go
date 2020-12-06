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
