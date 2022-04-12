package service

import (
	. "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAggregateTotalTimesByTagSortsAlphabetically(t *testing.T) {
	r1 := NewRecord(Ɀ_Date_(2020, 1, 1))
	r1.AddDuration(NewDuration(1, 0), Ɀ_EntrySummary_("#bbb"))
	r1.AddDuration(NewDuration(3, 0), Ɀ_EntrySummary_("#aaa"))
	r1.AddDuration(NewDuration(3, 0), Ɀ_EntrySummary_("#ddd"))

	r2 := NewRecord(Ɀ_Date_(2020, 1, 2))
	r2.AddDuration(NewDuration(0, 30), Ɀ_EntrySummary_("#ccc=1"))
	r2.AddDuration(NewDuration(0, 30), Ɀ_EntrySummary_("#ccc=2"))

	totals := AggregateTotalsByTags(r1, r2)
	require.Len(t, totals, 6)

	i := 0
	assert.Equal(t, NewTagOrPanic("aaa", ""), totals[i].Tag)
	i += 1
	assert.Equal(t, NewTagOrPanic("bbb", ""), totals[i].Tag)
	i += 1
	assert.Equal(t, NewTagOrPanic("ccc", ""), totals[i].Tag)
	i += 1
	assert.Equal(t, NewTagOrPanic("ccc", "1"), totals[i].Tag)
	i += 1
	assert.Equal(t, NewTagOrPanic("ccc", "2"), totals[i].Tag)
	i += 1
	assert.Equal(t, NewTagOrPanic("ddd", ""), totals[i].Tag)
}

func TestAggregateTotalTimesByTag(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.SetSummary(Ɀ_RecordSummary_("#foo"))
	r.AddDuration(NewDuration(1, 0), Ɀ_EntrySummary_("#foo=1"))
	r.AddDuration(NewDuration(3, 0), Ɀ_EntrySummary_("#foo"))

	totals := AggregateTotalsByTags(r)
	require.Len(t, totals, 2)
	i := 0

	assert.Equal(t, NewTagOrPanic("foo", ""), totals[i].Tag)
	assert.Equal(t, NewDuration(4, 0), totals[i].Total)
	i += 1
	assert.Equal(t, NewTagOrPanic("foo", "1"), totals[i].Tag)
	assert.Equal(t, NewDuration(1, 0), totals[i].Total)
}

func TestAggregateTotalIgnoresRedundantTags(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.SetSummary(Ɀ_RecordSummary_("#foo #foo #foo=1"))
	r.AddDuration(NewDuration(1, 0), Ɀ_EntrySummary_("#foo=1 #foo"))
	r.AddDuration(NewDuration(3, 0), Ɀ_EntrySummary_("#foo=2 #foo=1 #foo"))

	totals := AggregateTotalsByTags(r)
	require.Len(t, totals, 3)
	i := 0

	assert.Equal(t, NewTagOrPanic("foo", ""), totals[i].Tag)
	assert.Equal(t, NewDuration(4, 0), totals[i].Total)
	i += 1
	assert.Equal(t, NewTagOrPanic("foo", "1"), totals[i].Tag)
	assert.Equal(t, NewDuration(4, 0), totals[i].Total)
	i += 1
	assert.Equal(t, NewTagOrPanic("foo", "2"), totals[i].Tag)
	assert.Equal(t, NewDuration(3, 0), totals[i].Total)
}
