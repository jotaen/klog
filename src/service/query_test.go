package service

import (
	. "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func sampleRecordsForQuerying() []Record {
	return []Record{
		func() Record {
			r := NewRecord(Ɀ_Date_(1999, 12, 30))
			_ = r.SetSummary(NewSummary("#foo"))
			return r
		}(), func() Record {
			r := NewRecord(Ɀ_Date_(1999, 12, 31))
			r.AddDuration(NewDuration(5, 0), NewSummary("#bar"))
			return r
		}(), func() Record {
			r := NewRecord(Ɀ_Date_(2000, 1, 1))
			_ = r.SetSummary(NewSummary("#foo"))
			r.AddDuration(NewDuration(0, 15), NewSummary())
			r.AddDuration(NewDuration(6, 0), NewSummary("#bar"))
			r.AddDuration(NewDuration(0, -30), NewSummary())
			return r
		}(), func() Record {
			r := NewRecord(Ɀ_Date_(2000, 1, 2))
			_ = r.SetSummary(NewSummary("#foo"))
			r.AddDuration(NewDuration(7, 0), NewSummary())
			return r
		}(), func() Record {
			r := NewRecord(Ɀ_Date_(2000, 1, 3))
			_ = r.SetSummary(NewSummary("#foo"))
			r.AddDuration(NewDuration(4, 0), NewSummary("#bar"))
			r.AddDuration(NewDuration(4, 0), NewSummary("#bar"))
			return r
		}(),
	}
}

func TestQueryWithNoClauses(t *testing.T) {
	rs := Filter(sampleRecordsForQuerying(), FilterQry{})
	require.Len(t, rs, 5)
	assert.Equal(t, NewDuration(5+6+7+8, -30+15), Total(rs...))
}

func TestQueryWithAfter(t *testing.T) {
	rs := Filter(sampleRecordsForQuerying(), FilterQry{AfterOrEqual: Ɀ_Date_(2000, 1, 1)})
	require.Len(t, rs, 3)
	assert.Equal(t, 1, rs[0].Date().Day())
	assert.Equal(t, 2, rs[1].Date().Day())
	assert.Equal(t, 3, rs[2].Date().Day())
}

func TestQueryWithBefore(t *testing.T) {
	rs := Filter(sampleRecordsForQuerying(), FilterQry{BeforeOrEqual: Ɀ_Date_(2000, 1, 1)})
	require.Len(t, rs, 3)
	assert.Equal(t, 30, rs[0].Date().Day())
	assert.Equal(t, 31, rs[1].Date().Day())
	assert.Equal(t, 1, rs[2].Date().Day())
}

func TestQueryWithTagOnEntries(t *testing.T) {
	rs := Filter(sampleRecordsForQuerying(), FilterQry{Tags: []string{"bar"}})
	require.Len(t, rs, 3)
	assert.Equal(t, 31, rs[0].Date().Day())
	assert.Equal(t, 1, rs[1].Date().Day())
	assert.Equal(t, 3, rs[2].Date().Day())
	assert.Equal(t, NewDuration(5+8+6, 0), Total(rs...))
}

func TestQueryWithTagOnOverallSummary(t *testing.T) {
	rs := Filter(sampleRecordsForQuerying(), FilterQry{Tags: []string{"foo"}})
	require.Len(t, rs, 4)
	assert.Equal(t, 30, rs[0].Date().Day())
	assert.Equal(t, 1, rs[1].Date().Day())
	assert.Equal(t, 2, rs[2].Date().Day())
	assert.Equal(t, 3, rs[3].Date().Day())
	assert.Equal(t, NewDuration(6+7+8, -30+15), Total(rs...))
}

func TestQueryWithTagOnEntriesAndInSummary(t *testing.T) {
	rs := Filter(sampleRecordsForQuerying(), FilterQry{Tags: []string{"foo", "bar"}})
	require.Len(t, rs, 2)
	assert.Equal(t, 1, rs[0].Date().Day())
	assert.Equal(t, 3, rs[1].Date().Day())
	assert.Equal(t, NewDuration(8+6, 0), Total(rs...))
}

func TestQueryWithFuzzyTags(t *testing.T) {
	rs := Filter(sampleRecordsForQuerying(), FilterQry{Tags: []string{"fo..."}})
	require.Len(t, rs, 4)
	assert.Equal(t, 30, rs[0].Date().Day())
	assert.Equal(t, 1, rs[1].Date().Day())
	assert.Equal(t, 2, rs[2].Date().Day())
	assert.Equal(t, 3, rs[3].Date().Day())
}

func TestQueryWithSorting(t *testing.T) {
	ss := sampleRecordsForQuerying()
	for _, x := range []struct{ rs []Record }{
		{ss},
		{[]Record{ss[3], ss[1], ss[2], ss[0], ss[4]}},
		{[]Record{ss[1], ss[4], ss[0], ss[3], ss[2]}},
	} {
		ascending := Sort(x.rs, true)
		assert.Equal(t, []Record{ss[0], ss[1], ss[2], ss[3], ss[4]}, ascending)

		descending := Sort(x.rs, false)
		assert.Equal(t, []Record{ss[4], ss[3], ss[2], ss[1], ss[0]}, descending)
	}
}
