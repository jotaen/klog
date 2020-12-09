package datetime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOnlyConstructsValidTimes(t *testing.T) {
	tm, err := CreateTime(22, 13)
	assert.Equal(t, tm.Hour(), 22)
	assert.Equal(t, tm.Minute(), 13)
	assert.Nil(t, err)
}

func TestDetectsInvalidTimes(t *testing.T) {
	invalidHour, err := CreateTime(25, 30)
	assert.Error(t, err)
	assert.Nil(t, invalidHour)

	invalidMinute, err := CreateTime(4, 85)
	assert.Error(t, err)
	assert.Nil(t, invalidMinute)
}

func TestSerialiseTime(t *testing.T) {
	tm, _ := CreateTime(13, 45)
	assert.Equal(t, "13:45", tm.ToString())
}

func TestSerialiseTimePadsLeadingZeros(t *testing.T) {
	tm, _ := CreateTime(8, 5)
	assert.Equal(t, "08:05", tm.ToString())
}

func TestSerialiseDuration(t *testing.T) {
	assert.Equal(t, "00:01", Duration(1).ToString())
	assert.Equal(t, "02:20", Duration(140).ToString())
	assert.Equal(t, "15:00", Duration(900).ToString())
	assert.Equal(t, "68:59", Duration(4139).ToString())
}
