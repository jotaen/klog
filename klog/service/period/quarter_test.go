package period

import (
	"github.com/jotaen/klog/klog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQuarterPeriod(t *testing.T) {
	for _, x := range []struct {
		actual   Period
		expected Period
	}{
		// Q1
		{NewQuarterFromDate(klog.Ɀ_Date_(1999, 1, 19)).Period(), NewPeriod(klog.Ɀ_Date_(1999, 1, 1), klog.Ɀ_Date_(1999, 3, 31))},

		// Q2
		{NewQuarterFromDate(klog.Ɀ_Date_(2005, 5, 19)).Period(), NewPeriod(klog.Ɀ_Date_(2005, 4, 1), klog.Ɀ_Date_(2005, 6, 30))},

		// Q3
		{NewQuarterFromDate(klog.Ɀ_Date_(1589, 8, 3)).Period(), NewPeriod(klog.Ɀ_Date_(1589, 7, 1), klog.Ɀ_Date_(1589, 9, 30))},

		// Q4
		{NewQuarterFromDate(klog.Ɀ_Date_(2134, 12, 30)).Period(), NewPeriod(klog.Ɀ_Date_(2134, 10, 1), klog.Ɀ_Date_(2134, 12, 31))},

		// Since is same as original date
		{NewQuarterFromDate(klog.Ɀ_Date_(1998, 4, 1)).Period(), NewPeriod(klog.Ɀ_Date_(1998, 4, 1), klog.Ɀ_Date_(1998, 6, 30))},

		// Until is same as original date
		{NewQuarterFromDate(klog.Ɀ_Date_(1998, 9, 30)).Period(), NewPeriod(klog.Ɀ_Date_(1998, 7, 1), klog.Ɀ_Date_(1998, 9, 30))},
	} {
		assert.Equal(t, x.expected, x.actual)
	}
}

func TestParseValidQuarter(t *testing.T) {
	for _, x := range []struct {
		text   string
		expect Period
	}{
		{"0000-Q1", NewPeriod(klog.Ɀ_Date_(0, 1, 1), klog.Ɀ_Date_(0, 3, 31))},
		{"0475-Q2", NewPeriod(klog.Ɀ_Date_(475, 4, 1), klog.Ɀ_Date_(475, 6, 30))},
		{"2008-Q3", NewPeriod(klog.Ɀ_Date_(2008, 7, 1), klog.Ɀ_Date_(2008, 9, 30))},
		{"8641-Q4", NewPeriod(klog.Ɀ_Date_(8641, 10, 1), klog.Ɀ_Date_(8641, 12, 31))},
	} {
		quarter, err := NewQuarterFromString(x.text)
		require.Nil(t, err)
		assert.True(t, x.expect.Since().IsEqualTo(quarter.Period().Since()))
		assert.True(t, x.expect.Until().IsEqualTo(quarter.Period().Until()))
	}
}

func TestParseRejectsInvalidQuarter(t *testing.T) {
	for _, x := range []string{
		"2000-Q5",
		"2000-Q0",
		"2000-Q-1",
		"2000-Q",
		"2000-q2",
		"2000-asdf",
		"2000-Q01",
		"2000-",
		"273888-Q2",
		"Q3",
	} {
		_, err := NewQuarterFromString(x)
		require.Error(t, err)
	}
}

func TestQuarterPreviousQuarter(t *testing.T) {
	for _, x := range []struct {
		initial  Quarter
		expected Period
	}{
		// In same year
		{NewQuarterFromDate(klog.Ɀ_Date_(1987, 5, 19)), NewPeriod(klog.Ɀ_Date_(1987, 1, 1), klog.Ɀ_Date_(1987, 3, 31))},
		{NewQuarterFromDate(klog.Ɀ_Date_(1987, 4, 19)), NewPeriod(klog.Ɀ_Date_(1987, 1, 1), klog.Ɀ_Date_(1987, 3, 31))},
		{NewQuarterFromDate(klog.Ɀ_Date_(1444, 8, 13)), NewPeriod(klog.Ɀ_Date_(1444, 4, 1), klog.Ɀ_Date_(1444, 6, 30))},
		{NewQuarterFromDate(klog.Ɀ_Date_(2009, 12, 31)), NewPeriod(klog.Ɀ_Date_(2009, 7, 1), klog.Ɀ_Date_(2009, 9, 30))},
		{NewQuarterFromDate(klog.Ɀ_Date_(2009, 10, 1)), NewPeriod(klog.Ɀ_Date_(2009, 7, 1), klog.Ɀ_Date_(2009, 9, 30))},

		// In last year
		{NewQuarterFromDate(klog.Ɀ_Date_(1987, 1, 1)), NewPeriod(klog.Ɀ_Date_(1986, 10, 1), klog.Ɀ_Date_(1986, 12, 31))},
		{NewQuarterFromDate(klog.Ɀ_Date_(2400, 2, 27)), NewPeriod(klog.Ɀ_Date_(2399, 10, 1), klog.Ɀ_Date_(2399, 12, 31))},
	} {
		previous := x.initial.Previous().Period()
		assert.Equal(t, x.expected, previous)
	}
}
