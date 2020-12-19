package datetime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerialiseDurationOnlyWithMeaningfulValues(t *testing.T) {
	assert.Equal(t, "0m", NewDuration(0, 0).ToString())
	assert.Equal(t, "1m", NewDuration(0, 1).ToString())
	assert.Equal(t, "15h", NewDuration(15, 0).ToString())
}

func TestSerialiseDurationOfLargeValues(t *testing.T) {
	assert.Equal(t, "265h 45m", NewDuration(265, 45).ToString())
}

func TestSerialiseDurationWithoutLeadingZeros(t *testing.T) {
	assert.Equal(t, "2h 6m", NewDuration(2, 6).ToString())
}

func TestSerialiseDurationOfNegativeValues(t *testing.T) {
	assert.Equal(t, "-3h 18m", NewDuration(-3, -18).ToString())
	assert.Equal(t, "-3h", NewDuration(-3, 0).ToString())
	assert.Equal(t, "-18m", NewDuration(0, -18).ToString())
}

func TestParsingDurationWithHoursAndMinutes(t *testing.T) {
	duration, err := NewDurationFromString("2h 6m")
	assert.Nil(t, err)
	assert.Equal(t, NewDuration(2, 6), duration)
}

func TestParsingDurationWithHoursOnly(t *testing.T) {
	for _, d := range []string{
		"13h",
		"13h 0m",
	} {
		duration, err := NewDurationFromString(d)
		assert.Nil(t, err)
		assert.Equal(t, NewDuration(13, 0), duration)
	}
}

func TestParsingDurationWithMinutesOnly(t *testing.T) {
	for _, d := range []string{
		"48m",
		"0h 48m",
	} {
		duration, err := NewDurationFromString(d)
		assert.Nil(t, err)
		assert.Equal(t, NewDuration(0, 48), duration)
	}
}

func TestParsingNegativeDuration(t *testing.T) {
	duration, err := NewDurationFromString("-2h 5m")
	assert.Nil(t, err)
	assert.Equal(t, NewDuration(-2, -5), duration)
}

func TestParsingIgnoresWhiteSpace(t *testing.T) {
	for _, d := range []string{
		"1h11m",
		"1h 11m",
		"  1h 11m  ",
		"1h     11m",
	} {
		duration, err := NewDurationFromString(d)
		assert.Nil(t, err)
		assert.Equal(t, NewDuration(0, 71), duration)
	}
}

func TestParsingFailsWithInvalidValue(t *testing.T) {
	for _, d := range []string{
		"asdf",
		"6h asdf",
		"qwer 30m",
	} {
		duration, err := NewDurationFromString(d)
		assert.EqualError(t, err, "MALFORMED_DURATION")
		assert.Equal(t, nil, duration)
	}
}

func TestParsingFailsWithMinutesGreaterThan60(t *testing.T) {
	duration, err := NewDurationFromString("8h 1653m")
	assert.EqualError(t, err, "UNREPRESENTABLE_DURATION")
	assert.Equal(t, nil, duration)
}
