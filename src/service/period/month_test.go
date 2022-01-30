package period

import (
	. "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseValidMonth(t *testing.T) {
	for _, x := range []struct {
		text   string
		expect Period
	}{
		{"0000-01", NewPeriod(Ɀ_Date_(0, 1, 1), Ɀ_Date_(0, 01, 31))},
		{"0000-12", NewPeriod(Ɀ_Date_(0, 12, 1), Ɀ_Date_(0, 12, 31))},
		{"0475-05", NewPeriod(Ɀ_Date_(475, 5, 1), Ɀ_Date_(475, 5, 31))},
		{"2008-11", NewPeriod(Ɀ_Date_(2008, 11, 1), Ɀ_Date_(2008, 11, 30))},
		{"8641-04", NewPeriod(Ɀ_Date_(8641, 4, 1), Ɀ_Date_(8641, 4, 30))},
		{"9999-12", NewPeriod(Ɀ_Date_(9999, 12, 1), Ɀ_Date_(9999, 12, 31))},
	} {
		month, err := NewMonthFromString(x.text)
		require.Nil(t, err)
		assert.True(t, x.expect.Since().IsEqualTo(month.Period().Since()))
		assert.True(t, x.expect.Until().IsEqualTo(month.Period().Until()))
	}
}

func TestMonthEnds(t *testing.T) {
	for _, x := range []struct {
		text    string
		month   int
		lastDay int
	}{
		{"2018-01", 1, 31},
		{"2018-02", 2, 28},
		{"2018-03", 3, 31},
		{"2018-04", 4, 30},
		{"2018-05", 5, 31},
		{"2018-06", 6, 30},
		{"2018-07", 7, 31},
		{"2018-08", 8, 31},
		{"2018-09", 9, 30},
		{"2018-10", 10, 31},
		{"2018-11", 11, 30},
		{"2018-12", 12, 31},
	} {
		m, err := NewMonthFromString(x.text)
		require.Nil(t, err)
		p := m.Period()
		assert.Equal(t, p.Since(), Ɀ_Date_(2018, x.month, 1))
		assert.Equal(t, p.Until(), Ɀ_Date_(2018, x.month, x.lastDay))
	}
}

func TestParseMonthInLeapYear(t *testing.T) {
	m, _ := NewMonthFromString("2016-02")
	assert.Equal(t, m.Period().Until(), Ɀ_Date_(2016, 2, 29))
}

func TestRejectsInvalidMonth(t *testing.T) {
	for _, x := range []string{
		"4000-00",
		"4000-13",
		"1833716-01",
		"2008-1",
	} {
		_, err := NewMonthFromString(x)
		require.Error(t, err)
	}
}

func TestRejectsMalformedMonth(t *testing.T) {
	for _, x := range []string{
		"",
		"asdf",
		"2005",
		"2005_12",
		"2005--12",
	} {
		_, err := NewMonthFromString(x)
		require.Error(t, err)
	}
}

func TestMonthPeriod(t *testing.T) {
	for _, x := range []struct {
		actual   Period
		expected Period
	}{
		// Range in same year
		{NewMonthFromDate(Ɀ_Date_(1987, 5, 19)).Period(), NewPeriod(Ɀ_Date_(1987, 5, 1), Ɀ_Date_(1987, 5, 31))},
		{NewMonthFromDate(Ɀ_Date_(2004, 11, 16)).Period(), NewPeriod(Ɀ_Date_(2004, 11, 1), Ɀ_Date_(2004, 11, 30))},

		// Since is same as original date
		{NewMonthFromDate(Ɀ_Date_(1998, 10, 1)).Period(), NewPeriod(Ɀ_Date_(1998, 10, 1), Ɀ_Date_(1998, 10, 31))},

		// Until is same as original date
		{NewMonthFromDate(Ɀ_Date_(1998, 2, 28)).Period(), NewPeriod(Ɀ_Date_(1998, 2, 1), Ɀ_Date_(1998, 2, 28))},

		// Leap year
		{NewMonthFromDate(Ɀ_Date_(2000, 2, 4)).Period(), NewPeriod(Ɀ_Date_(2000, 2, 1), Ɀ_Date_(2000, 2, 29))},
	} {
		assert.Equal(t, x.expected, x.actual)
	}
}

func TestMonthPreviousMonth(t *testing.T) {
	for _, x := range []struct {
		initial  Month
		expected Period
	}{
		// In same year
		{NewMonthFromDate(Ɀ_Date_(1987, 5, 19)), NewPeriod(Ɀ_Date_(1987, 4, 1), Ɀ_Date_(1987, 4, 30))},
		{NewMonthFromDate(Ɀ_Date_(1987, 3, 31)), NewPeriod(Ɀ_Date_(1987, 2, 1), Ɀ_Date_(1987, 2, 28))},
		{NewMonthFromDate(Ɀ_Date_(1987, 3, 1)), NewPeriod(Ɀ_Date_(1987, 2, 1), Ɀ_Date_(1987, 2, 28))},

		// In last year
		{NewMonthFromDate(Ɀ_Date_(1987, 1, 19)), NewPeriod(Ɀ_Date_(1986, 12, 1), Ɀ_Date_(1986, 12, 31))},
	} {
		previous := x.initial.Previous().Period()
		assert.Equal(t, x.expected, previous)
	}
}
