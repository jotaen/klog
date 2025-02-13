package klog

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSerialiseDurationOnlyWithMeaningfulValues(t *testing.T) {
	assert.Equal(t, "0m", NewDuration(0, 0).ToString())
	assert.Equal(t, "1m", NewDuration(0, 1).ToString())
	assert.Equal(t, "34m", NewDuration(0, 34).ToString())
	assert.Equal(t, "59m", NewDuration(0, 59).ToString())
	assert.Equal(t, "1h", NewDuration(1, 0).ToString())
	assert.Equal(t, "15h", NewDuration(15, 0).ToString())
	assert.Equal(t, "15h3m", NewDuration(15, 3).ToString())
	assert.Equal(t, "265h45m", NewDuration(265, 45).ToString())
	assert.Equal(t, "4716278h48m", NewDuration(4716278, 48).ToString())
	assert.Equal(t, "153722867280912930h7m", NewDuration(0, 9223372036854775807).ToString())
}

func TestSerialiseDurationWithoutLeadingZeros(t *testing.T) {
	assert.Equal(t, "2h6m", NewDuration(2, 6).ToString())
}

func TestSerialiseDurationOfNegativeValues(t *testing.T) {
	assert.Equal(t, "-2h4m", NewDuration(-2, -4).ToString())
	assert.Equal(t, "-3h18m", NewDuration(-3, -18).ToString())
	assert.Equal(t, "-812747h", NewDuration(-812747, 0).ToString())
	assert.Equal(t, "-18m", NewDuration(0, -18).ToString())
	assert.Equal(t, "-153722867280912930h7m", NewDuration(0, -9223372036854775807).ToString())
}

func TestSerialiseDurationWithSign(t *testing.T) {
	// Zero is neutral by default:
	assert.Equal(t, "0m", NewDuration(0, 0).ToStringWithSign())

	// Positive values:
	assert.Equal(t, "+3h18m", NewDuration(3, 18).ToStringWithSign())
	assert.Equal(t, "+3h", NewDuration(3, 0).ToStringWithSign())
	assert.Equal(t, "+18m", NewDuration(0, 18).ToStringWithSign())

	// Negative values:
	assert.Equal(t, "-3h18m", NewDuration(-3, -18).ToStringWithSign())
	assert.Equal(t, "-3h", NewDuration(-3, 0).ToStringWithSign())
	assert.Equal(t, "-18m", NewDuration(0, -18).ToStringWithSign())
}

func TestSerialisePreservesOriginalFormatting(t *testing.T) {
	for _, x := range []string{
		"0m",
		"+0m",
		"-0m",

		"15m",
		"+15m",
		"-15m",
	} {
		neutralZero, _ := NewDurationFromString(x)
		assert.Equal(t, x, neutralZero.ToString())
	}
}

func TestNormaliseDurationsWhenSerialising(t *testing.T) {
	assert.Equal(t, "2h", NewDuration(0, 120).ToString())
	assert.Equal(t, "2h30m", NewDuration(0, 150).ToString())

	d, _ := NewDurationFromString("120m")
	assert.Equal(t, "2h", d.ToString())
}

func TestParsingDurationWithHoursAndMinutes(t *testing.T) {
	d, err := NewDurationFromString("2h6m")
	assert.Nil(t, err)
	assert.Equal(t, NewDuration(2, 6), d)
}

func TestParsingDurationWithHourValueOnly(t *testing.T) {
	for _, d := range []struct {
		text   string
		expect Duration
	}{
		{"0h", NewDuration(0, 0)},
		{"1h", NewDuration(1, 0)},
		{"13h", NewDuration(13, 0)},
		{"9882187612h", NewDuration(9882187612, 0)},
		{"13h0m", NewDuration(13, 0)},
	} {
		duration, err := NewDurationFromString(d.text)
		assert.Nil(t, err)
		assert.Equal(t, d.expect, duration)
	}
}

func TestParsingDurationWithMinuteValueOnly(t *testing.T) {
	for _, d := range []struct {
		text   string
		expect Duration
	}{
		{"1m", NewDuration(0, 1)},
		{"48m", NewDuration(0, 48)},
		{"59m", NewDuration(0, 59)},

		{"0h48m", NewDuration(0, 48)},

		// Minutes >60 are okay if there is no hour part present
		{"60m", NewDuration(1, 0)},
		{"120m", NewDuration(2, 0)},
		{"568721940327m", NewDuration(0, 568721940327)},
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
	assert.Equal(t, NewDurationWithFormat(2, 5, DurationFormat{ForcePlus: true}), duration)
	assert.Equal(t, "+2h5m", duration.ToString())
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
		"⠙⠛m",   // Braille digits
		"四二h", // Japanese digits
		"᠒h᠐᠒m", // Mongolean digits
	} {
		duration, err := NewDurationFromString(d)
		assert.EqualError(t, err, "MALFORMED_DURATION")
		assert.Equal(t, nil, duration)
	}
}

func TestParsingFailsWithMinutesGreaterThan60WhenHourPartPresent(t *testing.T) {
	for _, d := range []string{
		"1h60m",
		"0h60m",
		"8h1653m",
		"-8h1653m",
	} {
		duration, err := NewDurationFromString(d)
		assert.EqualError(t, err, "UNREPRESENTABLE_DURATION")
		assert.Equal(t, nil, duration)
	}
}

func TestParsingDurationWithMaxValue(t *testing.T) {
	t.Run("max", func(t *testing.T) {
		d, err := NewDurationFromString("9223372036854775807m")
		require.Nil(t, err)
		assert.Equal(t, NewDuration(0, 9223372036854775807), d)
	})
	t.Run("max", func(t *testing.T) {
		d, err := NewDurationFromString("153722867280912930h7m")
		require.Nil(t, err)
		assert.Equal(t, NewDuration(153722867280912930, 7), d)
	})
	t.Run("min", func(t *testing.T) {
		d, err := NewDurationFromString("-9223372036854775807m")
		require.Nil(t, err)
		assert.Equal(t, NewDuration(0, -9223372036854775807), d)
	})
	t.Run("max", func(t *testing.T) {
		d, err := NewDurationFromString("-153722867280912930h7m")
		require.Nil(t, err)
		assert.Equal(t, NewDuration(-153722867280912930, -7), d)
	})
}

func TestParsingDurationTooBigToRepresent(t *testing.T) {
	for _, d := range []string{
		"9223372036854775808m",
		"-9223372036854775808m",
		"9223372036854775808h",
		"-9223372036854775808h",
		"153722867280912930h08m",
		"-153722867280912930h08m",
	} {
		assert.Panics(t, func() {
			_, _ = NewDurationFromString(d)
		}, d)
	}
}

func TestDurationPlusMinus(t *testing.T) {
	for _, d := range []struct {
		sum    Duration
		expect int
	}{
		{NewDuration(0, 0).Plus(NewDuration(0, 0)), 0},
		{NewDuration(0, 0).Plus(NewDuration(0, 1)), 1},
		{NewDuration(0, 0).Plus(NewDuration(1, 2)), 62},
		{NewDuration(1382, 9278).Plus(NewDuration(4718, 5010)), 380288},
		{NewDuration(0, 9223372036854775806).Plus(NewDuration(0, 1)), 9223372036854775807},
		{NewDuration(0, 0).Plus(NewDuration(0, -9223372036854775807)), -9223372036854775807},

		{NewDuration(0, 0).Minus(NewDuration(0, 0)), 0},
		{NewDuration(0, 0).Minus(NewDuration(0, 1)), -1},
		{NewDuration(0, 0).Minus(NewDuration(1, 2)), -62},
		{NewDuration(1382, 9278).Minus(NewDuration(4718, 5010)), -195892},
	} {
		assert.Equal(t, d.sum.InMinutes(), d.expect)
	}
}

func TestPanicsIfAdditionOverflows(t *testing.T) {
	assert.Panics(t, func() {
		NewDuration(0, 9223372036854775807).Plus(NewDuration(0, 1))
	})

	assert.Panics(t, func() {
		NewDuration(0, -9223372036854775807).Plus(NewDuration(0, -1))
	})
}
