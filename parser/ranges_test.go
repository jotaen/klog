package parser

import (
	"github.com/stretchr/testify/assert"
	"klog/workday"
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
	assert.Equal(t, e.Ranges(), [][]workday.Time{
		[]workday.Time{workday.Time{Hour: 9, Minute: 12}, workday.Time{Hour: 12, Minute: 5}},
		[]workday.Time{workday.Time{Hour: 13, Minute: 3}, workday.Time{Hour: 16, Minute: 47}},
	})
}
