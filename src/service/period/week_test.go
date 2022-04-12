package period

import (
	. "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWeekPeriod(t *testing.T) {
	for _, x := range []struct {
		actual   Period
		expected Period
	}{
		// Range in same month
		{NewWeekFromDate(Ɀ_Date_(1987, 5, 19)).Period(), NewPeriod(Ɀ_Date_(1987, 5, 18), Ɀ_Date_(1987, 5, 24))},
		{NewWeekFromDate(Ɀ_Date_(2004, 12, 16)).Period(), NewPeriod(Ɀ_Date_(2004, 12, 13), Ɀ_Date_(2004, 12, 19))},

		// Range across months
		{NewWeekFromDate(Ɀ_Date_(1983, 6, 1)).Period(), NewPeriod(Ɀ_Date_(1983, 5, 30), Ɀ_Date_(1983, 6, 5))},
		{NewWeekFromDate(Ɀ_Date_(1998, 10, 27)).Period(), NewPeriod(Ɀ_Date_(1998, 10, 26), Ɀ_Date_(1998, 11, 1))},

		// Range across years
		{NewWeekFromDate(Ɀ_Date_(2009, 1, 2)).Period(), NewPeriod(Ɀ_Date_(2008, 12, 29), Ɀ_Date_(2009, 1, 4))},
		{NewWeekFromDate(Ɀ_Date_(2009, 12, 30)).Period(), NewPeriod(Ɀ_Date_(2009, 12, 28), Ɀ_Date_(2010, 1, 3))},

		// Since is same as original date
		{NewWeekFromDate(Ɀ_Date_(1998, 10, 26)).Period(), NewPeriod(Ɀ_Date_(1998, 10, 26), Ɀ_Date_(1998, 11, 1))},

		// Until is same as original date
		{NewWeekFromDate(Ɀ_Date_(1998, 11, 1)).Period(), NewPeriod(Ɀ_Date_(1998, 10, 26), Ɀ_Date_(1998, 11, 1))},
	} {
		assert.Equal(t, x.expected, x.actual)
	}
}

func TestParseValidWeek(t *testing.T) {
	for _, x := range []struct {
		text   string
		expect Period
	}{
		{"2022-W01", NewPeriod(Ɀ_Date_(2022, 1, 3), Ɀ_Date_(2022, 1, 9))},
		{"2022-W1", NewPeriod(Ɀ_Date_(2022, 1, 3), Ɀ_Date_(2022, 1, 9))},
		{"2017-W26", NewPeriod(Ɀ_Date_(2017, 6, 26), Ɀ_Date_(2017, 7, 2))},
		{"2017-W27", NewPeriod(Ɀ_Date_(2017, 7, 3), Ɀ_Date_(2017, 7, 9))},
		{"2012-W09", NewPeriod(Ɀ_Date_(2012, 2, 27), Ɀ_Date_(2012, 3, 4))},
		{"2022-W02", NewPeriod(Ɀ_Date_(2022, 1, 10), Ɀ_Date_(2022, 1, 16))},
		{"2022-W52", NewPeriod(Ɀ_Date_(2022, 12, 26), Ɀ_Date_(2023, 1, 1))},
		{"2025-W01", NewPeriod(Ɀ_Date_(2024, 12, 30), Ɀ_Date_(2025, 1, 5))},
	} {
		week, err := NewWeekFromString(x.text)
		require.Nil(t, err)
		assert.True(t, x.expect.Since().IsEqualTo(week.Period().Since()), x.text)
		assert.True(t, x.expect.Until().IsEqualTo(week.Period().Until()))
	}
}

func TestParseRejectsInvalidWeekString(t *testing.T) {
	for _, x := range []string{
		"2000-W00",
		"2000-W-1",
		"2000-W001",
		"2000-W54",
		"2000-W",
		"2000-w14",
		"2000-w14",
		"2000-asdf",
		"12873612-W02",
	} {
		_, err := NewWeekFromString(x)
		require.Error(t, err)
	}
}

func TestWeekPreviousWeek(t *testing.T) {
	for _, x := range []struct {
		initial  Week
		expected Period
	}{
		// Same month
		{NewWeekFromDate(Ɀ_Date_(1987, 5, 19)), NewPeriod(Ɀ_Date_(1987, 5, 11), Ɀ_Date_(1987, 5, 17))},

		// `Since` in other month
		{NewWeekFromDate(Ɀ_Date_(2014, 8, 6)), NewPeriod(Ɀ_Date_(2014, 7, 28), Ɀ_Date_(2014, 8, 3))},

		// `Since`&`Until` in other month
		{NewWeekFromDate(Ɀ_Date_(2014, 8, 2)), NewPeriod(Ɀ_Date_(2014, 7, 21), Ɀ_Date_(2014, 7, 27))},

		// `Since` in other year
		{NewWeekFromDate(Ɀ_Date_(2014, 1, 9)), NewPeriod(Ɀ_Date_(2013, 12, 30), Ɀ_Date_(2014, 1, 5))},

		// `Since`&`Until` in other year
		{NewWeekFromDate(Ɀ_Date_(2029, 1, 2)), NewPeriod(Ɀ_Date_(2028, 12, 25), Ɀ_Date_(2028, 12, 31))},
	} {
		previous := x.initial.Previous().Period()
		assert.Equal(t, x.expected, previous)
	}
}
