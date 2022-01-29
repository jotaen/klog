package period

import (
	. "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseValidYear(t *testing.T) {
	for _, x := range []struct {
		text         string
		expectYear   Year
		expectPeriod Period
	}{
		{"0000", Year{0}, Period{Ɀ_Date_(0, 1, 1), Ɀ_Date_(0, 12, 31)}},
		{"0475", Year{475}, Period{Ɀ_Date_(475, 1, 1), Ɀ_Date_(475, 12, 31)}},
		{"2008", Year{2008}, Period{Ɀ_Date_(2008, 1, 1), Ɀ_Date_(2008, 12, 31)}},
		{"8641", Year{8641}, Period{Ɀ_Date_(8641, 1, 1), Ɀ_Date_(8641, 12, 31)}},
		{"9999", Year{9999}, Period{Ɀ_Date_(9999, 1, 1), Ɀ_Date_(9999, 12, 31)}},
	} {
		year, err := NewYearFromString(x.text)
		require.Nil(t, err)
		assert.Equal(t, x.expectYear, year)
		assert.True(t, x.expectPeriod.Since.IsEqualTo(year.Period().Since))
		assert.True(t, x.expectPeriod.Until.IsEqualTo(year.Period().Until))
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
