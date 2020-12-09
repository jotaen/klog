package workday

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	"testing"
)

func TestOkayWhenAddingValidDuration(t *testing.T) {
	date, _ := datetime.CreateDate(2020, 1, 1)
	w := Create(date)
	err := w.AddDuration(datetime.Duration(1))
	assert.Nil(t, err)
	assert.Equal(t, w.Times(), []datetime.Duration{datetime.Duration(1)})
}

func TestErrorWhenAddingInvalidDuration(t *testing.T) {
	date, _ := datetime.CreateDate(2020, 1, 1)
	w := Create(date)
	err := w.AddDuration(datetime.Duration(-1))
	assert.Equal(t, err.Error(), "NEGATIVE_DURATION")
	assert.Equal(t, len(w.Times()), 0)
}
