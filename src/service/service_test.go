package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "klog/record"
	"testing"
)

func TestSumUpTimes(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.AddDuration(NewDuration(1, 0), "")
	r.AddDuration(NewDuration(2, 0), "")
	assert.Equal(t, NewDuration(3, 0), Total(r))
}

func TestSumUpZeroIfNoTimesAvailable(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	assert.Equal(t, NewDuration(0, 0), Total(r))
}

func TestSumUpRanges(t *testing.T) {
	range1 := Ɀ_Range_(Ɀ_Time_(9, 7), Ɀ_Time_(12, 59))
	range2 := Ɀ_Range_(Ɀ_Time_(13, 49), Ɀ_Time_(17, 12))
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.AddRange(range1, "")
	r.AddRange(range2, "")
	assert.Equal(t, NewDuration(7, 15), Total(r))
}

func TestSumUpTimesAndRanges(t *testing.T) {
	range1 := Ɀ_Range_(Ɀ_Time_(8, 0), Ɀ_Time_(12, 0))
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.AddDuration(NewDuration(1, 33), "")
	r.AddRange(range1, "")
	assert.Equal(t, NewDuration(5, 33), Total(r))
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

func TestFindFilterWithNoClauses(t *testing.T) {
	rs, es := FindFilter([]Record{
		func() Record {
			r := NewRecord(Ɀ_Date_(2000, 1, 15))
			r.AddDuration(NewDuration(1, 0), "")
			return r
		}(),
		func() Record {
			r := NewRecord(Ɀ_Date_(2000, 1, 16))
			r.AddDuration(NewDuration(1, 0), "")
			return r
		}(),
		func() Record {
			r := NewRecord(Ɀ_Date_(2000, 1, 16))
			r.AddDuration(NewDuration(1, 0), "")
			return r
		}(),
	}, Filter{})

	require.Len(t, rs, 3)
	assert.Equal(t, NewDuration(3, 0), TotalEntries(es))
}

func TestFindFilterWithAfter(t *testing.T) {
	rs, _ := FindFilter([]Record{
		NewRecord(Ɀ_Date_(2000, 1, 15)),
		NewRecord(Ɀ_Date_(2000, 1, 16)),
		NewRecord(Ɀ_Date_(2000, 1, 17)),
	}, Filter{AfterEq: Ɀ_Date_(2000, 1, 16)})

	require.Len(t, rs, 2)
	assert.Equal(t, 16, rs[0].Date().Day())
	assert.Equal(t, 17, rs[1].Date().Day())
}

func TestFindFilterWithBefore(t *testing.T) {
	rs, _ := FindFilter([]Record{
		NewRecord(Ɀ_Date_(2000, 1, 15)),
		NewRecord(Ɀ_Date_(2000, 1, 16)),
		NewRecord(Ɀ_Date_(2000, 1, 17)),
	}, Filter{BeforeEq: Ɀ_Date_(2000, 1, 16)})

	require.Len(t, rs, 2)
	assert.Equal(t, 15, rs[0].Date().Day())
	assert.Equal(t, 16, rs[1].Date().Day())
}

func TestFindFilterWithHash(t *testing.T) {
	rs, es := FindFilter([]Record{
		func() Record {
			r := NewRecord(Ɀ_Date_(2000, 1, 15))
			_ = r.SetSummary("Contains #foo")
			r.AddDuration(NewDuration(3, 0), "")
			return r
		}(),
		func() Record {
			r := NewRecord(Ɀ_Date_(2000, 1, 16))
			r.AddDuration(NewDuration(1, 0), "Contains #foo too")
			return r
		}(),
		func() Record {
			r := NewRecord(Ɀ_Date_(2000, 1, 17))
			return r
		}(),
	}, Filter{Tags: []string{"foo"}})

	require.Len(t, rs, 2)
	require.Len(t, es, 2)
	assert.Equal(t, 15, rs[0].Date().Day())
	assert.Equal(t, 16, rs[1].Date().Day())
	assert.Equal(t, NewDuration(4, 0), TotalEntries(es))
}
