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

func TestSerialiseTimeWithoutLeadingZeros(t *testing.T) {
	tm, _ := CreateTime(8, 5)
	assert.Equal(t, "8:05", tm.ToString())
}
