package reconciler

import (
	"testing"
	"cloud.google.com/go/civil"
	"github.com/stretchr/testify/assert"
)

func TestParseSimpleEntry(t *testing.T) {
	yaml := `
date: 2020-01-01
summary: Just a normal day
`
	entry, _ := Parse(yaml)
	date, _ := civil.ParseDate("2020-01-01")
	assert.Equal(t, entry, Entry{
		Date: date,
		Summary: "Just a normal day",
	})
}

func TestAbsentDatePropertyFails(t *testing.T) {
	yaml := `
summary: Just a normal day
`
	_, err := Parse(yaml)
	assert.Error(t, err)
}

func TestSummaryIsOptional(t *testing.T) {
	yaml := `
date: 2020-01-01
`
	entry, err := Parse(yaml)
	assert.Equal(t, entry.Summary, "")
	assert.Equal(t, err, nil)
}

func TestParseEntryWithTimes(t *testing.T) {
	yaml := `
date: 2020-01-01
hours:
- time: 2:00
- time: 5:00
`
	entry, _ := Parse(yaml)
	date, _ := civil.ParseDate("2020-01-01")
	assert.Equal(t, entry, Entry{
		Date: date,
		Times: []Minutes{ Minutes(2*60), Minutes(5*60) },
	})
}

func TestParseEntryWithRanges(t *testing.T) {
	yaml := `
date: 2020-01-01
hours:
- start: 9:12
  end: 12:00
- start: 13:00
  end: 16:30
`
	entry, _ := Parse(yaml)
	date, _ := civil.ParseDate("2020-01-01")
	x1_start, _ := civil.ParseTime("9:12:00")
	x1_end, _ := civil.ParseTime("12:00:00")
	x2_start, _ := civil.ParseTime("13:00:00")
	x2_end, _ := civil.ParseTime("16:30:00")
	assert.Equal(t, entry, Entry{
		Date: date,
		Ranges: []Range{
			Range{ Start: x1_start, End: x1_end },
			Range{ Start: x2_start, End: x2_end },
		},
	})
}
