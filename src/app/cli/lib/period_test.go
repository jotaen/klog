package lib

import (
	. "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseValidPeriodWithYear(t *testing.T) {
	p, err := NewPeriodFromString("2018")
	require.Nil(t, err)
	assert.Equal(t, p.Since, Ɀ_Date_(2018, 1, 1))
	assert.Equal(t, p.Until, Ɀ_Date_(2018, 12, 31))
}

func TestParseValidPeriodWithYearAndMonth(t *testing.T) {
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
		p, err := NewPeriodFromString(x.text)
		require.Nil(t, err)
		assert.Equal(t, p.Since, Ɀ_Date_(2018, x.month, 1))
		assert.Equal(t, p.Until, Ɀ_Date_(2018, x.month, x.lastDay))
	}
}

func TestParsePeriodWithLeapYear(t *testing.T) {
	p, _ := NewPeriodFromString("2016-02")
	assert.Equal(t, p.Until, Ɀ_Date_(2016, 2, 29))
}

func TestFailParsingWithMalformedInput(t *testing.T) {
	for _, x := range []string{
		"",
		"asdf",
		"2018-",
		"2018-3",
		"2018-a",
		"20-03",
		"-03",
		"03",
	} {
		_, err := NewPeriodFromString(x)
		require.Error(t, err)
	}
}
