package parser

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	"testing"
)

func TestParseWorkDayWithRanges(t *testing.T) {
	yaml := `
date: 2008-12-03
hours:
- start: 9:12
  end: 12:05
- start: 13:03
  end: 16:47
`
	e, _ := Parse(yaml)
	time1, _ := datetime.CreateTime(9, 12)
	time2, _ := datetime.CreateTime(12, 05)
	time3, _ := datetime.CreateTime(13, 3)
	time4, _ := datetime.CreateTime(16, 47)
	assert.Equal(t, e.Ranges(), [][]datetime.Time{
		[]datetime.Time{time1, time2},
		[]datetime.Time{time3, time4},
	})
}
