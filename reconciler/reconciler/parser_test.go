package reconciler

import (
	"cloud.google.com/go/civil"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseSimpleEntry(t *testing.T) {
	yaml := `
date: 2020-01-01
summary: Just a normal day
`
	entry, _ := Parse(yaml)
	assert.Equal(t, entry, Entry{
		Date:    civil.Date{Year: 2020, Month: time.January, Day: 1},
		Summary: "Just a normal day",
	})
}

func TestAbsentDatePropertyFails(t *testing.T) {
	yaml := `
summary: Just a normal day
`
	entry, err := Parse(yaml)
	assert.Equal(t, entry, Entry{})
	assert.Error(t, err)
}

func TestMinimalValidEntry(t *testing.T) {
	yaml := `
date: 2020-01-01
`
	entry, err := Parse(yaml)
	assert.Equal(t, entry.Summary, "")
	assert.Equal(t, err, nil)
}

func TestFailOnUnknownProperties(t *testing.T) {
	yaml := `
date: 2020-01-01
foo: 1
bar: test
`
	entry, err := Parse(yaml)
	assert.Equal(t, entry, Entry{})
	assert.Error(t, err)
}

func TestParseEntryWithTimes(t *testing.T) {
	yaml := `
date: 1985-03-14
hours:
- time: 2:00
- time: 5:00
`
	entry, _ := Parse(yaml)
	assert.Equal(t, entry, Entry{
		Date:  civil.Date{Year: 1985, Month: time.March, Day: 14},
		Times: []Minutes{Minutes(2 * 60), Minutes(5 * 60)},
	})
}

func TestParseEntryWithMalformedTimesFails(t *testing.T) {
	yaml := `
date: 1985-03-14
hours:
- time: asdf
`
	entry, err := Parse(yaml)
	assert.Equal(t, entry, Entry{})
	assert.Error(t, err)
}

func TestParseEntryWithInvalidTimesFails(t *testing.T) {
	yaml := `
date: 1985-03-14
hours:
- time: 23:87
`
	entry, err := Parse(yaml)
	assert.Equal(t, entry, Entry{})
	assert.Error(t, err)
}

func TestParseEntryWithRanges(t *testing.T) {
	yaml := `
date: 2008-12-03
hours:
- start: 9:12
  end: 12:05
- start: 13:03
  end: 16:47
`
	entry, _ := Parse(yaml)
	assert.Equal(t, entry, Entry{
		Date: civil.Date{Year: 2008, Month: time.December, Day: 3},
		Ranges: []Range{
			Range{
				Start: civil.Time{Hour: 9, Minute: 12},
				End:   civil.Time{Hour: 12, Minute: 5},
			},
			Range{
				Start: civil.Time{Hour: 13, Minute: 3},
				End:   civil.Time{Hour: 16, Minute: 47},
			},
		},
	})
}
