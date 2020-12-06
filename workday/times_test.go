package workday

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	"testing"
)

func TestOkayWhenAddingValidTimes(t *testing.T) {
	date := datetime.Date{Year: 2020, Month: 1, Day: 1}
	w, _ := Create(date)
	err := w.AddTime(datetime.Duration(1))
	assert.Nil(t, err)
	assert.Equal(t, w.Times(), []datetime.Duration{datetime.Duration(1)})
}

func TestErrorWhenAddingInvalidTimes(t *testing.T) {
	date := datetime.Date{Year: 2020, Month: 1, Day: 1}
	w, _ := Create(date)
	err := w.AddTime(datetime.Duration(-1))
	assert.Equal(t, err.(*WorkDayError).Code, NEGATIVE_TIME)
	assert.Equal(t, len(w.Times()), 0)
}
