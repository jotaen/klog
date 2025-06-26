package service

import (
	"github.com/jotaen/klog/klog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAggregateTotalTimesByTag(t *testing.T) {
	rs := []klog.Record{
		func() klog.Record {
			r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
			r.SetSummary(klog.Ɀ_RecordSummary_("#foo"))
			r.AddDuration(klog.NewDuration(1, 0), klog.Ɀ_EntrySummary_("#foo=1"))
			r.AddDuration(klog.NewDuration(3, 0), klog.Ɀ_EntrySummary_("#foo"))
			r.AddDuration(klog.NewDuration(0, 30), klog.Ɀ_EntrySummary_("#test"))
			return r
		}(),
		func() klog.Record {
			r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 2))
			r.AddDuration(klog.NewDuration(1, 0), klog.Ɀ_EntrySummary_("#foo=2"))
			r.AddDuration(klog.NewDuration(8, 0), klog.Ɀ_EntrySummary_("#bar"))
			r.AddDuration(klog.NewDuration(0, 45), klog.Ɀ_EntrySummary_("no tag"))
			return r
		}(),
	}

	tagStats, untagged := AggregateTotalsByTags(rs...)
	require.Len(t, tagStats, 5)

	assert.Equal(t, klog.NewDuration(0, 45), untagged.Total)
	assert.Equal(t, 1, untagged.Count)

	i := 0
	assert.Equal(t, klog.NewTagOrPanic("bar", ""), tagStats[i].Tag)
	assert.Equal(t, klog.NewDuration(8, 0), tagStats[i].Total)
	assert.Equal(t, 1, tagStats[i].Count)

	i++
	assert.Equal(t, klog.NewTagOrPanic("foo", ""), tagStats[i].Tag)
	assert.Equal(t, klog.NewDuration(5, 30), tagStats[i].Total)
	assert.Equal(t, 4, tagStats[i].Count)

	i++
	assert.Equal(t, klog.NewTagOrPanic("foo", "1"), tagStats[i].Tag)
	assert.Equal(t, klog.NewDuration(1, 0), tagStats[i].Total)
	assert.Equal(t, 1, tagStats[i].Count)

	i++
	assert.Equal(t, klog.NewTagOrPanic("foo", "2"), tagStats[i].Tag)
	assert.Equal(t, klog.NewDuration(1, 0), tagStats[i].Total)
	assert.Equal(t, 1, tagStats[i].Count)

	i++
	assert.Equal(t, klog.NewTagOrPanic("test", ""), tagStats[i].Tag)
	assert.Equal(t, klog.NewDuration(0, 30), tagStats[i].Total)
	assert.Equal(t, 1, tagStats[i].Count)
}

func TestAggregateTotalIgnoresRedundantTags(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r.SetSummary(klog.Ɀ_RecordSummary_("#foo #foo #foo=1"))
	r.AddDuration(klog.NewDuration(1, 0), klog.Ɀ_EntrySummary_("#foo=1 #foo"))
	r.AddDuration(klog.NewDuration(3, 0), klog.Ɀ_EntrySummary_("#foo=2 #foo=1 #foo"))

	tagStats, _ := AggregateTotalsByTags(r)
	require.Len(t, tagStats, 3)

	i := 0
	assert.Equal(t, klog.NewTagOrPanic("foo", ""), tagStats[i].Tag)
	assert.Equal(t, klog.NewDuration(4, 0), tagStats[i].Total)
	assert.Equal(t, 2, tagStats[i].Count)

	i++
	assert.Equal(t, klog.NewTagOrPanic("foo", "1"), tagStats[i].Tag)
	assert.Equal(t, klog.NewDuration(4, 0), tagStats[i].Total)
	assert.Equal(t, 2, tagStats[i].Count)

	i++
	assert.Equal(t, klog.NewTagOrPanic("foo", "2"), tagStats[i].Tag)
	assert.Equal(t, klog.NewDuration(3, 0), tagStats[i].Total)
	assert.Equal(t, 1, tagStats[i].Count)
}

func TestAggregateTotalTimesByTagSortsAlphabetically(t *testing.T) {
	r1 := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r1.AddDuration(klog.NewDuration(1, 0), klog.Ɀ_EntrySummary_("#bbb"))
	r1.AddDuration(klog.NewDuration(3, 0), klog.Ɀ_EntrySummary_("#aaa"))
	r1.AddDuration(klog.NewDuration(3, 0), klog.Ɀ_EntrySummary_("#ddd"))

	r2 := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 2))
	r2.AddDuration(klog.NewDuration(0, 30), klog.Ɀ_EntrySummary_("#ccc=1"))
	r2.AddDuration(klog.NewDuration(0, 30), klog.Ɀ_EntrySummary_("#ccc=2"))

	tagStats, _ := AggregateTotalsByTags(r1, r2)
	require.Len(t, tagStats, 6)

	i := 0
	assert.Equal(t, klog.NewTagOrPanic("aaa", ""), tagStats[i].Tag)
	i += 1
	assert.Equal(t, klog.NewTagOrPanic("bbb", ""), tagStats[i].Tag)
	i += 1
	assert.Equal(t, klog.NewTagOrPanic("ccc", ""), tagStats[i].Tag)
	i += 1
	assert.Equal(t, klog.NewTagOrPanic("ccc", "1"), tagStats[i].Tag)
	i += 1
	assert.Equal(t, klog.NewTagOrPanic("ccc", "2"), tagStats[i].Tag)
	i += 1
	assert.Equal(t, klog.NewTagOrPanic("ddd", ""), tagStats[i].Tag)
}
