package datetime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerialiseDate(t *testing.T) {
	d := Date{Year: 2005, Month: 3, Day: 9}
	assert.Equal(t, d.ToString(), "2005-03-09")
}

func TestSerialiseTime(t *testing.T) {
	tm := Time{Hour: 8, Minute: 5}
	assert.Equal(t, tm.ToString(), "08:05")
}
