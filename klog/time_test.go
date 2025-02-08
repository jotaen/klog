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

func TestHandle2400SpecialCase(t *testing.T) {
	{
		tm, err := NewTime(24, 00)
		require.Nil(t, err)
		assert.Equal(t, tm.Hour(), 0)
		assert.Equal(t, tm.Minute(), 0)
		assert.Equal(t, tm.IsToday(), false)
		assert.Equal(t, tm.IsYesterday(), false)
		assert.Equal(t, tm.IsTomorrow(), true)
	}
	{
		tm, err := NewTimeYesterday(24, 00)
		require.Nil(t, err)
		assert.Equal(t, tm.Hour(), 0)
		assert.Equal(t, tm.Minute(), 0)
		assert.Equal(t, tm.IsToday(), true)
		assert.Equal(t, tm.IsYesterday(), false)
		assert.Equal(t, tm.IsTomorrow(), false)
	}
	{
		// 24:00 tomorrow cannot be represented.
		tm, err := NewTimeTomorrow(24, 00)
		require.Nil(t, tm)
		require.Error(t, err)
	}
}

func TestDetectsInvalidTimes(t *testing.T) {
	for _, invalidTime := range []struct {
		hours   int
		minutes int
	}{
		// Invalid hours
		{24, 01},
		{25, 30},
		{124, 34},
		{-12, 34},

		// Invalid minutes
		{05, 60},
		{05, 61},
		{05, 245},
		{05, -12},

		// Both invalid
		{1575, 28293},
	} {
		for _, constructor := range []func(int, int) (Time, error){
			NewTime, NewTimeYesterday, NewTimeTomorrow,
		} {
			invalidTime, err := constructor(invalidTime.hours, invalidTime.minutes)
			assert.EqualError(t, err, "INVALID_TIME")
			assert.Nil(t, invalidTime)
		}
	}
}

func TestSerialiseTime(t *testing.T) {
	tm, err := NewTime(13, 45)
	require.Nil(t, err)
	assert.Equal(t, "13:45", tm.ToString())
	assert.Equal(t, "13:45", tm.ToStringWithFormat(TimeFormat{Use24HourClock: true}))
	assert.Equal(t, "1:45pm", tm.ToStringWithFormat(TimeFormat{Use24HourClock: false}))
}

func TestSerialiseTimeWithoutLeadingZeros(t *testing.T) {
	tm, err := NewTime(8, 5)
	require.Nil(t, err)
	assert.Equal(t, "8:05", tm.ToString())
	assert.Equal(t, "8:05am", Ɀ_IsAmPm_(tm).ToString())
}

func TestSerialiseTimeYesterday(t *testing.T) {
	tm, err := NewTimeYesterday(23, 0)
	require.Nil(t, err)
	assert.Equal(t, "<23:00", tm.ToString())
	assert.Equal(t, "<11:00pm", Ɀ_IsAmPm_(tm).ToString())
}

func TestSerialiseTimeTomorrow(t *testing.T) {
	tm, err := NewTimeTomorrow(0, 2)
	require.Nil(t, err)
	assert.Equal(t, "0:02>", tm.ToString())
	assert.Equal(t, "12:02am>", Ɀ_IsAmPm_(tm).ToString())
}

func TestParseTime24Hours(t *testing.T) {
	for _, s := range []struct {
		val string
		exp Time
	}{
		{"9:42", Ɀ_Time_(9, 42)},
		{"09:42", Ɀ_Time_(9, 42)},
		{"16:01", Ɀ_Time_(16, 01)},
	} {
		tm, err := NewTimeFromString(s.val)
		require.Nil(t, err)
		require.NotNil(t, tm)
		assert.Equal(t, s.exp, tm)
		assert.True(t, s.exp.IsEqualTo(tm), s.val)
		assert.Equal(t, TimeFormat{Use24HourClock: true}, tm.Format())
	}
}

func TestParseTime12Hours(t *testing.T) {
	for _, s := range []struct {
		val string
		exp Time
	}{
		{"12:00am", Ɀ_Time_(0, 00)},
		{"12:37am", Ɀ_Time_(0, 37)},
		{"1:00am", Ɀ_Time_(1, 0)},
		{"1:00am", Ɀ_Time_(1, 0)},
		{"12:00pm", Ɀ_Time_(12, 00)},
		{"12:22pm", Ɀ_Time_(12, 22)},
		{"1:59pm", Ɀ_Time_(13, 59)},
		{"7:33pm", Ɀ_Time_(19, 33)},
	} {
		tm, err := NewTimeFromString(s.val)
		require.Nil(t, err)
		require.NotNil(t, tm)
		assert.Equal(t, Ɀ_IsAmPm_(s.exp), tm)
		assert.True(t, s.exp.IsEqualTo(tm), s.val)
		assert.Equal(t, TimeFormat{Use24HourClock: false}, tm.Format())
	}
}

func TestParseTimeYesterday(t *testing.T) {
	for _, s := range []struct {
		val string
		exp Time
	}{
		{"<3:43", Ɀ_TimeYesterday_(3, 43)},
		{"<03:43", Ɀ_TimeYesterday_(3, 43)},
		{"<03:43am", Ɀ_IsAmPm_(Ɀ_TimeYesterday_(3, 43))},
		{"<3:43pm", Ɀ_IsAmPm_(Ɀ_TimeYesterday_(15, 43))},
	} {
		tm, err := NewTimeFromString(s.val)
		require.Nil(t, err)
		assert.Equal(t, s.exp, tm)
		assert.Equal(t, false, tm.IsToday())
		assert.Equal(t, true, tm.IsYesterday())
		assert.Equal(t, false, tm.IsTomorrow())
	}
}

func TestParseTimeTomorrow(t *testing.T) {
	for _, s := range []struct {
		val string
		exp Time
	}{
		{"2:12>", Ɀ_TimeTomorrow_(2, 12)},
		{"02:12>", Ɀ_TimeTomorrow_(2, 12)},
		{"2:12am>", Ɀ_IsAmPm_(Ɀ_TimeTomorrow_(2, 12))},
		{"02:12pm>", Ɀ_IsAmPm_(Ɀ_TimeTomorrow_(14, 12))},
	} {
		tm, err := NewTimeFromString(s.val)
		require.Nil(t, err)
		assert.Equal(t, s.exp, tm)
		assert.Equal(t, false, tm.IsToday())
		assert.Equal(t, false, tm.IsYesterday())
		assert.Equal(t, true, tm.IsTomorrow())
	}
}

func TestParseTime2400SpecialCase(t *testing.T) {
	for _, s := range []struct {
		val string
		exp Time
	}{
		{"<24:00", Ɀ_Time_(0, 0)},
		{"24:00", Ɀ_TimeTomorrow_(0, 0)},
	} {
		tm, err := NewTimeFromString(s.val)
		require.Nil(t, err)
		require.NotNil(t, tm)
		assert.True(t, s.exp.IsEqualTo(tm), s.val)
	}
}

func TestParseMalformedTimesFail(t *testing.T) {
	for _, s := range []string{
		"009:42", // Hours cannot have infinite leading 0s
		"09:042", // Minutes cannot have infinite leading 0s
		"<2:15>", // Shift-markers cannot appear on both sides
		"asdf",
		"12",
		"12am",   // Minutes missing
		"13:3",   // Minutes must have 2 digits
		"-14:12", // Hours cannot be negative
		"14:-12", // Minutes cannot be negative
		"⠃⠚:⠙⠛",  // Braille digits
		"四:二八",   // Japanese digits
		"᠒᠐:᠑᠒",  // Mongolean digits
	} {
		tm, err := NewTimeFromString(s)
		require.Nil(t, tm, s)
		assert.EqualError(t, err, "MALFORMED_TIME", s)
	}
}

func TestParseUnrepresentableTimesFail(t *testing.T) {
	for _, s := range []string{
		"49:12",  // Invalid hours
		"25:12",  // Invalid hours
		"3:60",   // Invalid minutes
		"3:87",   // Invalid minutes
		"24:00>", // This would require shifting twice
		"24:01",  // The 24-hour special case can’t have minutes
		"13:00am",
		"13:00pm",
		"0:00am", // There is no `0` hour when using am/pm
		"0:00pm", // There is no `0` hour when using am/pm
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
	midnight := Ɀ_Time_(0, 0)
	midnight2 := Ɀ_Time_(0, 0)
	noon := Ɀ_Time_(12, 30)
	noon2 := Ɀ_Time_(12, 31)
	yesterday := Ɀ_TimeYesterday_(22, 43)
	tomorrow := Ɀ_TimeTomorrow_(9, 50)
	assert.True(t, midnight2.IsAfterOrEqual(midnight))
	assert.True(t, noon.IsAfterOrEqual(noon))
	assert.True(t, noon2.IsAfterOrEqual(noon))
	assert.True(t, noon.IsAfterOrEqual(yesterday))
	assert.True(t, tomorrow.IsAfterOrEqual(noon))
}

func TestAddDuration(t *testing.T) {
	for _, x := range []struct {
		initial   Time
		increment Duration
		expect    Time
	}{
		{Ɀ_Time_(11, 30), NewDuration(0, 00), Ɀ_Time_(11, 30)},
		{Ɀ_Time_(11, 30), NewDuration(0, 30), Ɀ_Time_(12, 00)},
		{Ɀ_Time_(18, 00), NewDuration(-6, 0), Ɀ_Time_(12, 00)},
		{Ɀ_Time_(3, 59), NewDuration(8, 1), Ɀ_Time_(12, 00)},
		{Ɀ_Time_(23, 59), NewDuration(0, 1), Ɀ_TimeTomorrow_(0, 00)},
		{Ɀ_TimeYesterday_(23, 45), NewDuration(12, 15), Ɀ_Time_(12, 00)},
		{Ɀ_TimeYesterday_(12, 12), NewDuration(1, 19), Ɀ_TimeYesterday_(13, 31)},
		{Ɀ_TimeYesterday_(0, 1), NewDuration(0, -1), Ɀ_TimeYesterday_(0, 0)},
		{Ɀ_TimeTomorrow_(4, 12), NewDuration(-16, -12), Ɀ_Time_(12, 00)},
		{Ɀ_TimeTomorrow_(18, 38), NewDuration(-1, -1), Ɀ_TimeTomorrow_(17, 37)},
		{Ɀ_TimeTomorrow_(23, 58), NewDuration(0, 1), Ɀ_TimeTomorrow_(23, 59)},
	} {
		result, err := x.initial.Plus(x.increment)
		require.Nil(t, err)
		assert.Equal(t, x.expect, result, x.initial)
	}
}

func TestAddDurationPreservesFormat(t *testing.T) {
	hour24, _ := Ɀ_Time_(11, 30).Plus(NewDuration(0, 1))
	assert.Equal(t, hour24.Format().Use24HourClock, true)

	hour12, _ := Ɀ_IsAmPm_(Ɀ_Time_(11, 30)).Plus(NewDuration(0, 1))
	assert.Equal(t, hour12.Format().Use24HourClock, false)
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
		result, err := x.initial.Plus(x.increment)
		require.Nil(t, result)
		assert.Error(t, err)
	}
}
