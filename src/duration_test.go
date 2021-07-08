package klog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerialiseDurationOnlyWithMeaningfulValues(t *testing.T) {
	assert.Equal(t, "0m", NewDuration(0, 0).ToString())
	assert.Equal(t, "1m", NewDuration(0, 1).ToString())
	assert.Equal(t, "15h", NewDuration(15, 0).ToString())
}

func TestSerialiseDurationOfLargeHourValues(t *testing.T) {
	assert.Equal(t, "265h45m", NewDuration(265, 45).ToString())
}

func TestSerialiseDurationWithoutLeadingZeros(t *testing.T) {
	assert.Equal(t, "2h6m", NewDuration(2, 6).ToString())
}

func TestSerialiseDurationOfNegativeValues(t *testing.T) {
	assert.Equal(t, "-3h18m", NewDuration(-3, -18).ToString())
	assert.Equal(t, "-3h", NewDuration(-3, 0).ToString())
	assert.Equal(t, "-18m", NewDuration(0, -18).ToString())
}

func TestSerialiseDurationOfExplicitlyPositiveValues(t *testing.T) {
	assert.Equal(t, "+3h18m", NewDuration(3, 18).ToStringWithSign())
	assert.Equal(t, "+3h", NewDuration(3, 0).ToStringWithSign())
	assert.Equal(t, "+18m", NewDuration(0, 18).ToStringWithSign())

	// 0 is an exception, as it doesn’t make sense to sign it
	assert.Equal(t, "0m", NewDuration(0, 0).ToStringWithSign())
}

func TestNormaliseDurationsWhenSerialising(t *testing.T) {
	assert.Equal(t, "2h", NewDuration(0, 120).ToString())
	assert.Equal(t, "2h30m", NewDuration(0, 150).ToString())
}

func TestParsingDurationWithHoursAndMinutes(t *testing.T) {
	duration, err := NewDurationFromString("2h6m")
	assert.Nil(t, err)
	assert.Equal(t, NewDuration(2, 6), duration)
}

func TestParsingDurationWithHourValueOnly(t *testing.T) {
	for _, d := range []string{
		"13h",
		"13h0m",
	} {
		duration, err := NewDurationFromString(d)
		assert.Nil(t, err)
		assert.Equal(t, NewDuration(13, 0), duration)
	}
}

func TestParsingDurationWithMinuteValueOnly(t *testing.T) {
	for _, d := range []struct {
		text   string
		expect Duration
	}{
		{"48m", NewDuration(0, 48)},
		{"0h48m", NewDuration(0, 48)},

		// Minutes >60 are okay if there is no hour part present
		{"120m", NewDuration(2, 0)},
		{"150m", NewDuration(2, 30)},
	} {
		duration, err := NewDurationFromString(d.text)
		assert.Nil(t, err)
		assert.Equal(t, d.expect, duration)
	}
}

func TestParsingNegativeDuration(t *testing.T) {
	duration, err := NewDurationFromString("-2h5m")
	assert.Nil(t, err)
	assert.Equal(t, NewDuration(-2, -5), duration)
}

func TestParsingExplicitlyPositiveDuration(t *testing.T) {
	duration, err := NewDurationFromString("+2h5m")
	assert.Nil(t, err)
	assert.Equal(t, NewDuration(2, 5), duration)
}

func TestParsingWithLeadingZeros(t *testing.T) {
	for _, d := range []string{
		"000009h00000000001m",
		"9h001m",
		"09h1m",
	} {
		duration, err := NewDurationFromString(d)
		assert.Nil(t, err)
		assert.Equal(t, NewDuration(9, 1), duration)
	}
}

func TestParsingFailsWithInvalidValue(t *testing.T) {
	for _, d := range []string{
		"",
		"1h 11m",
		"asdf",
		"6h asdf",
		"qwer 30m",
	} {
		duration, err := NewDurationFromString(d)
		assert.EqualError(t, err, "MALFORMED_DURATION")
		assert.Equal(t, nil, duration)
	}
}

func TestParsingFailsWithMinutesGreaterThan60WhenHourPartPresent(t *testing.T) {
	for _, d := range []string{
		"1h60m",
		"8h1653m",
		"-8h1653m",
	} {
		duration, err := NewDurationFromString(d)
		assert.EqualError(t, err, "UNREPRESENTABLE_DURATION")
		assert.Equal(t, nil, duration)
	}
}
