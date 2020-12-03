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

func TestParseEntryWithTime(t *testing.T) {
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
