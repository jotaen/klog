package service

import (
	"github.com/jotaen/klog/klog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func r(f int) Rounding {
	result, err := NewRounding(f)
	if err != nil {
		panic(err)
	}
	return result
}

func TestValidRoundingValues(t *testing.T) {
	for _, m := range []int{
		5, 10, 15, 30, 60,
	} {
		fr, err := NewRounding(m)
		require.Nil(t, err)
		assert.Equal(t, m, fr.toInt())
	}
}

func TestInvalidRoundingValues(t *testing.T) {
	for _, m := range []int{
		-60, -30, -15, -10, -5, 0, 1, 2, 3, 4, 6, 7, 8, 9, 11, 12, 13, 14, 16, 17, 18, 19, 20, 25, 35, 45, 55, 120, 600,
	} {
		fr, err := NewRounding(m)
		require.Nil(t, fr)
		assert.Error(t, err)
	}
}

func TestParseRoundingValuesFromString(t *testing.T) {
	for _, x := range []struct {
		value    string
		expected int
	}{
		{"5", 5},
		{"5m", 5},
		{"10", 10},
		{"10m", 10},
		{"15m", 15},
		{"30m", 30},
		{"60m", 60},
		{"1h", 60},
	} {
		fr, err := NewRoundingFromString(x.value)
		require.Nil(t, err)
		assert.Equal(t, x.expected, fr.toInt())
	}
}

func TestInvalidRoundingValuesFromStringFail(t *testing.T) {
	for _, v := range []string{
		"0", "1", "11", "a", "", "5h",
	} {
		fr, err := NewRoundingFromString(v)
		require.Nil(t, fr)
		require.Error(t, err)
	}
}

func TestRound(t *testing.T) {
	for _, tm := range []struct {
		original klog.Time
		exp      klog.Time
	}{
		// Round to 5s
		{RoundToNearest(klog.Ɀ_Time_(8, 00), r(5)), klog.Ɀ_Time_(8, 00)},
		{RoundToNearest(klog.Ɀ_Time_(8, 01), r(5)), klog.Ɀ_Time_(8, 00)},
		{RoundToNearest(klog.Ɀ_Time_(8, 02), r(5)), klog.Ɀ_Time_(8, 00)},
		{RoundToNearest(klog.Ɀ_Time_(8, 03), r(5)), klog.Ɀ_Time_(8, 05)},
		{RoundToNearest(klog.Ɀ_Time_(8, 04), r(5)), klog.Ɀ_Time_(8, 05)},
		{RoundToNearest(klog.Ɀ_Time_(8, 05), r(5)), klog.Ɀ_Time_(8, 05)},
		{RoundToNearest(klog.Ɀ_Time_(8, 06), r(5)), klog.Ɀ_Time_(8, 05)},
		{RoundToNearest(klog.Ɀ_Time_(8, 07), r(5)), klog.Ɀ_Time_(8, 05)},
		{RoundToNearest(klog.Ɀ_Time_(8, 8), r(5)), klog.Ɀ_Time_(8, 10)},
		{RoundToNearest(klog.Ɀ_Time_(8, 57), r(5)), klog.Ɀ_Time_(8, 55)},
		{RoundToNearest(klog.Ɀ_Time_(8, 58), r(5)), klog.Ɀ_Time_(9, 00)},
		{RoundToNearest(klog.Ɀ_Time_(8, 59), r(5)), klog.Ɀ_Time_(9, 00)},

		// Round to 10s
		{RoundToNearest(klog.Ɀ_Time_(8, 00), r(10)), klog.Ɀ_Time_(8, 00)},
		{RoundToNearest(klog.Ɀ_Time_(8, 04), r(10)), klog.Ɀ_Time_(8, 00)},
		{RoundToNearest(klog.Ɀ_Time_(8, 05), r(10)), klog.Ɀ_Time_(8, 10)},
		{RoundToNearest(klog.Ɀ_Time_(8, 10), r(10)), klog.Ɀ_Time_(8, 10)},
		{RoundToNearest(klog.Ɀ_Time_(8, 14), r(10)), klog.Ɀ_Time_(8, 10)},
		{RoundToNearest(klog.Ɀ_Time_(8, 15), r(10)), klog.Ɀ_Time_(8, 20)},
		{RoundToNearest(klog.Ɀ_Time_(8, 55), r(10)), klog.Ɀ_Time_(9, 00)},

		// Round to 15s
		{RoundToNearest(klog.Ɀ_Time_(8, 00), r(15)), klog.Ɀ_Time_(8, 00)},
		{RoundToNearest(klog.Ɀ_Time_(8, 07), r(15)), klog.Ɀ_Time_(8, 00)},
		{RoundToNearest(klog.Ɀ_Time_(8, 8), r(15)), klog.Ɀ_Time_(8, 15)},
		{RoundToNearest(klog.Ɀ_Time_(8, 15), r(15)), klog.Ɀ_Time_(8, 15)},
		{RoundToNearest(klog.Ɀ_Time_(8, 22), r(15)), klog.Ɀ_Time_(8, 15)},
		{RoundToNearest(klog.Ɀ_Time_(8, 23), r(15)), klog.Ɀ_Time_(8, 30)},

		// Round to 30s
		{RoundToNearest(klog.Ɀ_Time_(8, 00), r(30)), klog.Ɀ_Time_(8, 00)},
		{RoundToNearest(klog.Ɀ_Time_(8, 14), r(30)), klog.Ɀ_Time_(8, 00)},
		{RoundToNearest(klog.Ɀ_Time_(8, 15), r(30)), klog.Ɀ_Time_(8, 30)},
		{RoundToNearest(klog.Ɀ_Time_(8, 44), r(30)), klog.Ɀ_Time_(8, 30)},
		{RoundToNearest(klog.Ɀ_Time_(8, 45), r(30)), klog.Ɀ_Time_(9, 00)},

		// Round to 60s
		{RoundToNearest(klog.Ɀ_Time_(8, 00), r(60)), klog.Ɀ_Time_(8, 00)},
		{RoundToNearest(klog.Ɀ_Time_(8, 29), r(60)), klog.Ɀ_Time_(8, 00)},
		{RoundToNearest(klog.Ɀ_Time_(8, 30), r(60)), klog.Ɀ_Time_(9, 00)},
		{RoundToNearest(klog.Ɀ_Time_(8, 59), r(60)), klog.Ɀ_Time_(9, 00)},

		// Round near day border
		{RoundToNearest(klog.Ɀ_Time_(0, 01), r(5)), klog.Ɀ_Time_(0, 00)},
		{RoundToNearest(klog.Ɀ_Time_(23, 59), r(15)), klog.Ɀ_TimeTomorrow_(0, 00)},
		{RoundToNearest(klog.Ɀ_TimeYesterday_(23, 59), r(10)), klog.Ɀ_Time_(0, 00)},
		// It can’t get higher than `23:59>`:
		{RoundToNearest(klog.Ɀ_TimeTomorrow_(23, 59), r(60)), klog.Ɀ_TimeTomorrow_(23, 59)},
	} {
		assert.Equal(t, tm.exp, tm.original)
	}
}
