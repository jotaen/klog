package period

import (
	. "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseValidYear(t *testing.T) {
	for _, x := range []struct {
		text   string
		expect Period
	}{
		{"0000", NewPeriod(Ɀ_Date_(0, 1, 1), Ɀ_Date_(0, 12, 31))},
		{"0475", NewPeriod(Ɀ_Date_(475, 1, 1), Ɀ_Date_(475, 12, 31))},
		{"2008", NewPeriod(Ɀ_Date_(2008, 1, 1), Ɀ_Date_(2008, 12, 31))},
		{"8641", NewPeriod(Ɀ_Date_(8641, 1, 1), Ɀ_Date_(8641, 12, 31))},
		{"9999", NewPeriod(Ɀ_Date_(9999, 1, 1), Ɀ_Date_(9999, 12, 31))},
	} {
		year, err := NewYearFromString(x.text)
		require.Nil(t, err)
		assert.True(t, x.expect.Since().IsEqualTo(year.Period().Since()))
		assert.True(t, x.expect.Until().IsEqualTo(year.Period().Until()))
	}
}

func TestRejectsInvalidYear(t *testing.T) {
	for _, x := range []string{
		"-5",
		"10000",
		"9823746",
	} {
		_, err := NewYearFromString(x)
		require.Error(t, err)
	}
}

func TestRejectsMalformedYear(t *testing.T) {
	for _, x := range []string{
		"",
		"asdf",
		"2oo1",
	} {
		_, err := NewYearFromString(x)
		require.Error(t, err)
	}
}

func TestYearPeriod(t *testing.T) {
	for _, x := range []struct {
		initial  Year
		expected Period
	}{
		{NewYearFromDate(Ɀ_Date_(1987, 5, 19)), NewPeriod(Ɀ_Date_(1987, 1, 1), Ɀ_Date_(1987, 12, 31))},
		{NewYearFromDate(Ɀ_Date_(2000, 3, 31)), NewPeriod(Ɀ_Date_(2000, 1, 1), Ɀ_Date_(2000, 12, 31))},
		{NewYearFromDate(Ɀ_Date_(2555, 12, 31)), NewPeriod(Ɀ_Date_(2555, 1, 1), Ɀ_Date_(2555, 12, 31))},
	} {
		p := x.initial.Period()
		assert.Equal(t, x.expected, p)
	}
}

func TestYearPreviousYear(t *testing.T) {
	for _, x := range []struct {
		initial  Year
		expected Period
	}{
		{NewYearFromDate(Ɀ_Date_(1987, 5, 19)), NewPeriod(Ɀ_Date_(1986, 1, 1), Ɀ_Date_(1986, 12, 31))},
		{NewYearFromDate(Ɀ_Date_(2000, 3, 31)), NewPeriod(Ɀ_Date_(1999, 1, 1), Ɀ_Date_(1999, 12, 31))},
		{NewYearFromDate(Ɀ_Date_(2555, 12, 31)), NewPeriod(Ɀ_Date_(2554, 1, 1), Ɀ_Date_(2554, 12, 31))},
	} {
		previous := x.initial.Previous().Period()
		assert.Equal(t, x.expected, previous)
	}
}
