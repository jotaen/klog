package workday

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOkayWhenValidDateIsPresent(t *testing.T) {
	date := Date{Year: 2020, Month: time.January, Day: 1}
	workDay, err := Create(date)
	assert.Nil(t, err)
	assert.Equal(t, workDay.Date(), date)
}

func TestErrorWhenDateIsInvalid(t *testing.T) {
	workDay, err := Create(Date{Year: 2020, Month: time.January, Day: 99})
	assert.Nil(t, workDay)
	assert.Equal(t, err.(*WorkDayError).Code, INVALID_DATE)
}
