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

func TestCreateVoidTimeRange(t *testing.T) {
	time1, _ := NewTime(12, 00)
	time2, _ := NewTime(12, 00)
	tr, err := NewTimeRange(time1, time2)
	require.Nil(t, err)
	require.NotNil(t, tr)
	assert.Equal(t, time1, tr.Start())
	assert.Equal(t, time2, tr.End())
	assert.Equal(t, NewDuration(0, 00), tr.Duration())
}

func TestCreateOverlappingTimeRangeStartingYesterday(t *testing.T) {
	time1, _ := NewTimeYesterday(23, 30)
	time2, _ := NewTime(8, 10)
	tr, err := NewTimeRange(time1, time2)
	require.Nil(t, err)
	require.NotNil(t, tr)
	assert.Equal(t, time1, tr.Start())
	assert.Equal(t, time2, tr.End())
	assert.Equal(t, NewDuration(8, 40), tr.Duration())
}

func TestCreateOverlappingTimeRangeEndingTomorrow(t *testing.T) {
	time1, _ := NewTime(18, 15)
	time2, _ := NewTimeTomorrow(1, 45)
	tr, err := NewTimeRange(time1, time2)
	require.Nil(t, err)
	require.NotNil(t, tr)
	assert.Equal(t, time1, tr.Start())
	assert.Equal(t, time2, tr.End())
	assert.Equal(t, NewDuration(7, 30), tr.Duration())
}

func TestCreationFailsIfStartIsBeforeEnd(t *testing.T) {
	for _, p := range []func() (TimeRange, error){
		func() (TimeRange, error) {
			start, _ := NewTime(15, 00)
			end, _ := NewTime(14, 00)
			return NewTimeRange(start, end)
		},
		func() (TimeRange, error) {
			start, _ := NewTime(14, 00)
			end, _ := NewTimeYesterday(15, 00)
			return NewTimeRange(start, end)
		},
		func() (TimeRange, error) {
			start, _ := NewTimeTomorrow(14, 00)
			end, _ := NewTime(15, 00)
			return NewTimeRange(start, end)
		},
	} {
		tr, err := p()
		assert.Nil(t, tr)
		assert.EqualError(t, err, "ILLEGAL_RANGE")
	}
}
