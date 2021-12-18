package parser

import (
	klog "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestDefaultStyle(t *testing.T) {
	assert.Equal(t, &Style{
		LineEnding:     "\n",
		Indentation:    "    ",
		DateFormat:     klog.DateFormat{UseDashes: true},
		TimeFormat:     klog.TimeFormat{Is24HourClock: true},
		SpacingInRange: " ",
	}, DefaultStyle())
}

func TestDetectsStyleFromMinimalFile(t *testing.T) {
	rs := parseOrPanic("2000-01-01")
	assert.Equal(t, &Style{
		LineEnding:     "\n",
		Indentation:    "    ",
		DateFormat:     klog.DateFormat{UseDashes: true},
		dateFormatSet:  true,
		TimeFormat:     klog.TimeFormat{Is24HourClock: true},
		SpacingInRange: " ",
	}, rs[0].Style)
}

func TestDetectCanonicalStyle(t *testing.T) {
	rs := parseOrPanic("2000-01-01\nTest\n    8:00 - 9:00\n")
	assert.Equal(t, &Style{
		LineEnding:        "\n",
		lineEndingSet:     true,
		Indentation:       "    ",
		indentationSet:    true,
		SpacingInRange:    " ",
		spacingInRangeSet: true,
		DateFormat:        klog.DateFormat{UseDashes: true},
		dateFormatSet:     true,
		TimeFormat:        klog.TimeFormat{Is24HourClock: true},
		timeFormatSet:     true,
	}, rs[0].Style)
}

func TestDetectsCustomStyle(t *testing.T) {
	rs := parseOrPanic("2000/01/01\r\nTest\r\n\t8:00am-9:00am\r\n")
	assert.Equal(t, &Style{
		LineEnding:        "\r\n",
		lineEndingSet:     true,
		Indentation:       "\t",
		indentationSet:    true,
		SpacingInRange:    "",
		spacingInRangeSet: true,
		DateFormat:        klog.DateFormat{UseDashes: false},
		dateFormatSet:     true,
		TimeFormat:        klog.TimeFormat{Is24HourClock: false},
		timeFormatSet:     true,
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
		LineEnding:        "\r\n",
		lineEndingSet:     true,
		Indentation:       "  ",
		indentationSet:    true,
		SpacingInRange:    "",
		spacingInRangeSet: true,
		DateFormat:        klog.DateFormat{UseDashes: true},
		dateFormatSet:     true,
		TimeFormat:        klog.TimeFormat{Is24HourClock: false},
		timeFormatSet:     true,
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
		LineEnding:        "\n",
		lineEndingSet:     true,
		Indentation:       "\t",
		indentationSet:    true,
		SpacingInRange:    " ",
		spacingInRangeSet: true,
		DateFormat:        klog.DateFormat{UseDashes: false},
		dateFormatSet:     true,
		TimeFormat:        klog.TimeFormat{Is24HourClock: true},
		timeFormatSet:     true,
	}, result)
}

func parseOrPanic(recordsAsText ...string) []ParsedRecord {
	rs, err := Parse(strings.Join(recordsAsText, ""))
	if err != nil {
		panic("Invalid data")
	}
	return rs
}
