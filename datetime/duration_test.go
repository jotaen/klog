package datetime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerialiseDurationOnlyWithMeaningfulValues(t *testing.T) {
	assert.Equal(t, "1m", Duration(1).ToString())
	assert.Equal(t, "15h", Duration(15*60).ToString())
}

func TestSerialiseDurationOfLargeValues(t *testing.T) {
	assert.Equal(t, "265h 45m", Duration(265*60+45).ToString())
}

func TestSerialiseDurationWithoutLeadingZeros(t *testing.T) {
	assert.Equal(t, "2h 6m", Duration(2*60+6).ToString())
}

func TestSerialiseDurationOfNegativeValues(t *testing.T) {
	assert.Equal(t, "-3h 18m", Duration(-(3*60 + 18)).ToString())
	assert.Equal(t, "-3h", Duration(-(3 * 60)).ToString())
	assert.Equal(t, "-18m", Duration(-(18)).ToString())
}

func TestParsingDurationWithHoursAndMinutes(t *testing.T) {
	duration, err := CreateDurationFromString("2h 6m")
	assert.Nil(t, err)
	assert.Equal(t, Duration(2*60+6), duration)
}

func TestParsingDurationWithHoursOnly(t *testing.T) {
	for _, d := range []string{
		"13h",
		"13h 0m",
	} {
		duration, err := CreateDurationFromString(d)
		assert.Nil(t, err)
		assert.Equal(t, Duration(13*60), duration)
	}
}

func TestParsingDurationWithMinutesOnly(t *testing.T) {
	for _, d := range []string{
		"48m",
		"0h 48m",
	} {
		duration, err := CreateDurationFromString(d)
		assert.Nil(t, err)
		assert.Equal(t, Duration(48), duration)
	}
}

func TestParsingNegativeDuration(t *testing.T) {
	duration, err := CreateDurationFromString("-2h 5m")
	assert.Nil(t, err)
	assert.Equal(t, Duration(-(2*60 + 5)), duration)
}

func TestParsingIgnoresWhiteSpace(t *testing.T) {
	for _, d := range []string{
		"1h11m",
		"1h 11m",
		"  1h 11m  ",
		"1h     11m",
	} {
		duration, err := CreateDurationFromString(d)
		assert.Nil(t, err)
		assert.Equal(t, Duration(71), duration)
	}
}

func TestParsingFailsWithInvalidValue(t *testing.T) {
	for _, d := range []string{
		"asdf",
		"6h asdf",
		"qwer 30m",
	} {
		duration, err := CreateDurationFromString(d)
		assert.Error(t, err)
		assert.Equal(t, Duration(0), duration)
	}
}

func TestParsingFailsWithMinutesGreaterThan60(t *testing.T) {
	duration, err := CreateDurationFromString("8h 1653m")
	assert.Error(t, err)
	assert.Equal(t, Duration(0), duration)
}
