package parser

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	"testing"
)

func TestMinimalValidDocument(t *testing.T) {
	yaml := `
date: 2020-01-01
`
	w, errs := Parse(yaml)
	assert.Equal(t, w.Summary(), "")
	assert.Nil(t, errs)
}

func TestParsingAllFieldsCorrectly(t *testing.T) {
	yaml := `
date: 2008-12-03
summary: Just a normal day
hours:
- start: 9:12
  end: 12:05
- start: 13:03
  end: 16:47
- time: 2:00
- time: 5:00
`
	time1, _ := datetime.CreateTime(9, 12)
	time2, _ := datetime.CreateTime(12, 05)
	time3, _ := datetime.CreateTime(13, 3)
	time4, _ := datetime.CreateTime(16, 47)

	w, errs := Parse(yaml)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, w.Ranges(), [][]datetime.Time{
		[]datetime.Time{time1, time2},
		[]datetime.Time{time3, time4},
	})
	assert.Equal(t, w.Times(), []datetime.Duration{datetime.Duration(2 * 60), datetime.Duration(5 * 60)})
	assert.Equal(t, w.Summary(), "Just a normal day")
}

func TestAbsentRequiredPropertiesFails(t *testing.T) {
	yaml := `
summary: Just a normal day
`
	w, errs := Parse(yaml)
	assert.Equal(t, w, nil)
	assert.Contains(t, errs, errors.New("INVALID_DATE"))
}

func TestMalformedDateFails(t *testing.T) {
	yaml := `
date: 01.01.2020
`
	w, errs := Parse(yaml)
	assert.Equal(t, w, nil)
	assert.Contains(t, errs, errors.New("INVALID_DATE"))
}

func TestFailOnUnknownProperties(t *testing.T) {
	yaml := `
date: 2020-01-01
foo: 1
bar: test
`
	w, errs := Parse(yaml)
	assert.Equal(t, w, nil)
	assert.Contains(t, errs, errors.New("MALFORMED_YAML"))
}

func TestParseWithMalformedTimesFails(t *testing.T) {
	yaml := `
date: 1985-03-14
hours:
- time: asdf
`
	w, errs := Parse(yaml)
	assert.Equal(t, w, nil)
	assert.Contains(t, errs, errors.New("INVALID_TIME"))
}

func TestParseWithInvalidTimesFails(t *testing.T) {
	yaml := `
date: 1985-03-14
hours:
- time: 23:87
`
	w, errs := Parse(yaml)
	assert.Equal(t, w, nil)
	assert.Contains(t, errs, errors.New("INVALID_TIME"))
}
