package datetime

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOnlyConstructsValidTimes(t *testing.T) {
	tm, err := NewTime(22, 13)
	require.Nil(t, err)
	assert.Equal(t, tm.Hour(), 22)
	assert.Equal(t, tm.Minute(), 13)
	assert.Equal(t, tm.IsToday(), true)
	assert.Equal(t, tm.IsYesterday(), false)
	assert.Equal(t, tm.IsTomorrow(), false)
}

func TestDetectsInvalidTimes(t *testing.T) {
	invalidHour, err := NewTime(25, 30)
	assert.EqualError(t, err, "INVALID_TIME")
	assert.Nil(t, invalidHour)

	invalidMinute, err := NewTime(4, 85)
	assert.EqualError(t, err, "INVALID_TIME")
	assert.Nil(t, invalidMinute)

	invalidTime, err := NewTime(24, 00)
	assert.EqualError(t, err, "INVALID_TIME")
	assert.Nil(t, invalidTime)
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
	require.Nil(t, err)
	should, _ := NewTime(9, 42)
	assert.Equal(t, should, tm)
}

func TestParseTimeYesterday(t *testing.T) {
	tm, err := NewTimeFromString("22:43 yesterday")
	require.Nil(t, err)
	should, _ := NewTimeYesterday(22, 43)
	assert.Equal(t, should, tm)
	assert.Equal(t, false, tm.IsToday())
	assert.Equal(t, true, tm.IsYesterday())
	assert.Equal(t, false, tm.IsTomorrow())
}

func TestParseTimeTomorrow(t *testing.T) {
	tm, err := NewTimeFromString("02:12 tomorrow")
	require.Nil(t, err)
	should, _ := NewTimeTomorrow(2, 12)
	assert.Equal(t, should, tm)
	assert.Equal(t, false, tm.IsToday())
	assert.Equal(t, false, tm.IsYesterday())
	assert.Equal(t, true, tm.IsTomorrow())
}

func TestParseTimeFailsfMalformed(t *testing.T) {
	for _, s := range []string{
		"009:42",
		"asdf",
		"12",
		"13:3",
	} {
		tm, err := NewTimeFromString(s)
		require.Nil(t, tm)
		assert.EqualError(t, err, "MALFORMED_TIME")
	}
}

func TestCalculateMinutesSinceMidnight(t *testing.T) {
	for _, s := range []struct {
		in  string
		exp Duration
	}{
		{in: "0:00", exp: NewDuration(0, 0)},
		{in: "0:01", exp: NewDuration(0, 1)},
		{in: "14:59", exp: NewDuration(14, 59)},
		{in: "23:59", exp: NewDuration(23, 59)},
		{in: "18:35 yesterday", exp: NewDuration(-5, -25)},
		{in: "5:35 tomorrow", exp: NewDuration(24+5, 35)},
	} {
		tm, err := NewTimeFromString(s.in)
		require.Nil(t, err)
		assert.Equal(t, s.exp, tm.MidnightOffset())
	}
}

func TestTimeComparison(t *testing.T) {
	today1, _ := NewTime(12, 30)
	today2, _ := NewTime(12, 31)
	yesterday, _ := NewTimeYesterday(22, 43)
	tomorrow, _ := NewTimeTomorrow(9, 50)

	assert.True(t, today1.IsAfterOrEqual(today1))
	assert.True(t, today2.IsAfterOrEqual(today1))
	assert.True(t, today1.IsAfterOrEqual(yesterday))
	assert.True(t, tomorrow.IsAfterOrEqual(today1))
}
