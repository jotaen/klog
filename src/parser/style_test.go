package parser

import (
	. "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestDefaultStyle(t *testing.T) {
	assert.Equal(t, &Style{
		LineEnding:     styleProp[string]{"\n", false},
		Indentation:    styleProp[string]{"    ", false},
		DateFormat:     styleProp[DateFormat]{DateFormat{UseDashes: true}, false},
		TimeFormat:     styleProp[TimeFormat]{TimeFormat{Use24HourClock: true}, false},
		SpacingInRange: styleProp[string]{" ", false},
	}, DefaultStyle())
}

func TestDetectsStyleFromMinimalFile(t *testing.T) {
	rs := parseOrPanic("2000-01-01")
	assert.Equal(t, &Style{
		LineEnding:     styleProp[string]{"\n", false},
		Indentation:    styleProp[string]{"    ", false},
		DateFormat:     styleProp[DateFormat]{DateFormat{UseDashes: true}, true},
		TimeFormat:     styleProp[TimeFormat]{TimeFormat{Use24HourClock: true}, false},
		SpacingInRange: styleProp[string]{" ", false},
	}, rs[0].Style)
}

func TestDetectCanonicalStyle(t *testing.T) {
	rs := parseOrPanic("2000-01-01\nTest\n    8:00 - 9:00\n")
	assert.Equal(t, &Style{
		LineEnding:     styleProp[string]{"\n", true},
		Indentation:    styleProp[string]{"    ", true},
		DateFormat:     styleProp[DateFormat]{DateFormat{UseDashes: true}, true},
		TimeFormat:     styleProp[TimeFormat]{TimeFormat{Use24HourClock: true}, true},
		SpacingInRange: styleProp[string]{" ", true},
	}, rs[0].Style)
}

func TestDetectsCustomStyle(t *testing.T) {
	rs := parseOrPanic("2000/01/01\r\nTest\r\n\t8:00am-9:00am\r\n")
	assert.Equal(t, &Style{
		LineEnding:     styleProp[string]{"\r\n", true},
		Indentation:    styleProp[string]{"\t", true},
		DateFormat:     styleProp[DateFormat]{DateFormat{UseDashes: false}, true},
		TimeFormat:     styleProp[TimeFormat]{TimeFormat{Use24HourClock: false}, true},
		SpacingInRange: styleProp[string]{"", true},
	}, rs[0].Style)
}

func TestElectStyle(t *testing.T) {
	rs := parseOrPanic(
		"2001-05-19\n\t1:00 - 2:00\n\n",
		"2001/05/19\r\n  1:00am-2:00pm\r\n\r\n",
		"2001-05-19\n   1:00am-2:00pm\n   2:00pm-3:00pm\n\n",
		"2001/05/19\r\n  1:00 - 2:00\r\n\r\n",
		"2001-05-19\r\n    1:00am-2:00pm\r\n\r\n",
	)
	result := Elect(*DefaultStyle(), rs)
	assert.Equal(t, &Style{
		LineEnding:     styleProp[string]{"\r\n", true},
		Indentation:    styleProp[string]{"  ", true},
		DateFormat:     styleProp[DateFormat]{DateFormat{UseDashes: true}, true},
		TimeFormat:     styleProp[TimeFormat]{TimeFormat{Use24HourClock: false}, true},
		SpacingInRange: styleProp[string]{"", true},
	}, result)
}

func TestElectStyleDoesNotOverrideSetPreferences(t *testing.T) {
	rs := parseOrPanic(
		"2001-05-19\n\t1:00 - 2:00\n\n",
		"2001/05/19\r\n  1:00am-2:00pm\r\n\r\n",
		"2001-05-19\n   1:00am-2:00pm\n   2:00pm-3:00pm\n\n",
		"2001/05/19\r\n  1:00 - 2:00\r\n\r\n",
		"2001-05-19\r\n    1:00am-2:00pm\r\n\r\n",
	)
	result := Elect(*parseOrPanic("2018/01/01\n\t8:00 - 9:00")[0].Style, rs)
	assert.Equal(t, &Style{
		LineEnding:     styleProp[string]{"\n", true},
		Indentation:    styleProp[string]{"\t", true},
		DateFormat:     styleProp[DateFormat]{DateFormat{UseDashes: false}, true},
		TimeFormat:     styleProp[TimeFormat]{TimeFormat{Use24HourClock: true}, true},
		SpacingInRange: styleProp[string]{" ", true},
	}, result)
}

func parseOrPanic(recordsAsText ...string) []ParsedRecord {
	rs, err := Parse(strings.Join(recordsAsText, ""))
	if err != nil {
		panic("Invalid data")
	}
	return rs
}
