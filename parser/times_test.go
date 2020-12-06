package parser

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	"testing"
)

func TestParseWorkDayWithTimes(t *testing.T) {
	yaml := `
date: 1985-03-14
hours:
- time: 2:00
- time: 5:00
`
	e, _ := Parse(yaml)
	assert.Equal(t, e.Times(), []datetime.Minutes{datetime.Minutes(2 * 60), datetime.Minutes(5 * 60)})
}

func TestParseWorkDayWithMalformedTimesFails(t *testing.T) {
	yaml := `
date: 1985-03-14
hours:
- time: asdf
`
	e, errs := Parse(yaml)
	assert.Equal(t, e, nil)
	assert.Contains(t, errs, parserError(INVALID_TIME))
}

func TestParseWorkDayWithInvalidTimesFails(t *testing.T) {
	yaml := `
date: 1985-03-14
hours:
- time: 23:87
`
	e, errs := Parse(yaml)
	assert.Equal(t, e, nil)
	assert.Contains(t, errs, parserError(INVALID_TIME))
}
