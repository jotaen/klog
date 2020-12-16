package datetime

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

func TestDetectsInvalidDates(t *testing.T) {
	invalidMonth, err := NewDate(2005, 13, 15)
	assert.Nil(t, invalidMonth)
	assert.Error(t, err)

	invalidDay, err := NewDate(2005, 2, 30)
	assert.Nil(t, invalidDay)
	assert.Error(t, err)
}

func TestSerialiseDate(t *testing.T) {
	d, _ := NewDate(2005, 12, 31)
	assert.Equal(t, "2005-12-31", d.ToString())
}

func TestSerialiseDatePadsLeadingZeros(t *testing.T) {
	d, _ := NewDate(2005, 3, 9)
	assert.Equal(t, "2005-03-09", d.ToString())
}
