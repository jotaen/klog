package record

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInitialiseRecord(t *testing.T) {
	date := Ɀ_Date_(2020, 1, 1)
	r := NewRecord(date)
	assert.Equal(t, r.Date(), date)
	assert.Equal(t, nil, r.ShouldTotal())
	assert.Equal(t, "", r.Summary())
	assert.Len(t, r.Entries(), 0)
}

func TestSavesSummary(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	err := r.SetSummary("Hello World")
	require.Nil(t, err)
	assert.Equal(t, "Hello World", r.Summary())
}

func TestSummaryCannotContainWhitespaceAtBeginningOfLine(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	require.Error(t, r.SetSummary("Hello\n World"))
	require.Error(t, r.SetSummary(" Hello"))
	assert.Equal(t, "", r.Summary()) // Still empty
}

func TestSavesShouldTotal(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.SetShouldTotal(NewDuration(5, 0))
	assert.Equal(t, NewDuration(5, 0), r.ShouldTotal())
}

func TestAddRanges(t *testing.T) {
	range1 := Ɀ_Range_(Ɀ_Time_(9, 7), Ɀ_Time_(12, 59))
	range2 := Ɀ_Range_(Ɀ_Time_(13, 49), Ɀ_Time_(17, 12))
	w := NewRecord(Ɀ_Date_(2020, 1, 1))
	w.AddRange(range1, "Range 1")
	w.AddRange(range2, "Range 2")
	require.Len(t, w.Entries(), 2)
	assert.Equal(t, range1, w.Entries()[0].Value())
	assert.Equal(t, "Range 1", w.Entries()[0].SummaryAsString())
	assert.Equal(t, range2, w.Entries()[1].Value())
	assert.Equal(t, "Range 2", w.Entries()[1].SummaryAsString())
}

func TestStartOpenRange(t *testing.T) {
	time := Ɀ_Time_(11, 23)
	w := NewRecord(Ɀ_Date_(2020, 1, 1))
	assert.Equal(t, nil, w.OpenRange())
	_ = w.StartOpenRange(time, "Open Range")
	require.Len(t, w.Entries(), 1)
	assert.Equal(t, time, w.Entries()[0].Value())
	assert.Equal(t, "Open Range", w.Entries()[0].SummaryAsString())
}

func TestCannotStartSecondOpenRange(t *testing.T) {
	time := Ɀ_Time_(11, 23)
	w := NewRecord(Ɀ_Date_(2020, 1, 1))
	assert.Equal(t, nil, w.OpenRange())
	_ = w.StartOpenRange(time, "Open Range")
	err := w.StartOpenRange(time, "Open Range")
	require.Error(t, err)
}

func TestCloseOpenRange(t *testing.T) {
	start := Ɀ_Time_(19, 22)
	w := NewRecord(Ɀ_Date_(2012, 6, 17))
	_ = w.StartOpenRange(start, "Started")
	end := Ɀ_Time_(20, 55)
	err := w.EndOpenRange(end)
	require.Nil(t, err)
	assert.Nil(t, w.OpenRange())
	require.Len(t, w.Entries(), 1)
	assert.Equal(t, Ɀ_Range_(start, end), w.Entries()[0].Value())
	assert.Equal(t, "Started", w.Entries()[0].SummaryAsString())
}

func TestCloseOpenRangeFailsIfResultingRangeIsInvalid(t *testing.T) {
	start := Ɀ_Time_(19, 22)
	w := NewRecord(Ɀ_Date_(2012, 6, 17))
	_ = w.StartOpenRange(start, "Started")
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
	w.AddDuration(d1, "Duration 1")
	w.AddDuration(d2, "Duration 2")
	require.Len(t, w.Entries(), 2)
	assert.Equal(t, d1, w.Entries()[0].Value())
	assert.Equal(t, "Duration 1", w.Entries()[0].SummaryAsString())
	assert.Equal(t, d2, w.Entries()[1].Value())
	assert.Equal(t, "Duration 2", w.Entries()[1].SummaryAsString())
}
