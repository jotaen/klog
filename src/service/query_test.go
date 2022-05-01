package service

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/service/period"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func sampleRecordsForQuerying() []Record {
	return []Record{
		func() Record {
			r := NewRecord(Ɀ_Date_(1999, 12, 30))
			r.SetSummary(Ɀ_RecordSummary_("Hello World", "#foo"))
			return r
		}(), func() Record {
			r := NewRecord(Ɀ_Date_(1999, 12, 31))
			r.AddDuration(NewDuration(5, 0), Ɀ_EntrySummary_("#bar"))
			return r
		}(), func() Record {
			r := NewRecord(Ɀ_Date_(2000, 1, 1))
			r.SetSummary(Ɀ_RecordSummary_("#foo"))
			r.AddDuration(NewDuration(0, 15), nil)
			r.AddDuration(NewDuration(6, 0), Ɀ_EntrySummary_("#bar"))
			r.AddDuration(NewDuration(0, -30), nil)
			return r
		}(), func() Record {
			r := NewRecord(Ɀ_Date_(2000, 1, 2))
			r.SetSummary(Ɀ_RecordSummary_("#foo"))
			r.AddDuration(NewDuration(7, 0), nil)
			return r
		}(), func() Record {
			r := NewRecord(Ɀ_Date_(2000, 1, 3))
			r.SetSummary(Ɀ_RecordSummary_("#foo=a"))
			r.AddDuration(NewDuration(4, 0), Ɀ_EntrySummary_("test", "foo #bar=1"))
			r.AddDuration(NewDuration(4, 0), Ɀ_EntrySummary_("#bar=2"))
			return r
		}(),
	}
}

func TestQueryWithNoClauses(t *testing.T) {
	qry := Query{}
	rs := Filter(qry.ToMatcher(), sampleRecordsForQuerying())
	require.Len(t, rs, 5)
	assert.Equal(t, NewDuration(5+6+7+8, -30+15), Total(rs...))
}

func TestQueryWithAtDate(t *testing.T) {
	qry := Query{AtDate: Ɀ_Date_(2000, 1, 2)}
	rs := Filter(qry.ToMatcher(), sampleRecordsForQuerying())
	require.Len(t, rs, 1)
	assert.Equal(t, NewDuration(7, 0), Total(rs...))
}

func TestQueryWithAfter(t *testing.T) {
	qry := Query{FromDate: Ɀ_Date_(2000, 1, 1)}
	rs := Filter(qry.ToMatcher(), sampleRecordsForQuerying())
	require.Len(t, rs, 3)
	assert.Equal(t, 1, rs[0].Date().Day())
	assert.Equal(t, 2, rs[1].Date().Day())
	assert.Equal(t, 3, rs[2].Date().Day())
}

func TestQueryWithBefore(t *testing.T) {
	qry := Query{UpToDate: Ɀ_Date_(2000, 1, 1)}
	rs := Filter(qry.ToMatcher(), sampleRecordsForQuerying())
	require.Len(t, rs, 3)
	assert.Equal(t, 30, rs[0].Date().Day())
	assert.Equal(t, 31, rs[1].Date().Day())
	assert.Equal(t, 1, rs[2].Date().Day())
}

func TestQueryInPeriod(t *testing.T) {
	qry := Query{InPeriod: []period.Period{period.NewPeriod(Ɀ_Date_(2000, 1, 1), Ɀ_Date_(2000, 1, 31))}}
	rs := Filter(qry.ToMatcher(), sampleRecordsForQuerying())
	require.Len(t, rs, 3)
	assert.Equal(t, 1, rs[0].Date().Day())
	assert.Equal(t, 2, rs[1].Date().Day())
	assert.Equal(t, 3, rs[2].Date().Day())
}

func TestQueryWithTagOnEntries(t *testing.T) {
	qry := Query{WithTags: []Tag{NewTagOrPanic("bar", "")}}
	rs := Filter(qry.ToMatcher(), sampleRecordsForQuerying())
	require.Len(t, rs, 3)
	assert.Equal(t, 31, rs[0].Date().Day())
	assert.Equal(t, 1, rs[1].Date().Day())
	assert.Equal(t, 3, rs[2].Date().Day())
	assert.Equal(t, NewDuration(5+8+6, 0), Total(rs...))
}

func TestQueryWithTagOnOverallSummary(t *testing.T) {
	qry := Query{WithTags: []Tag{NewTagOrPanic("foo", "")}}
	rs := Filter(qry.ToMatcher(), sampleRecordsForQuerying())
	require.Len(t, rs, 4)
	assert.Equal(t, 30, rs[0].Date().Day())
	assert.Equal(t, 1, rs[1].Date().Day())
	assert.Equal(t, 2, rs[2].Date().Day())
	assert.Equal(t, 3, rs[3].Date().Day())
	assert.Equal(t, NewDuration(6+7+8, -30+15), Total(rs...))
}

func TestQueryWithTagOnEntriesAndInSummary(t *testing.T) {
	qry := Query{WithTags: []Tag{NewTagOrPanic("foo", ""), NewTagOrPanic("bar", "")}}
	rs := Filter(qry.ToMatcher(), sampleRecordsForQuerying())
	require.Len(t, rs, 2)
	assert.Equal(t, 1, rs[0].Date().Day())
	assert.Equal(t, 3, rs[1].Date().Day())
	assert.Equal(t, NewDuration(8+6, 0), Total(rs...))
}

func TestQueryWithTagValues(t *testing.T) {
	qry := Query{WithTags: []Tag{NewTagOrPanic("foo", "a")}}
	rs := Filter(qry.ToMatcher(), sampleRecordsForQuerying())
	require.Len(t, rs, 1)
	assert.Equal(t, 3, rs[0].Date().Day())
	assert.Equal(t, NewDuration(8, 0), Total(rs...))
}

func TestQueryWithTagValuesInEntries(t *testing.T) {
	qry := Query{WithTags: []Tag{NewTagOrPanic("bar", "1")}}
	rs := Filter(qry.ToMatcher(), sampleRecordsForQuerying())
	require.Len(t, rs, 1)
	assert.Equal(t, 3, rs[0].Date().Day())
	assert.Equal(t, NewDuration(4, 0), Total(rs...))
}

func TestQueryWithTagNonMatchingValues(t *testing.T) {
	qry := Query{WithTags: []Tag{NewTagOrPanic("bar", "3")}}
	rs := Filter(qry.ToMatcher(), sampleRecordsForQuerying())
	require.Len(t, rs, 0)
}
