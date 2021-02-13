package klog

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
	tm := Ɀ_Time_(13, 45)
	assert.Equal(t, "13:45", tm.ToString())
}

func TestSerialiseTimeWithoutLeadingZeros(t *testing.T) {
	tm := Ɀ_Time_(8, 5)
	assert.Equal(t, "8:05", tm.ToString())
}

func TestSerialiseTimeYesterday(t *testing.T) {
	tm := Ɀ_TimeYesterday_(23, 0)
	assert.Equal(t, "<23:00", tm.ToString())
}

func TestSerialiseTimeTomorrow(t *testing.T) {
	tm := Ɀ_TimeTomorrow_(0, 2)
	assert.Equal(t, "0:02>", tm.ToString())
}

func TestParseTime24Hours(t *testing.T) {
	tm, err := NewTimeFromString("9:42")
	require.Nil(t, err)
	should := Ɀ_Time_(9, 42)
	assert.Equal(t, should, tm)
}

func TestParseTime12Hours(t *testing.T) {
	for _, s := range []struct {
		val string
		exp Time
	}{
		{"12:37am", Ɀ_Time_(0, 37)},
		{"1:00am", Ɀ_Time_(1, 0)},
		{"1:00am", Ɀ_Time_(1, 0)},
		{"12:22pm", Ɀ_Time_(12, 22)},
		{"1:59pm", Ɀ_Time_(13, 59)},
		{"7:33pm", Ɀ_Time_(19, 33)},
	} {
		tm, err := NewTimeFromString(s.val)
		require.Nil(t, err)
		require.NotNil(t, tm)
		assert.True(t, s.exp.IsEqualTo(tm), s.val)
		assert.Equal(t, s.val, tm.ToString())
	}
}

func TestParseTimeYesterday(t *testing.T) {
	tm, err := NewTimeFromString("<22:43")
	require.Nil(t, err)
	should := Ɀ_TimeYesterday_(22, 43)
	assert.Equal(t, should, tm)
	assert.Equal(t, false, tm.IsToday())
	assert.Equal(t, true, tm.IsYesterday())
	assert.Equal(t, false, tm.IsTomorrow())
}

func TestParseTimeTomorrow(t *testing.T) {
	tm, err := NewTimeFromString("02:12>")
	require.Nil(t, err)
	should := Ɀ_TimeTomorrow_(2, 12)
	assert.Equal(t, should, tm)
	assert.Equal(t, false, tm.IsToday())
	assert.Equal(t, false, tm.IsYesterday())
	assert.Equal(t, true, tm.IsTomorrow())
}

func TestParseMalformedTimesFail(t *testing.T) {
	for _, s := range []string{
		"009:42", // Hours cannot have infinite leading 0s
		"09:042", // Minutes cannot have infinite leading 0s
		"<2:15>", // Markers cannot appear on both sides
		"asdf",
		"12",
		"13:3",   // Minutes must have 2 digits
		"-14:12", // Cannot be negative
		"14:-12", // Cannot be negative
	} {
		tm, err := NewTimeFromString(s)
		require.Nil(t, tm, s)
		assert.EqualError(t, err, "MALFORMED_TIME", s)
	}
}

func TestParseUnrepresentableTimesFail(t *testing.T) {
	for _, s := range []string{
		"25:12",
		"3:87",
		"00:00pm",
		"13:00am",
		"13:00pm",
	} {
		tm, err := NewTimeFromString(s)
		require.Nil(t, tm, s)
		assert.EqualError(t, err, "INVALID_TIME", s)
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
		{in: "<18:35", exp: NewDuration(-5, -25)},
		{in: "5:35>", exp: NewDuration(24+5, 35)},
	} {
		tm, err := NewTimeFromString(s.in)
		require.Nil(t, err)
		assert.Equal(t, s.exp, tm.MidnightOffset())
	}
}

func TestTimeComparison(t *testing.T) {
	today1 := Ɀ_Time_(12, 30)
	today2 := Ɀ_Time_(12, 31)
	yesterday := Ɀ_TimeYesterday_(22, 43)
	tomorrow := Ɀ_TimeTomorrow_(9, 50)

	assert.True(t, today1.IsAfterOrEqual(today1))
	assert.True(t, today2.IsAfterOrEqual(today1))
	assert.True(t, today1.IsAfterOrEqual(yesterday))
	assert.True(t, tomorrow.IsAfterOrEqual(today1))
}

func TestAddDuration(t *testing.T) {
	for _, x := range []struct {
		initial   Time
		increment Duration
		expect    Time
	}{
		{Ɀ_Time_(11, 30), NewDuration(0, 30), Ɀ_Time_(12, 00)},
		{Ɀ_Time_(18, 00), NewDuration(-6, 0), Ɀ_Time_(12, 00)},
		{Ɀ_Time_(3, 59), NewDuration(8, 1), Ɀ_Time_(12, 00)},
		{Ɀ_TimeYesterday_(23, 45), NewDuration(12, 15), Ɀ_Time_(12, 00)},
		{Ɀ_TimeYesterday_(12, 12), NewDuration(1, 19), Ɀ_TimeYesterday_(13, 31)},
		{Ɀ_TimeYesterday_(0, 1), NewDuration(0, -1), Ɀ_TimeYesterday_(0, 0)},
		{Ɀ_TimeTomorrow_(4, 12), NewDuration(-16, -12), Ɀ_Time_(12, 00)},
		{Ɀ_TimeTomorrow_(18, 38), NewDuration(-1, -1), Ɀ_TimeTomorrow_(17, 37)},
		{Ɀ_TimeTomorrow_(23, 58), NewDuration(0, 1), Ɀ_TimeTomorrow_(23, 59)},
	} {
		result, err := x.initial.Add(x.increment)
		require.Nil(t, err)
		assert.Equal(t, x.expect, result, x.initial)
	}
}

func TestAddDurationImpossible(t *testing.T) {
	for _, x := range []struct {
		initial   Time
		increment Duration
	}{
		{Ɀ_Time_(11, 30), NewDuration(353, 0)},
		{Ɀ_Time_(11, 30), NewDuration(-353, 0)},
		{Ɀ_TimeYesterday_(0, 0), NewDuration(0, -1)},
		{Ɀ_TimeTomorrow_(23, 59), NewDuration(0, 1)},
	} {
		result, err := x.initial.Add(x.increment)
		require.Nil(t, result)
		assert.Error(t, err)
	}
}
