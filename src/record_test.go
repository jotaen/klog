package klog

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInitialiseRecord(t *testing.T) {
	date := Ɀ_Date_(2020, 1, 1)
	r := NewRecord(date)
	assert.Equal(t, r.Date(), date)
	assert.Equal(t, NewDuration(0, 0).InMinutes(), r.ShouldTotal().InMinutes())
	assert.Equal(t, NewSummary(), r.Summary())
	assert.Len(t, r.Entries(), 0)
}

func TestSavesShouldTotal(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	assert.Equal(t, NewDuration(0, 0).InMinutes(), r.ShouldTotal().InMinutes())
	r.SetShouldTotal(NewDuration(5, 0))
	assert.Equal(t, NewDuration(5, 0).InMinutes(), r.ShouldTotal().InMinutes())
}

func TestAddRanges(t *testing.T) {
	range1 := Ɀ_Range_(Ɀ_Time_(9, 7), Ɀ_Time_(12, 59))
	range2 := Ɀ_Range_(Ɀ_Time_(13, 49), Ɀ_Time_(17, 12))
	w := NewRecord(Ɀ_Date_(2020, 1, 1))
	w.AddRange(range1, NewSummary("Range 1"))
	w.AddRange(range2, NewSummary("Range 2"))
	require.Len(t, w.Entries(), 2)
	assert.Equal(t, range1, w.Entries()[0].value)
	assert.Equal(t, NewSummary("Range 1"), w.Entries()[0].Summary())
	assert.Equal(t, range2, w.Entries()[1].value)
	assert.Equal(t, NewSummary("Range 2"), w.Entries()[1].Summary())
}

func TestStartOpenRange(t *testing.T) {
	time := Ɀ_Time_(11, 23)
	w := NewRecord(Ɀ_Date_(2020, 1, 1))
	assert.Equal(t, nil, w.OpenRange())
	_ = w.StartOpenRange(time, NewSummary("Open Range"))
	require.Len(t, w.Entries(), 1)
	assert.Equal(t, NewOpenRange(time), w.Entries()[0].value)
	assert.Equal(t, NewSummary("Open Range"), w.Entries()[0].Summary())
}

func TestCannotStartSecondOpenRange(t *testing.T) {
	time := Ɀ_Time_(11, 23)
	w := NewRecord(Ɀ_Date_(2020, 1, 1))
	assert.Equal(t, nil, w.OpenRange())
	_ = w.StartOpenRange(time, NewSummary("Open Range"))
	err := w.StartOpenRange(time, NewSummary("Open Range"))
	require.Error(t, err)
}

func TestCloseOpenRange(t *testing.T) {
	start := Ɀ_Time_(19, 22)
	w := NewRecord(Ɀ_Date_(2012, 6, 17))
	_ = w.StartOpenRange(start, NewSummary("Started"))
	end := Ɀ_Time_(20, 55)
	err := w.EndOpenRange(end)
	require.Nil(t, err)
	assert.Nil(t, w.OpenRange())
	require.Len(t, w.Entries(), 1)
	assert.Equal(t, Ɀ_Range_(start, end), w.Entries()[0].value)
	assert.Equal(t, NewSummary("Started"), w.Entries()[0].Summary())
}

func TestCloseOpenRangeFailsIfResultingRangeIsInvalid(t *testing.T) {
	start := Ɀ_Time_(19, 22)
	w := NewRecord(Ɀ_Date_(2012, 6, 17))
	_ = w.StartOpenRange(start, NewSummary("Started"))
	oldEntry := w.OpenRange()
	end := Ɀ_Time_(1, 30)
	err := w.EndOpenRange(end)
	require.Error(t, err)
	assert.Equal(t, oldEntry, w.OpenRange())
}

func TestAddDurations(t *testing.T) {
	d1 := NewDuration(0, 1)
	d2 := NewDuration(2, 50)
	w := NewRecord(Ɀ_Date_(2020, 1, 1))
	w.AddDuration(d1, NewSummary("Duration 1"))
	w.AddDuration(d2, NewSummary("Duration 2"))
	require.Len(t, w.Entries(), 2)
	assert.Equal(t, d1, w.Entries()[0].value)
	assert.Equal(t, NewSummary("Duration 1"), w.Entries()[0].Summary())
	assert.Equal(t, d2, w.Entries()[1].value)
	assert.Equal(t, NewSummary("Duration 2"), w.Entries()[1].Summary())
}
