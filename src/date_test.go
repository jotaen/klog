package klog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRecognisesValidDate(t *testing.T) {
	d, err := NewDate(2005, 4, 15)
	assert.Nil(t, err)
	assert.Equal(t, 2005, d.Year())
	assert.Equal(t, 4, d.Month())
	assert.Equal(t, 15, d.Day())
}

func TestReconWithDate(t *testing.T) {
	d, _ := NewDate(2005, 12, 31)
	assert.Equal(t, Ɀ_Date_(2006, 1, 1), d.PlusDays(1))
	assert.Equal(t, Ɀ_Date_(2006, 2, 1), d.PlusDays(32))
	assert.Equal(t, Ɀ_Date_(2005, 12, 30), d.PlusDays(-1))
}

func TestHashYieldsDistinctValues(t *testing.T) {
	hashes := make(map[DateHash]bool)
	for i, d := 0, Ɀ_Date_(1000, 1, 1); i < 1000; i++ {
		d = d.PlusDays(i)
		hashes[d.Hash()] = true
	}
	assert.Len(t, hashes, 1000)
}

func TestDetectsUnrepresentableDates(t *testing.T) {
	invalidMonth, err := NewDate(2005, 13, 15)
	assert.Nil(t, invalidMonth)
	assert.EqualError(t, err, "UNREPRESENTABLE_DATE")

	invalidDay, err := NewDate(2005, 2, 30)
	assert.Nil(t, invalidDay)
	assert.EqualError(t, err, "UNREPRESENTABLE_DATE")
}

func TestSerialiseDate(t *testing.T) {
	d := Ɀ_Date_(2005, 12, 31)
	assert.Equal(t, "2005-12-31", d.ToString())
}

func TestSerialiseDatePadsLeadingZeros(t *testing.T) {
	d := Ɀ_Date_(2005, 3, 9)
	assert.Equal(t, "2005-03-09", d.ToString())
}

func TestParseDateWithDashes(t *testing.T) {
	d, err := NewDateFromString("1856-10-22")
	assert.Nil(t, err)
	should, _ := NewDate(1856, 10, 22)
	assert.Equal(t, d, should)
}

func TestEquality(t *testing.T) {
	a := Ɀ_Date_(2005, 1, 1)
	b := Ɀ_Date_(2005, 1, 1)
	c := Ɀ_Date_(1982, 12, 31)
	assert.True(t, a.IsEqualTo(b))
	assert.False(t, a.IsEqualTo(c))
	assert.False(t, b.IsEqualTo(c))
}

func TestComparison(t *testing.T) {
	a := Ɀ_Date_(2005, 3, 15)
	b := Ɀ_Date_(2005, 3, 15)
	c := Ɀ_Date_(2005, 3, 16)
	d := Ɀ_Date_(2004, 3, 16)
	e := Ɀ_Date_(2005, 4, 1)
	assert.True(t, b.IsAfterOrEqual(a))
	assert.True(t, c.IsAfterOrEqual(a))
	assert.True(t, a.IsAfterOrEqual(d))
	assert.True(t, e.IsAfterOrEqual(c))
}

func TestParseDateWithSlashes(t *testing.T) {
	original := "1856/10/22"
	d, err := NewDateFromString(original)
	assert.Nil(t, err)
	should, _ := NewDate(1856, 10, 22)
	assert.True(t, should.IsEqualTo(d))
	assert.Equal(t, original, d.ToString())
}

func TestParseDateFailsIfMalformed(t *testing.T) {
	for _, s := range []string{
		"1856-1-2",
		"1856/01-02",
		"20-12-12",
		"asdf",
		"01.01.2000",
	} {
		d, err := NewDateFromString(s)
		assert.Nil(t, d)
		assert.EqualError(t, err, "MALFORMED_DATE")
	}
}

func TestCalculateWeekday(t *testing.T) {
	for _, d := range []struct {
		d Date
		w int
	}{
		{Ɀ_Date_(2021, 01, 15), 5},
		{Ɀ_Date_(2021, 01, 16), 6},
		{Ɀ_Date_(2021, 01, 17), 7}, // Sunday
		{Ɀ_Date_(2021, 01, 18), 1},
		{Ɀ_Date_(2021, 01, 19), 2},
		{Ɀ_Date_(2021, 01, 20), 3},
		{Ɀ_Date_(2021, 01, 21), 4},
		{Ɀ_Date_(2021, 01, 22), 5},
	} {
		assert.Equal(t, d.w, d.d.Weekday())
	}
}
