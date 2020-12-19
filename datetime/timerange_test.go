package datetime

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	// "klog/testutil" REMINDER: can not use `testutil` because of circular import
	"testing"
)

func TestCreateTimeRange(t *testing.T) {
	time1, _ := NewTime(11, 25)
	time2, _ := NewTime(17, 10)
	tr, err := NewTimeRange(time1, time2)
	require.Nil(t, err)
	require.NotNil(t, tr)
	assert.Equal(t, time1, tr.Start())
	assert.Equal(t, time2, tr.End())
}

func TestCreateOverlappingTimeRangeYesterday(t *testing.T) {
	time1, _ := NewTime(23, 30)
	time2, _ := NewTime(8, 10)
	tr, err := NewOverlappingTimeRange(time1, true, time2, false)
	require.Nil(t, err)
	require.NotNil(t, tr)
	assert.Equal(t, time1, tr.Start())
	assert.Equal(t, time2, tr.End())
	assert.Equal(t, true, tr.IsStartYesterday())
	assert.Equal(t, false, tr.IsEndTomorrow())
	assert.Equal(t, NewDuration(8, 40), tr.Duration())
}

func TestCreateOverlappingTimeRangeTomorrow(t *testing.T) {
	time1, _ := NewTime(18, 15)
	time2, _ := NewTime(1, 45)
	tr, err := NewOverlappingTimeRange(time1, false, time2, true)
	require.Nil(t, err)
	require.NotNil(t, tr)
	assert.Equal(t, time1, tr.Start())
	assert.Equal(t, time2, tr.End())
	assert.Equal(t, false, tr.IsStartYesterday())
	assert.Equal(t, true, tr.IsEndTomorrow())
	assert.Equal(t, NewDuration(7, 30), tr.Duration())
}

func TestCreationFailsIfStartIsBeforeEnd(t *testing.T) {
	time1, _ := NewTime(14, 00)
	time2, _ := NewTime(13, 59)
	tr, err := NewTimeRange(time1, time2)
	assert.Nil(t, tr)
	assert.EqualError(t, err, "ILLEGAL_RANGE")
}
