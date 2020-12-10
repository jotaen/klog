package parser

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"klog/datetime"
	"klog/testutil"
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
- start: 8:12
  end: 09:05
- start: 10:15
- time: 2h
- time: 05h 03m
- time: -4h 45m
`
	range1 := testutil.Range_(testutil.Time_(8, 12), testutil.Time_(9, 05))

	w, errs := Parse(yaml)
	require.Equal(t, 0, len(errs))

	assert.Equal(t, w.Ranges(), []datetime.TimeRange{range1})
	assert.Equal(t, w.OpenRangeStart(), testutil.Time_(10, 15))
	assert.Equal(t, w.Times(), []datetime.Duration{
		datetime.Duration(2 * 60),
		datetime.Duration(5*60 + 3),
		datetime.Duration(-(4*60 + 45)),
	})
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

func TestMalformedValuesFails(t *testing.T) {
	yaml := `
date: 01.01.2020
hours:
- time: asdf
`
	w, errs := Parse(yaml)
	assert.Equal(t, w, nil)
	assert.Contains(t, errs, errors.New("INVALID_DATE"))
}
