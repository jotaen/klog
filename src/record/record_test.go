package record

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSavesDateUponCreation(t *testing.T) {
	date := Ɀ_Date_(2020, 1, 1)
	r := NewRecord(date)
	assert.Equal(t, r.Date(), date)
}

func TestAddRanges(t *testing.T) {
	range1 := Ɀ_Range_(Ɀ_Time_(9, 7), Ɀ_Time_(12, 59))
	range2 := Ɀ_Range_(Ɀ_Time_(13, 49), Ɀ_Time_(17, 12))
	w := NewRecord(Ɀ_Date_(2020, 1, 1))
	w.AddRange(range1)
	w.AddRange(range2)
	require.Len(t, w.Entries(), 2)
	assert.Equal(t, range1, w.Entries()[0].Value())
	assert.Equal(t, range2, w.Entries()[1].Value())
}

func TestStartOpenRange(t *testing.T) {
	time := Ɀ_Time_(11, 23)
	w := NewRecord(Ɀ_Date_(2020, 1, 1))
	assert.Equal(t, nil, w.OpenRange())
	w.StartOpenRange(time)
	require.Len(t, w.Entries(), 1)
	assert.Equal(t, time, w.Entries()[0].Value())
}

func TestCloseOpenRange(t *testing.T) {
	start := Ɀ_Time_(19, 22)
	w := NewRecord(Ɀ_Date_(2012, 6, 17))
	w.StartOpenRange(start)
	end := Ɀ_Time_(20, 55)
	w.EndOpenRange(end)
	assert.Nil(t, w.OpenRange())
	require.Len(t, w.Entries(), 1)
	assert.Equal(t, Ɀ_Range_(start, end), w.Entries()[0].Value())
}

func TestAddDurations(t *testing.T) {
	d1 := NewDuration(0, 1)
	d2 := NewDuration(2, 50)
	w := NewRecord(Ɀ_Date_(2020, 1, 1))
	w.AddDuration(d1)
	w.AddDuration(d2)
	require.Len(t, w.Entries(), 2)
	assert.Equal(t, d1, w.Entries()[0].Value())
	assert.Equal(t, d2, w.Entries()[1].Value())
}
