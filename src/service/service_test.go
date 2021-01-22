package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"klog"
	"testing"
)

func TestSumUpTimes(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r.AddDuration(klog.NewDuration(1, 0), "")
	r.AddDuration(klog.NewDuration(2, 0), "")
	assert.Equal(t, klog.NewDuration(3, 0), Total(r))
}

func TestSumUpZeroIfNoTimesAvailable(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	assert.Equal(t, klog.NewDuration(0, 0), Total(r))
}

func TestSumUpRanges(t *testing.T) {
	range1 := klog.Ɀ_Range_(klog.Ɀ_Time_(9, 7), klog.Ɀ_Time_(12, 59))
	range2 := klog.Ɀ_Range_(klog.Ɀ_Time_(13, 49), klog.Ɀ_Time_(17, 12))
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r.AddRange(range1, "")
	r.AddRange(range2, "")
	assert.Equal(t, klog.NewDuration(7, 15), Total(r))
}

func TestSumUpTimesAndRanges(t *testing.T) {
	range1 := klog.Ɀ_Range_(klog.Ɀ_Time_(8, 0), klog.Ɀ_Time_(12, 0))
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r.AddDuration(klog.NewDuration(1, 33), "")
	r.AddRange(range1, "")
	assert.Equal(t, klog.NewDuration(5, 33), Total(r))
}

func TestHashTagAllEntriesAreReturnedIfMatchIsInSummary(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	_ = r.SetSummary("This and #that, and other stuff as well")
	r.AddDuration(klog.NewDuration(2, 0), "Foo")
	r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(13, 49), klog.Ɀ_Time_(17, 12)), "Bar")
	es, hasMatched := FindEntriesWithHashtags(klog.NewTagSet("that"), r)
	require.Len(t, es, 2)
	assert.True(t, hasMatched)
	assert.Equal(t, klog.Summary("Foo"), es[0].Summary())
	assert.Equal(t, klog.Summary("Bar"), es[1].Summary())
}

func TestHashTagReturnsEntriesThatMatch(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	_ = r.SetSummary("This and #that, and other stuff as well")
	r.AddDuration(klog.NewDuration(2, 0), "Foo #fizzbuzz")
	r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(13, 49), klog.Ɀ_Time_(17, 12)), "Bar #barbaz")
	es, hasMatched := FindEntriesWithHashtags(klog.NewTagSet("fizzbuzz"), r)
	require.Len(t, es, 1)
	assert.True(t, hasMatched)
	assert.Equal(t, klog.Summary("Foo #fizzbuzz"), es[0].Summary())
}

func TestFindFilterWithNoClauses(t *testing.T) {
	rs, es := FindFilter([]klog.Record{
		func() klog.Record {
			r := klog.NewRecord(klog.Ɀ_Date_(2000, 1, 15))
			r.AddDuration(klog.NewDuration(1, 0), "")
			return r
		}(),
		func() klog.Record {
			r := klog.NewRecord(klog.Ɀ_Date_(2000, 1, 16))
			r.AddDuration(klog.NewDuration(1, 0), "")
			return r
		}(),
		func() klog.Record {
			r := klog.NewRecord(klog.Ɀ_Date_(2000, 1, 16))
			r.AddDuration(klog.NewDuration(1, 0), "")
			return r
		}(),
	}, Filter{})

	require.Len(t, rs, 3)
	assert.Equal(t, klog.NewDuration(3, 0), TotalEntries(es))
}

func TestFindFilterWithAfter(t *testing.T) {
	rs, _ := FindFilter([]klog.Record{
		klog.NewRecord(klog.Ɀ_Date_(2000, 1, 15)),
		klog.NewRecord(klog.Ɀ_Date_(2000, 1, 16)),
		klog.NewRecord(klog.Ɀ_Date_(2000, 1, 17)),
	}, Filter{AfterEq: klog.Ɀ_Date_(2000, 1, 16)})

	require.Len(t, rs, 2)
	assert.Equal(t, 16, rs[0].Date().Day())
	assert.Equal(t, 17, rs[1].Date().Day())
}

func TestFindFilterWithBefore(t *testing.T) {
	rs, _ := FindFilter([]klog.Record{
		klog.NewRecord(klog.Ɀ_Date_(2000, 1, 15)),
		klog.NewRecord(klog.Ɀ_Date_(2000, 1, 16)),
		klog.NewRecord(klog.Ɀ_Date_(2000, 1, 17)),
	}, Filter{BeforeEq: klog.Ɀ_Date_(2000, 1, 16)})

	require.Len(t, rs, 2)
	assert.Equal(t, 15, rs[0].Date().Day())
	assert.Equal(t, 16, rs[1].Date().Day())
}

func TestFindFilterWithHash(t *testing.T) {
	rs, es := FindFilter([]klog.Record{
		func() klog.Record {
			r := klog.NewRecord(klog.Ɀ_Date_(2000, 1, 15))
			_ = r.SetSummary("Contains #foo")
			r.AddDuration(klog.NewDuration(3, 0), "")
			return r
		}(),
		func() klog.Record {
			r := klog.NewRecord(klog.Ɀ_Date_(2000, 1, 16))
			r.AddDuration(klog.NewDuration(1, 0), "Contains #foo too")
			return r
		}(),
		func() klog.Record {
			r := klog.NewRecord(klog.Ɀ_Date_(2000, 1, 17))
			return r
		}(),
	}, Filter{Tags: []string{"foo"}})

	require.Len(t, rs, 2)
	require.Len(t, es, 2)
	assert.Equal(t, 15, rs[0].Date().Day())
	assert.Equal(t, 16, rs[1].Date().Day())
	assert.Equal(t, klog.NewDuration(4, 0), TotalEntries(es))
}
