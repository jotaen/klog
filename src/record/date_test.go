package record

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

func TestParseDate(t *testing.T) {
	d, err := NewDateFromString("1856-10-22")
	assert.Nil(t, err)
	should, _ := NewDate(1856, 10, 22)
	assert.Equal(t, d, should)
}

func TestParseDateFailsIfMalformed(t *testing.T) {
	for _, s := range []string{
		"1856-1-2",
		"20-12-12",
		"asdf",
		"01.01.2000",
	} {
		d, err := NewDateFromString(s)
		assert.Nil(t, d)
		assert.EqualError(t, err, "MALFORMED_DATE")
	}
}
