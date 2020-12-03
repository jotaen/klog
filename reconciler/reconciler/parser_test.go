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
