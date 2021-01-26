package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "klog"
	"testing"
	gotime "time"
)

func TestTotalSumUpZeroIfNoTimesSpecified(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	assert.Equal(t, NewDuration(0, 0), Total(r))
}

func TestTotalSumsUpTimesAndRangesButNotOpenRanges(t *testing.T) {
	r1 := NewRecord(Ɀ_Date_(2020, 1, 1))
	r1.AddDuration(NewDuration(3, 0), "")
	r1.AddDuration(NewDuration(1, 33), "")
	r1.AddRange(Ɀ_Range_(Ɀ_TimeYesterday_(8, 0), Ɀ_TimeTomorrow_(12, 0)), "")
	r1.AddRange(Ɀ_Range_(Ɀ_Time_(13, 49), Ɀ_Time_(17, 12)), "")
	_ = r1.StartOpenRange(Ɀ_Time_(1, 2), "")
	r2 := NewRecord(Ɀ_Date_(2020, 1, 2))
	r2.AddDuration(NewDuration(7, 55), "")
	assert.Equal(t, NewDuration(3+1+(16+24+12)+3+7, 33+11+12+55), Total(r1, r2))
}

func TestSumUpHypotheticalTotalAtGivenTime(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.AddDuration(NewDuration(2, 14), "")
	r.AddRange(Ɀ_Range_(Ɀ_TimeYesterday_(23, 0), Ɀ_Time_(4, 0)), "")
	_ = r.StartOpenRange(Ɀ_Time_(5, 7), "")

	time1, _ := gotime.Parse("2006-01-02T15:04:05-0700", "2020-01-01T05:06:59-0000")
	ht1, isOngoing1 := HypotheticalTotal(time1, r)
	assert.False(t, isOngoing1)
	assert.Equal(t, NewDuration(2+(1+4), 14), ht1)

	time2, _ := gotime.Parse("2006-01-02T15:04:05-0700", "2020-01-01T10:48:13-0000")
	ht2, isOngoing2 := HypotheticalTotal(time2, r)
	assert.True(t, isOngoing2)
	assert.Equal(t, NewDuration(2+(1+4)+4, 14+53+48), ht2)

	time3, _ := gotime.Parse("2006-01-02T15:04:05-0700", "2020-01-02T03:01:29-0000")
	ht3, isOngoing3 := HypotheticalTotal(time3, r)
	assert.True(t, isOngoing3)
	assert.Equal(t, NewDuration(2+(1+4)+18+3, 14+53+1), ht3)
}

func TestHashTagAllEntriesAreReturnedIfMatchIsInSummary(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	_ = r.SetSummary("This and #that, and other stuff as well")
	r.AddDuration(NewDuration(2, 0), "Foo")
	r.AddRange(Ɀ_Range_(Ɀ_Time_(13, 49), Ɀ_Time_(17, 12)), "Bar")
	es, hasMatched := FindEntriesWithHashtags(NewTagSet("that"), r)
	require.Len(t, es, 2)
	assert.True(t, hasMatched)
	assert.Equal(t, Summary("Foo"), es[0].Summary())
	assert.Equal(t, Summary("Bar"), es[1].Summary())
}

func TestHashTagReturnsEntriesThatMatch(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	_ = r.SetSummary("This and #that, and other stuff as well")
	r.AddDuration(NewDuration(2, 0), "Foo #fizzbuzz")
	r.AddRange(Ɀ_Range_(Ɀ_Time_(13, 49), Ɀ_Time_(17, 12)), "Bar #barbaz")
	es, hasMatched := FindEntriesWithHashtags(NewTagSet("fizzbuzz"), r)
	require.Len(t, es, 1)
	assert.True(t, hasMatched)
	assert.Equal(t, Summary("Foo #fizzbuzz"), es[0].Summary())
}

func sampleRecords() []Record {
	return []Record{
		func() Record {
			r := NewRecord(Ɀ_Date_(1999, 12, 30))
			r.AddDuration(NewDuration(4, 0), "")
			return r
		}(), func() Record {
			r := NewRecord(Ɀ_Date_(1999, 12, 31))
			r.AddDuration(NewDuration(5, 0), "#newYearsEve")
			return r
		}(), func() Record {
			r := NewRecord(Ɀ_Date_(2000, 1, 1))
			r.AddDuration(NewDuration(6, 0), "#millennium")
			r.AddDuration(NewDuration(6, 30), "")
			return r
		}(), func() Record {
			r := NewRecord(Ɀ_Date_(2000, 1, 2))
			r.AddDuration(NewDuration(7, 0), "")
			return r
		}(), func() Record {
			r := NewRecord(Ɀ_Date_(2000, 1, 3))
			r.AddDuration(NewDuration(8, 0), "")
			return r
		}(),
	}
}

func TestFindFilterWithNoClauses(t *testing.T) {
	rs, es := FindFilter(sampleRecords(), Filter{})
	require.Len(t, rs, 5)
	assert.Equal(t, NewDuration(4+5+6+6+7+8, 30), TotalEntries(es))
}

func TestFindFilterWithAfter(t *testing.T) {
	rs, _ := FindFilter(sampleRecords(), Filter{AfterEq: Ɀ_Date_(2000, 1, 1)})
	require.Len(t, rs, 3)
	assert.Equal(t, 1, rs[0].Date().Day())
	assert.Equal(t, 2, rs[1].Date().Day())
	assert.Equal(t, 3, rs[2].Date().Day())
}

func TestFindFilterWithBefore(t *testing.T) {
	rs, _ := FindFilter(sampleRecords(), Filter{BeforeEq: Ɀ_Date_(2000, 1, 1)})
	require.Len(t, rs, 3)
	assert.Equal(t, 30, rs[0].Date().Day())
	assert.Equal(t, 31, rs[1].Date().Day())
	assert.Equal(t, 1, rs[2].Date().Day())
}

func TestFindFilterWithHash(t *testing.T) {
	rs, es := FindFilter(sampleRecords(), Filter{Tags: []string{"newYearsEve", "millennium"}})
	require.Len(t, rs, 2)
	require.Len(t, es, 2)
	assert.Equal(t, 31, rs[0].Date().Day())
	assert.Equal(t, 1, rs[1].Date().Day())
	assert.Equal(t, NewDuration(5+6, 0), TotalEntries(es))
}
