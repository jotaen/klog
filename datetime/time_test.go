package datetime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOnlyConstructsValidTimes(t *testing.T) {
	tm, err := NewTime(22, 13)
	assert.Equal(t, tm.Hour(), 22)
	assert.Equal(t, tm.Minute(), 13)
	assert.Nil(t, err)
}

func TestDetectsInvalidTimes(t *testing.T) {
	invalidHour, err := NewTime(25, 30)
	assert.Error(t, err)
	assert.Nil(t, invalidHour)

	invalidMinute, err := NewTime(4, 85)
	assert.Error(t, err)
	assert.Nil(t, invalidMinute)
}

func TestSerialiseTime(t *testing.T) {
	tm, _ := NewTime(13, 45)
	assert.Equal(t, "13:45", tm.ToString())
}

func TestSerialiseTimeWithoutLeadingZeros(t *testing.T) {
	tm, _ := NewTime(8, 5)
	assert.Equal(t, "8:05", tm.ToString())
}

func TestParseTime(t *testing.T) {
	tm, err := NewTimeFromString("9:42")
	assert.Nil(t, err)
	should, _ := NewTime(9, 42)
	assert.Equal(t, tm, should)
}

func TestParseTimeFailsIfMalformed(t *testing.T) {
	for _, s := range []string{
		"009:42",
		"asdf",
		"12",
		"13:3",
	} {
		tm, err := NewTimeFromString(s)
		assert.Nil(t, tm)
		assert.EqualError(t, err, "MALFORMED_TIME")
	}
}
