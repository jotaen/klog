package workday

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	"klog/testutil"
	"testing"
)

func TestSavesDateUponCreation(t *testing.T) {
	date := testutil.Date_(2020, 1, 1)
	workDay := Create(date)
	assert.Equal(t, workDay.Date(), date)
}

func TestAddRanges(t *testing.T) {
	range1 := testutil.Range_(testutil.Time_(9, 7), testutil.Time_(12, 59))
	range2 := testutil.Range_(testutil.Time_(13, 49), testutil.Time_(17, 12))
	w := Create(testutil.Date_(2020, 1, 1))
	w.AddRange(range1)
	w.AddRange(range2)
	assert.Equal(t, []datetime.TimeRange{range1, range2}, w.Ranges())
}

func TestAddOpenRange(t *testing.T) {
	range1 := testutil.Range_(testutil.Time_(9, 7), nil)
	w := Create(testutil.Date_(2020, 1, 1))
	assert.Equal(t, nil, w.OpenRange())
	w.AddRange(range1)
	assert.Equal(t, range1, w.OpenRange())
}

func TestOkayWhenAddingValidDuration(t *testing.T) {
	w := Create(testutil.Date_(2020, 1, 1))
	err := w.AddDuration(datetime.Duration(1))
	assert.Nil(t, err)
	assert.Equal(t, []datetime.Duration{datetime.Duration(1)}, w.Times())
}
