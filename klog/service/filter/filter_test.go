package filter

import (
	"testing"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleRecordsForQuerying() []klog.Record {
	return []klog.Record{
		func() klog.Record {
			// Note that records without entries never match any query.
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
	rs := Filter(And{}, sampleRecordsForQuerying())
	require.Len(t, rs, 4)
	assert.Equal(t, klog.NewDuration(5+6+7+8, -30+15), service.Total(rs...))
}

func TestQueryWithAtDate(t *testing.T) {
	rs := Filter(IsInDateRange{
		From: klog.Ɀ_Date_(2000, 1, 2),
		To:   klog.Ɀ_Date_(2000, 1, 2),
	}, sampleRecordsForQuerying())
	require.Len(t, rs, 1)
	assert.Equal(t, klog.NewDuration(7, 0), service.Total(rs...))
}

func TestQueryWithAfter(t *testing.T) {
	rs := Filter(IsInDateRange{
		From: klog.Ɀ_Date_(2000, 1, 1),
		To:   nil,
	}, sampleRecordsForQuerying())
	require.Len(t, rs, 3)
	assert.Equal(t, 1, rs[0].Date().Day())
	assert.Equal(t, 2, rs[1].Date().Day())
	assert.Equal(t, 3, rs[2].Date().Day())
}

func TestQueryWithBefore(t *testing.T) {
	rs := Filter(IsInDateRange{
		From: nil,
		To:   klog.Ɀ_Date_(2000, 1, 1),
	}, sampleRecordsForQuerying())
	require.Len(t, rs, 2)
	assert.Equal(t, 31, rs[0].Date().Day())
	assert.Equal(t, 1, rs[1].Date().Day())
}

func TestQueryWithTagOnEntries(t *testing.T) {
	rs := Filter(HasTag{klog.NewTagOrPanic("bar", "")}, sampleRecordsForQuerying())
	require.Len(t, rs, 3)
	assert.Equal(t, 31, rs[0].Date().Day())
	assert.Equal(t, 1, rs[1].Date().Day())
	assert.Equal(t, 3, rs[2].Date().Day())
	assert.Equal(t, klog.NewDuration(5+8+6, 0), service.Total(rs...))
}

func TestQueryWithTagOnOverallSummary(t *testing.T) {
	rs := Filter(HasTag{klog.NewTagOrPanic("foo", "")}, sampleRecordsForQuerying())
	require.Len(t, rs, 3)
	assert.Equal(t, 1, rs[0].Date().Day())
	assert.Equal(t, 2, rs[1].Date().Day())
	assert.Equal(t, 3, rs[2].Date().Day())
	assert.Equal(t, klog.NewDuration(6+7+8, -30+15), service.Total(rs...))
}

func TestQueryWithTagOnEntriesAndInSummary(t *testing.T) {
	rs := Filter(And{[]Predicate{HasTag{klog.NewTagOrPanic("foo", "")}, HasTag{klog.NewTagOrPanic("bar", "")}}}, sampleRecordsForQuerying())
	require.Len(t, rs, 2)
	assert.Equal(t, 1, rs[0].Date().Day())
	assert.Equal(t, 3, rs[1].Date().Day())
	assert.Equal(t, klog.NewDuration(8+6, 0), service.Total(rs...))
}

func TestQueryWithTagValues(t *testing.T) {
	rs := Filter(HasTag{klog.NewTagOrPanic("foo", "a")}, sampleRecordsForQuerying())
	require.Len(t, rs, 1)
	assert.Equal(t, 3, rs[0].Date().Day())
	assert.Equal(t, klog.NewDuration(8, 0), service.Total(rs...))
}

func TestQueryWithTagValuesInEntries(t *testing.T) {
	rs := Filter(HasTag{klog.NewTagOrPanic("bar", "1")}, sampleRecordsForQuerying())
	require.Len(t, rs, 1)
	assert.Equal(t, 3, rs[0].Date().Day())
	assert.Equal(t, klog.NewDuration(4, 0), service.Total(rs...))
}

func TestQueryWithTagNonMatchingValues(t *testing.T) {
	rs := Filter(HasTag{klog.NewTagOrPanic("bar", "3")}, sampleRecordsForQuerying())
	require.Len(t, rs, 0)
}

func TestQueryWithEntryTypes(t *testing.T) {
	{
		rs := Filter(IsEntryType{ENTRY_TYPE_DURATION}, sampleRecordsForQuerying())
		require.Len(t, rs, 4)
		assert.Equal(t, klog.NewDuration(0, 1545), service.Total(rs...))
	}
	{
		rs := Filter(IsEntryType{ENTRY_TYPE_DURATION_NEGATIVE}, sampleRecordsForQuerying())
		require.Len(t, rs, 1)
		assert.Equal(t, 1, rs[0].Date().Day())
		assert.Equal(t, klog.NewDuration(0, -30), service.Total(rs...))
	}
	{
		rs := Filter(IsEntryType{ENTRY_TYPE_DURATION_POSITIVE}, sampleRecordsForQuerying())
		require.Len(t, rs, 4)
		assert.Equal(t, klog.NewDuration(0, 1575), service.Total(rs...))
	}
	{
		rs := Filter(IsEntryType{ENTRY_TYPE_RANGE}, sampleRecordsForQuerying())
		require.Len(t, rs, 0)
		assert.Equal(t, klog.NewDuration(0, 0), service.Total(rs...))
	}
	{
		rs := Filter(IsEntryType{ENTRY_TYPE_OPEN_RANGE}, sampleRecordsForQuerying())
		require.Len(t, rs, 1)
		assert.Equal(t, klog.NewDuration(0, 0), service.Total(rs...))
	}
}
