package workday

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOkayWhenAddingValidTimes(t *testing.T) {
	date := Date{Year: 2020, Month: time.January, Day: 1}
	w, _ := Create(date)
	err := w.AddTime(Minutes(1))
	assert.Nil(t, err)
	assert.Equal(t, w.Times(), []Minutes{Minutes(1)})
}

func TestErrorWhenAddingInvalidTimes(t *testing.T) {
	date := Date{Year: 2020, Month: time.January, Day: 1}
	w, _ := Create(date)
	err := w.AddTime(Minutes(-1))
	assert.Equal(t, err.(*WorkDayError).Code, NEGATIVE_TIME)
	assert.Equal(t, len(w.Times()), 0)
}
