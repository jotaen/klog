package service

import (
	"github.com/jotaen/klog/klog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func sampleRecordsForQuerying() []klog.Record {
	return []klog.Record{
		func() klog.Record {
			r := klog.NewRecord(klog.Ɀ_Date_(1999, 12, 30))
			r.SetSummary(klog.Ɀ_RecordSummary_("Hello World", "#foo"))
			return r
		}(), func() klog.Record {
			r := klog.NewRecord(klog.Ɀ_Date_(1999, 12, 31))
			r.AddDuration(klog.NewDuration(5, 0), klog.Ɀ_EntrySummary_("#bar"))
			return r
		}(), func() klog.Record {
			r := klog.NewRecord(klog.Ɀ_Date_(2000, 1, 1))
			r.SetSummary(klog.Ɀ_RecordSummary_("#foo"))
			r.AddDuration(klog.NewDuration(0, 15), nil)
			r.AddDuration(klog.NewDuration(6, 0), klog.Ɀ_EntrySummary_("#bar"))
			r.AddDuration(klog.NewDuration(0, -30), nil)
			return r
		}(), func() klog.Record {
			r := klog.NewRecord(klog.Ɀ_Date_(2000, 1, 2))
			r.SetSummary(klog.Ɀ_RecordSummary_("#foo"))
			r.AddDuration(klog.NewDuration(7, 0), nil)
			return r
		}(), func() klog.Record {
			r := klog.NewRecord(klog.Ɀ_Date_(2000, 1, 3))
			r.SetSummary(klog.Ɀ_RecordSummary_("#foo=a"))
			r.AddDuration(klog.NewDuration(4, 0), klog.Ɀ_EntrySummary_("test", "foo #bar=1"))
			r.AddDuration(klog.NewDuration(4, 0), klog.Ɀ_EntrySummary_("#bar=2"))
			r.Start(klog.NewOpenRange(klog.Ɀ_Time_(12, 00)), nil)
			return r
		}(),
	}
}

func TestQueryWithNoClauses(t *testing.T) {
	rs := Filter(sampleRecordsForQuerying(), FilterQry{})
	require.Len(t, rs, 5)
	assert.Equal(t, klog.NewDuration(5+6+7+8, -30+15), Total(rs...))
}

func TestQueryWithAtDate(t *testing.T) {
	rs := Filter(sampleRecordsForQuerying(), FilterQry{AtDate: klog.Ɀ_Date_(2000, 1, 2)})
	require.Len(t, rs, 1)
	assert.Equal(t, klog.NewDuration(7, 0), Total(rs...))
}

func TestQueryWithAfter(t *testing.T) {
	rs := Filter(sampleRecordsForQuerying(), FilterQry{AfterOrEqual: klog.Ɀ_Date_(2000, 1, 1)})
	require.Len(t, rs, 3)
	assert.Equal(t, 1, rs[0].Date().Day())
	assert.Equal(t, 2, rs[1].Date().Day())
	assert.Equal(t, 3, rs[2].Date().Day())
}

func TestQueryWithBefore(t *testing.T) {
	rs := Filter(sampleRecordsForQuerying(), FilterQry{BeforeOrEqual: klog.Ɀ_Date_(2000, 1, 1)})
	require.Len(t, rs, 3)
	assert.Equal(t, 30, rs[0].Date().Day())
	assert.Equal(t, 31, rs[1].Date().Day())
	assert.Equal(t, 1, rs[2].Date().Day())
}

func TestQueryWithTagOnEntries(t *testing.T) {
	rs := Filter(sampleRecordsForQuerying(), FilterQry{Tags: []klog.Tag{klog.NewTagOrPanic("bar", "")}})
	require.Len(t, rs, 3)
	assert.Equal(t, 31, rs[0].Date().Day())
	assert.Equal(t, 1, rs[1].Date().Day())
	assert.Equal(t, 3, rs[2].Date().Day())
	assert.Equal(t, klog.NewDuration(5+8+6, 0), Total(rs...))
}

func TestQueryWithTagOnOverallSummary(t *testing.T) {
	rs := Filter(sampleRecordsForQuerying(), FilterQry{Tags: []klog.Tag{klog.NewTagOrPanic("foo", "")}})
	require.Len(t, rs, 4)
	assert.Equal(t, 30, rs[0].Date().Day())
	assert.Equal(t, 1, rs[1].Date().Day())
	assert.Equal(t, 2, rs[2].Date().Day())
	assert.Equal(t, 3, rs[3].Date().Day())
	assert.Equal(t, klog.NewDuration(6+7+8, -30+15), Total(rs...))
}

func TestQueryWithTagOnEntriesAndInSummary(t *testing.T) {
	rs := Filter(sampleRecordsForQuerying(), FilterQry{Tags: []klog.Tag{klog.NewTagOrPanic("foo", ""), klog.NewTagOrPanic("bar", "")}})
	require.Len(t, rs, 2)
	assert.Equal(t, 1, rs[0].Date().Day())
	assert.Equal(t, 3, rs[1].Date().Day())
	assert.Equal(t, klog.NewDuration(8+6, 0), Total(rs...))
}

func TestQueryWithTagValues(t *testing.T) {
	rs := Filter(sampleRecordsForQuerying(), FilterQry{Tags: []klog.Tag{klog.NewTagOrPanic("foo", "a")}})
	require.Len(t, rs, 1)
	assert.Equal(t, 3, rs[0].Date().Day())
	assert.Equal(t, klog.NewDuration(8, 0), Total(rs...))
}

func TestQueryWithTagValuesInEntries(t *testing.T) {
	rs := Filter(sampleRecordsForQuerying(), FilterQry{Tags: []klog.Tag{klog.NewTagOrPanic("bar", "1")}})
	require.Len(t, rs, 1)
	assert.Equal(t, 3, rs[0].Date().Day())
	assert.Equal(t, klog.NewDuration(4, 0), Total(rs...))
}

func TestQueryWithTagNonMatchingValues(t *testing.T) {
	rs := Filter(sampleRecordsForQuerying(), FilterQry{Tags: []klog.Tag{klog.NewTagOrPanic("bar", "3")}})
	require.Len(t, rs, 0)
}

func TestQueryWithEntryTypes(t *testing.T) {
	{
		rs := Filter(sampleRecordsForQuerying(), FilterQry{EntryType: ENTRY_TYPE_DURATION})
		require.Len(t, rs, 4)
		assert.Equal(t, klog.NewDuration(0, 1545), Total(rs...))
	}
	{
		rs := Filter(sampleRecordsForQuerying(), FilterQry{EntryType: ENTRY_TYPE_NEGATIVE_DURATION})
		require.Len(t, rs, 1)
		assert.Equal(t, 1, rs[0].Date().Day())
		assert.Equal(t, klog.NewDuration(0, -30), Total(rs...))
	}
	{
		rs := Filter(sampleRecordsForQuerying(), FilterQry{EntryType: ENTRY_TYPE_POSITIVE_DURATION})
		require.Len(t, rs, 4)
		assert.Equal(t, klog.NewDuration(0, 1575), Total(rs...))
	}
	{
		rs := Filter(sampleRecordsForQuerying(), FilterQry{EntryType: ENTRY_TYPE_RANGE})
		require.Len(t, rs, 0)
		assert.Equal(t, klog.NewDuration(0, 0), Total(rs...))
	}
	{
		rs := Filter(sampleRecordsForQuerying(), FilterQry{EntryType: ENTRY_TYPE_OPEN_RANGE})
		require.Len(t, rs, 1)
		assert.Equal(t, klog.NewDuration(0, 0), Total(rs...))
	}
}

func TestQueryWithSorting(t *testing.T) {
	ss := sampleRecordsForQuerying()
	for _, x := range []struct{ rs []klog.Record }{
		{ss},
		{[]klog.Record{ss[3], ss[1], ss[2], ss[0], ss[4]}},
		{[]klog.Record{ss[1], ss[4], ss[0], ss[3], ss[2]}},
	} {
		ascending := Sort(x.rs, true)
		assert.Equal(t, []klog.Record{ss[0], ss[1], ss[2], ss[3], ss[4]}, ascending)

		descending := Sort(x.rs, false)
		assert.Equal(t, []klog.Record{ss[4], ss[3], ss[2], ss[1], ss[0]}, descending)
	}
}
