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
	assert.Nil(t, r.Summary())
	assert.Len(t, r.Entries(), 0)
}

func TestSavesSummary(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.SetSummary(Ɀ_RecordSummary_("Hello World"))
	assert.Equal(t, Ɀ_RecordSummary_("Hello World"), r.Summary())

	r.SetSummary(Ɀ_RecordSummary_("Two", "Lines"))
	assert.Equal(t, Ɀ_RecordSummary_("Two", "Lines"), r.Summary())

	r.SetSummary(nil)
	assert.Nil(t, r.Summary())
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
	range3 := Ɀ_Range_(Ɀ_Time_(23, 3), Ɀ_Time_(23, 3))
	w := NewRecord(Ɀ_Date_(2020, 1, 1))
	w.AddRange(range1, Ɀ_EntrySummary_("Range 1"))
	w.AddRange(range2, Ɀ_EntrySummary_("Range 2", "With second line"))
	w.AddRange(range3, nil)
	require.Len(t, w.Entries(), 3)

	assert.Equal(t, range1, w.Entries()[0].value)
	assert.Equal(t, Ɀ_EntrySummary_("Range 1"), w.Entries()[0].Summary())

	assert.Equal(t, range2, w.Entries()[1].value)
	assert.Equal(t, Ɀ_EntrySummary_("Range 2", "With second line"), w.Entries()[1].Summary())

	assert.Equal(t, range3, w.Entries()[2].value)
	assert.Nil(t, w.Entries()[2].Summary())
}

func TestStartOpenRange(t *testing.T) {
	time := Ɀ_Time_(11, 23)
	w := NewRecord(Ɀ_Date_(2020, 1, 1))
	assert.Equal(t, nil, w.OpenRange())
	_ = w.StartOpenRange(time, Ɀ_EntrySummary_("Open Range"))
	require.Len(t, w.Entries(), 1)
	assert.Equal(t, NewOpenRange(time), w.Entries()[0].value)
	assert.Equal(t, Ɀ_EntrySummary_("Open Range"), w.Entries()[0].Summary())
}

func TestCannotStartSecondOpenRange(t *testing.T) {
	time := Ɀ_Time_(11, 23)
	w := NewRecord(Ɀ_Date_(2020, 1, 1))
	assert.Equal(t, nil, w.OpenRange())
	_ = w.StartOpenRange(time, Ɀ_EntrySummary_("Open Range"))
	err := w.StartOpenRange(time, Ɀ_EntrySummary_("Open Range"))
	require.Error(t, err)
}

func TestCloseOpenRange(t *testing.T) {
	start := Ɀ_Time_(19, 22)
	w := NewRecord(Ɀ_Date_(2012, 6, 17))
	_ = w.StartOpenRange(start, Ɀ_EntrySummary_("Started"))
	end := Ɀ_Time_(20, 55)
	err := w.EndOpenRange(end)
	require.Nil(t, err)
	assert.Nil(t, w.OpenRange())
	require.Len(t, w.Entries(), 1)
	assert.Equal(t, Ɀ_Range_(start, end), w.Entries()[0].value)
	assert.Equal(t, Ɀ_EntrySummary_("Started"), w.Entries()[0].Summary())
}

func TestCloseOpenRangeFailsIfResultingRangeIsInvalid(t *testing.T) {
	start := Ɀ_Time_(19, 22)
	w := NewRecord(Ɀ_Date_(2012, 6, 17))
	_ = w.StartOpenRange(start, Ɀ_EntrySummary_("Started"))
	oldEntry := w.OpenRange()
	end := Ɀ_Time_(1, 30)
	err := w.EndOpenRange(end)
	require.Error(t, err)
	assert.Equal(t, oldEntry, w.OpenRange())
}

func TestAddDurations(t *testing.T) {
	d1 := NewDuration(0, 1)
	d2 := NewDuration(2, 50)
	d3 := NewDuration(1, 0)
	w := NewRecord(Ɀ_Date_(2020, 1, 1))
	w.AddDuration(d1, Ɀ_EntrySummary_("Duration 1"))
	w.AddDuration(d2, Ɀ_EntrySummary_("Duration 2", "With second line"))
	w.AddDuration(d3, nil)
	require.Len(t, w.Entries(), 3)

	assert.Equal(t, d1, w.Entries()[0].value)
	assert.Equal(t, Ɀ_EntrySummary_("Duration 1"), w.Entries()[0].Summary())

	assert.Equal(t, d2, w.Entries()[1].value)
	assert.Equal(t, Ɀ_EntrySummary_("Duration 2", "With second line"), w.Entries()[1].Summary())

	assert.Equal(t, d3, w.Entries()[2].value)
	assert.Nil(t, w.Entries()[2].Summary())
}
