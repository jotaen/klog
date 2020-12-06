package workday

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	"testing"
)

func TestOkayWhenValidDateIsPresent(t *testing.T) {
	date := datetime.Date{Year: 2020, Month: 1, Day: 1}
	workDay, err := Create(date)
	assert.Nil(t, err)
	assert.Equal(t, workDay.Date(), date)
}

func TestErrorWhenDateIsInvalid(t *testing.T) {
	workDay, err := Create(datetime.Date{Year: 2020, Month: 1, Day: 99})
	assert.Nil(t, workDay)
	assert.Equal(t, err.(*WorkDayError).Code, INVALID_DATE)
}
