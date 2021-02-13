package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "klog"
	"testing"
	gotime "time"
)

func TestNoWarningWhenAllIsOkay(t *testing.T) {
	ws := SanityCheck(
		gotime.Now(),
		sampleRecordsForChecking(gotime.Now().Add(gotime.Duration(1*24*60*60*1000000000))),
	)
	require.Nil(t, ws)
}

func TestWarnWhenUnclosedOpenRangeInThePast(t *testing.T) {
	ws := SanityCheck(
		gotime.Now(),
		sampleRecordsForChecking(gotime.Now()),
	)
	require.NotNil(t, ws)
	require.Len(t, ws, 1)
	assert.True(t, ws[0].Date.IsEqualTo(Ɀ_Date_(1999, 12, 30)))
}

func sampleRecordsForChecking(reference gotime.Time) []Record {
	today := NewDateFromTime(reference)
	now := NewTimeFromTime(gotime.Now())
	return []Record{
		func() Record {
			r := NewRecord(today.PlusDays(1))
			r.StartOpenRange(now, "")
			return r
		}(), func() Record {
			r := NewRecord(today)
			r.StartOpenRange(now, "")
			return r
		}(), func() Record {
			r := NewRecord(today.PlusDays(-1))
			r.StartOpenRange(now, "")
			return r
		}(), func() Record {
			r := NewRecord(today.PlusDays(-2))
			r.StartOpenRange(now, "")
			return r
		}(),
	}
}
