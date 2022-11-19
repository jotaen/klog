package service

import (
	"github.com/jotaen/klog/klog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAggregateTotalTimesByTag(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r.SetSummary(klog.Ɀ_RecordSummary_("#foo"))
	r.AddDuration(klog.NewDuration(1, 0), klog.Ɀ_EntrySummary_("#foo=1"))
	r.AddDuration(klog.NewDuration(3, 0), klog.Ɀ_EntrySummary_("#foo"))
	r.AddDuration(klog.NewDuration(0, 30), klog.Ɀ_EntrySummary_("#test"))

	totals := AggregateTotalsByTags(r)
	require.Len(t, totals, 3)

	i := 0
	assert.Equal(t, klog.NewTagOrPanic("foo", ""), totals[i].Tag)
	assert.Equal(t, klog.NewDuration(4, 30), totals[i].Total)
	assert.Equal(t, 3, totals[i].Count)

	i++
	assert.Equal(t, klog.NewTagOrPanic("foo", "1"), totals[i].Tag)
	assert.Equal(t, klog.NewDuration(1, 0), totals[i].Total)
	assert.Equal(t, 1, totals[i].Count)

	i++
	assert.Equal(t, klog.NewTagOrPanic("test", ""), totals[i].Tag)
	assert.Equal(t, klog.NewDuration(0, 30), totals[i].Total)
	assert.Equal(t, 1, totals[i].Count)
}

func TestAggregateTotalIgnoresRedundantTags(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r.SetSummary(klog.Ɀ_RecordSummary_("#foo #foo #foo=1"))
	r.AddDuration(klog.NewDuration(1, 0), klog.Ɀ_EntrySummary_("#foo=1 #foo"))
	r.AddDuration(klog.NewDuration(3, 0), klog.Ɀ_EntrySummary_("#foo=2 #foo=1 #foo"))

	totals := AggregateTotalsByTags(r)
	require.Len(t, totals, 3)

	i := 0
	assert.Equal(t, klog.NewTagOrPanic("foo", ""), totals[i].Tag)
	assert.Equal(t, klog.NewDuration(4, 0), totals[i].Total)
	assert.Equal(t, 2, totals[i].Count)

	i++
	assert.Equal(t, klog.NewTagOrPanic("foo", "1"), totals[i].Tag)
	assert.Equal(t, klog.NewDuration(4, 0), totals[i].Total)
	assert.Equal(t, 2, totals[i].Count)

	i++
	assert.Equal(t, klog.NewTagOrPanic("foo", "2"), totals[i].Tag)
	assert.Equal(t, klog.NewDuration(3, 0), totals[i].Total)
	assert.Equal(t, 1, totals[i].Count)
}

func TestAggregateTotalTimesByTagSortsAlphabetically(t *testing.T) {
	r1 := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 1))
	r1.AddDuration(klog.NewDuration(1, 0), klog.Ɀ_EntrySummary_("#bbb"))
	r1.AddDuration(klog.NewDuration(3, 0), klog.Ɀ_EntrySummary_("#aaa"))
	r1.AddDuration(klog.NewDuration(3, 0), klog.Ɀ_EntrySummary_("#ddd"))

	r2 := klog.NewRecord(klog.Ɀ_Date_(2020, 1, 2))
	r2.AddDuration(klog.NewDuration(0, 30), klog.Ɀ_EntrySummary_("#ccc=1"))
	r2.AddDuration(klog.NewDuration(0, 30), klog.Ɀ_EntrySummary_("#ccc=2"))

	totals := AggregateTotalsByTags(r1, r2)
	require.Len(t, totals, 6)

	i := 0
	assert.Equal(t, klog.NewTagOrPanic("aaa", ""), totals[i].Tag)
	i += 1
	assert.Equal(t, klog.NewTagOrPanic("bbb", ""), totals[i].Tag)
	i += 1
	assert.Equal(t, klog.NewTagOrPanic("ccc", ""), totals[i].Tag)
	i += 1
	assert.Equal(t, klog.NewTagOrPanic("ccc", "1"), totals[i].Tag)
	i += 1
	assert.Equal(t, klog.NewTagOrPanic("ccc", "2"), totals[i].Tag)
	i += 1
	assert.Equal(t, klog.NewTagOrPanic("ddd", ""), totals[i].Tag)
}
