package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"klog"
	"testing"
)

func TestSumUpTimes(t *testing.T) {
	r := src.NewRecord(src.Ɀ_Date_(2020, 1, 1))
	r.AddDuration(src.NewDuration(1, 0), "")
	r.AddDuration(src.NewDuration(2, 0), "")
	assert.Equal(t, src.NewDuration(3, 0), Total(r))
}

func TestSumUpZeroIfNoTimesAvailable(t *testing.T) {
	r := src.NewRecord(src.Ɀ_Date_(2020, 1, 1))
	assert.Equal(t, src.NewDuration(0, 0), Total(r))
}

func TestSumUpRanges(t *testing.T) {
	range1 := src.Ɀ_Range_(src.Ɀ_Time_(9, 7), src.Ɀ_Time_(12, 59))
	range2 := src.Ɀ_Range_(src.Ɀ_Time_(13, 49), src.Ɀ_Time_(17, 12))
	r := src.NewRecord(src.Ɀ_Date_(2020, 1, 1))
	r.AddRange(range1, "")
	r.AddRange(range2, "")
	assert.Equal(t, src.NewDuration(7, 15), Total(r))
}

func TestSumUpTimesAndRanges(t *testing.T) {
	range1 := src.Ɀ_Range_(src.Ɀ_Time_(8, 0), src.Ɀ_Time_(12, 0))
	r := src.NewRecord(src.Ɀ_Date_(2020, 1, 1))
	r.AddDuration(src.NewDuration(1, 33), "")
	r.AddRange(range1, "")
	assert.Equal(t, src.NewDuration(5, 33), Total(r))
}

func TestHashTagAllEntriesAreReturnedIfMatchIsInSummary(t *testing.T) {
	r := src.NewRecord(src.Ɀ_Date_(2020, 1, 1))
	_ = r.SetSummary("This and #that, and other stuff as well")
	r.AddDuration(src.NewDuration(2, 0), "Foo")
	r.AddRange(src.Ɀ_Range_(src.Ɀ_Time_(13, 49), src.Ɀ_Time_(17, 12)), "Bar")
	es, hasMatched := FindEntriesWithHashtags(src.NewTagSet("that"), r)
	require.Len(t, es, 2)
	assert.True(t, hasMatched)
	assert.Equal(t, src.Summary("Foo"), es[0].Summary())
	assert.Equal(t, src.Summary("Bar"), es[1].Summary())
}

func TestHashTagReturnsEntriesThatMatch(t *testing.T) {
	r := src.NewRecord(src.Ɀ_Date_(2020, 1, 1))
	_ = r.SetSummary("This and #that, and other stuff as well")
	r.AddDuration(src.NewDuration(2, 0), "Foo #fizzbuzz")
	r.AddRange(src.Ɀ_Range_(src.Ɀ_Time_(13, 49), src.Ɀ_Time_(17, 12)), "Bar #barbaz")
	es, hasMatched := FindEntriesWithHashtags(src.NewTagSet("fizzbuzz"), r)
	require.Len(t, es, 1)
	assert.True(t, hasMatched)
	assert.Equal(t, src.Summary("Foo #fizzbuzz"), es[0].Summary())
}

func TestFindFilterWithNoClauses(t *testing.T) {
	rs, es := FindFilter([]src.Record{
		func() src.Record {
			r := src.NewRecord(src.Ɀ_Date_(2000, 1, 15))
			r.AddDuration(src.NewDuration(1, 0), "")
			return r
		}(),
		func() src.Record {
			r := src.NewRecord(src.Ɀ_Date_(2000, 1, 16))
			r.AddDuration(src.NewDuration(1, 0), "")
			return r
		}(),
		func() src.Record {
			r := src.NewRecord(src.Ɀ_Date_(2000, 1, 16))
			r.AddDuration(src.NewDuration(1, 0), "")
			return r
		}(),
	}, Filter{})

	require.Len(t, rs, 3)
	assert.Equal(t, src.NewDuration(3, 0), TotalEntries(es))
}

func TestFindFilterWithAfter(t *testing.T) {
	rs, _ := FindFilter([]src.Record{
		src.NewRecord(src.Ɀ_Date_(2000, 1, 15)),
		src.NewRecord(src.Ɀ_Date_(2000, 1, 16)),
		src.NewRecord(src.Ɀ_Date_(2000, 1, 17)),
	}, Filter{AfterEq: src.Ɀ_Date_(2000, 1, 16)})

	require.Len(t, rs, 2)
	assert.Equal(t, 16, rs[0].Date().Day())
	assert.Equal(t, 17, rs[1].Date().Day())
}

func TestFindFilterWithBefore(t *testing.T) {
	rs, _ := FindFilter([]src.Record{
		src.NewRecord(src.Ɀ_Date_(2000, 1, 15)),
		src.NewRecord(src.Ɀ_Date_(2000, 1, 16)),
		src.NewRecord(src.Ɀ_Date_(2000, 1, 17)),
	}, Filter{BeforeEq: src.Ɀ_Date_(2000, 1, 16)})

	require.Len(t, rs, 2)
	assert.Equal(t, 15, rs[0].Date().Day())
	assert.Equal(t, 16, rs[1].Date().Day())
}

func TestFindFilterWithHash(t *testing.T) {
	rs, es := FindFilter([]src.Record{
		func() src.Record {
			r := src.NewRecord(src.Ɀ_Date_(2000, 1, 15))
			_ = r.SetSummary("Contains #foo")
			r.AddDuration(src.NewDuration(3, 0), "")
			return r
		}(),
		func() src.Record {
			r := src.NewRecord(src.Ɀ_Date_(2000, 1, 16))
			r.AddDuration(src.NewDuration(1, 0), "Contains #foo too")
			return r
		}(),
		func() src.Record {
			r := src.NewRecord(src.Ɀ_Date_(2000, 1, 17))
			return r
		}(),
	}, Filter{Tags: []string{"foo"}})

	require.Len(t, rs, 2)
	require.Len(t, es, 2)
	assert.Equal(t, 15, rs[0].Date().Day())
	assert.Equal(t, 16, rs[1].Date().Day())
	assert.Equal(t, src.NewDuration(4, 0), TotalEntries(es))
}
