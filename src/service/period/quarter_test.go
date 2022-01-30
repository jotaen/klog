package period

import (
	. "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuarterPeriod(t *testing.T) {
	for _, x := range []struct {
		actual   Period
		expected Period
	}{
		// Q1
		{NewQuarterFromDate(Ɀ_Date_(1999, 1, 19)).Period(), NewPeriod(Ɀ_Date_(1999, 1, 1), Ɀ_Date_(1999, 3, 31))},

		// Q2
		{NewQuarterFromDate(Ɀ_Date_(2005, 5, 19)).Period(), NewPeriod(Ɀ_Date_(2005, 4, 1), Ɀ_Date_(2005, 6, 30))},

		// Q3
		{NewQuarterFromDate(Ɀ_Date_(1589, 8, 3)).Period(), NewPeriod(Ɀ_Date_(1589, 7, 1), Ɀ_Date_(1589, 9, 30))},

		// Q4
		{NewQuarterFromDate(Ɀ_Date_(2134, 12, 30)).Period(), NewPeriod(Ɀ_Date_(2134, 10, 1), Ɀ_Date_(2134, 12, 31))},

		// Since is same as original date
		{NewQuarterFromDate(Ɀ_Date_(1998, 4, 1)).Period(), NewPeriod(Ɀ_Date_(1998, 4, 1), Ɀ_Date_(1998, 6, 30))},

		// Until is same as original date
		{NewQuarterFromDate(Ɀ_Date_(1998, 9, 30)).Period(), NewPeriod(Ɀ_Date_(1998, 7, 1), Ɀ_Date_(1998, 9, 30))},
	} {
		assert.Equal(t, x.expected, x.actual)
	}
}

func TestQuarterPreviousQuarter(t *testing.T) {
	for _, x := range []struct {
		initial  Quarter
		expected Period
	}{
		// In same year
		{NewQuarterFromDate(Ɀ_Date_(1987, 5, 19)), NewPeriod(Ɀ_Date_(1987, 1, 1), Ɀ_Date_(1987, 3, 31))},
		{NewQuarterFromDate(Ɀ_Date_(1987, 4, 19)), NewPeriod(Ɀ_Date_(1987, 1, 1), Ɀ_Date_(1987, 3, 31))},
		{NewQuarterFromDate(Ɀ_Date_(1444, 8, 13)), NewPeriod(Ɀ_Date_(1444, 4, 1), Ɀ_Date_(1444, 6, 30))},
		{NewQuarterFromDate(Ɀ_Date_(2009, 12, 31)), NewPeriod(Ɀ_Date_(2009, 7, 1), Ɀ_Date_(2009, 9, 30))},
		{NewQuarterFromDate(Ɀ_Date_(2009, 10, 1)), NewPeriod(Ɀ_Date_(2009, 7, 1), Ɀ_Date_(2009, 9, 30))},

		// In last year
		{NewQuarterFromDate(Ɀ_Date_(1987, 1, 1)), NewPeriod(Ɀ_Date_(1986, 10, 1), Ɀ_Date_(1986, 12, 31))},
		{NewQuarterFromDate(Ɀ_Date_(2400, 2, 27)), NewPeriod(Ɀ_Date_(2399, 10, 1), Ɀ_Date_(2399, 12, 31))},
	} {
		previous := x.initial.Previous().Period()
		assert.Equal(t, x.expected, previous)
	}
}
