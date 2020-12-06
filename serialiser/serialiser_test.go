package serialiser

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	"klog/workday"
	"testing"
)

func TestSerialiseDate(t *testing.T) {
	workDay, _ := workday.Create(datetime.Date{Year: 1859, Month: 6, Day: 2})
	text := Serialise(workDay)
	assert.Equal(t, `date: 1859-06-02
`, text)
}

func TestSerialiseSummaryIfPresent(t *testing.T) {
	workDay, _ := workday.Create(datetime.Date{Year: 1859, Month: 6, Day: 2})
	workDay.SetSummary("Hello World!")
	text := Serialise(workDay)
	assert.Equal(t, `date: 1859-06-02
summary: Hello World!
`, text)
}
