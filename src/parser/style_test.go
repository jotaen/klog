package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultStyle(t *testing.T) {
	assert.Equal(t, Style{
		LineEnding:       "\n",
		Indentation:      "    ",
		UsesDashesInDate: true,
		Uses24HourClock:  true,
		UsesSpaceInRange: true,
	}, DefaultStyle())
}

func TestDetectsStyleFromMinimalFile(t *testing.T) {
	pRecord, _ := Parse("2000-01-01")
	assert.Equal(t, DefaultStyle(), pRecord[0].Style)
}

func TestDetectCanonicalStyle(t *testing.T) {
	pRecord, _ := Parse("2000-01-01\nTest\n    8:00 - 9:00\n")
	assert.Equal(t, DefaultStyle(), pRecord[0].Style)
}

func TestDetectsCustomStyle(t *testing.T) {
	pRecord, _ := Parse("2000/01/01\r\nTest\r\n\t8:00am-9:00am\r\n")
	assert.Equal(t, Style{
		LineEnding:       "\r\n",
		Indentation:      "\t",
		UsesDashesInDate: false,
		Uses24HourClock:  false,
		UsesSpaceInRange: false,
	}, pRecord[0].Style)
}

func TestCondense(t *testing.T) {

}
