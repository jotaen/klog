package workday

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	. "klog/testutil/datetime"
	"testing"
)

func TestSavesDateUponCreation(t *testing.T) {
	date := Date_(2020, 1, 1)
	workDay := NewWorkDay(date)
	assert.Equal(t, workDay.Date(), date)
}

func TestAddRanges(t *testing.T) {
	range1 := Range_(Time_(9, 7), Time_(12, 59))
	range2 := Range_(Time_(13, 49), Time_(17, 12))
	w := NewWorkDay(Date_(2020, 1, 1))
	w.AddRange(range1)
	w.AddRange(range2)
	assert.Equal(t, []datetime.TimeRange{range1, range2}, w.Ranges())
}

func TestStartOpenRange(t *testing.T) {
	time := Time_(11, 23)
	w := NewWorkDay(Date_(2020, 1, 1))
	assert.Equal(t, nil, w.OpenRange())
	w.StartOpenRange(time)
	assert.Equal(t, time, w.OpenRange())
}

func TestCloseOpenRange(t *testing.T) {
	start := Time_(19, 22)
	w := NewWorkDay(Date_(2012, 6, 17))
	w.StartOpenRange(start)
	end := Time_(20, 55)
	w.EndOpenRange(end)
	assert.Nil(t, w.OpenRange())
	assert.Equal(t, []datetime.TimeRange{Range_(start, end)}, w.Ranges())
}

func TestOkayWhenAddingValidDuration(t *testing.T) {
	w := NewWorkDay(Date_(2020, 1, 1))
	w.AddDuration(datetime.Duration(1))
	assert.Equal(t, []datetime.Duration{datetime.Duration(1)}, w.Times())
}
