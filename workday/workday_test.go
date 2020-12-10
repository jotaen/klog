package workday

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	"testing"
)

func TestSavesDateUponCreation(t *testing.T) {
	date, _ := datetime.CreateDate(2020, 1, 1)
	workDay := Create(date)
	assert.Equal(t, workDay.Date(), date)
}

func TestAddRanges(t *testing.T) {
	date, _ := datetime.CreateDate(2020, 1, 1)
	time1, _ := datetime.CreateTime(9, 7)
	time2, _ := datetime.CreateTime(12, 59)
	time3, _ := datetime.CreateTime(13, 49)
	time4, _ := datetime.CreateTime(17, 12)
	range1, _ := datetime.CreateTimeRange(time1, time2)
	range2, _ := datetime.CreateTimeRange(time3, time4)
	w := Create(date)
	w.AddRange(range1)
	w.AddRange(range2)
	assert.Equal(t, []datetime.TimeRange{range1, range2}, w.Ranges())
}

func TestAddOpenRange(t *testing.T) {
	date, _ := datetime.CreateDate(2020, 1, 1)
	time1, _ := datetime.CreateTime(9, 7)
	range1, _ := datetime.CreateTimeRange(time1, nil)
	w := Create(date)
	assert.Equal(t, nil, w.OpenRange())
	w.AddRange(range1)
	assert.Equal(t, range1, w.OpenRange())
}

func TestOkayWhenAddingValidDuration(t *testing.T) {
	date, _ := datetime.CreateDate(2020, 1, 1)
	w := Create(date)
	err := w.AddDuration(datetime.Duration(1))
	assert.Nil(t, err)
	assert.Equal(t, []datetime.Duration{datetime.Duration(1)}, w.Times())
}
