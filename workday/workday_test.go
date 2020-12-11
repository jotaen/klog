package workday

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	datetime2 "klog/testutil/datetime"
	"testing"
)

func TestSavesDateUponCreation(t *testing.T) {
	date := datetime2.Date_(2020, 1, 1)
	workDay := Create(date)
	assert.Equal(t, workDay.Date(), date)
}

func TestAddRanges(t *testing.T) {
	range1 := datetime2.Range_(datetime2.Time_(9, 7), datetime2.Time_(12, 59))
	range2 := datetime2.Range_(datetime2.Time_(13, 49), datetime2.Time_(17, 12))
	w := Create(datetime2.Date_(2020, 1, 1))
	w.AddRange(range1)
	w.AddRange(range2)
	assert.Equal(t, []datetime.TimeRange{range1, range2}, w.Ranges())
}

func TestAddOpenRange(t *testing.T) {
	time := datetime2.Time_(11, 23)
	w := Create(datetime2.Date_(2020, 1, 1))
	assert.Equal(t, nil, w.OpenRangeStart())
	w.SetOpenRangeStart(time)
	assert.Equal(t, time, w.OpenRangeStart())
}

func TestOkayWhenAddingValidDuration(t *testing.T) {
	w := Create(datetime2.Date_(2020, 1, 1))
	w.AddDuration(datetime.Duration(1))
	assert.Equal(t, []datetime.Duration{datetime.Duration(1)}, w.Times())
}
