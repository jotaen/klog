package datetime

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateTimeRange(t *testing.T) {
	time1, _ := CreateTime(11, 25)
	time2, _ := CreateTime(17, 10)
	tr, err := CreateTimeRange(time1, time2)
	require.Nil(t, err)
	require.NotNil(t, tr)
	assert.Equal(t, tr.Start(), time1)
	assert.Equal(t, tr.End(), time2)
}

func TestCreationFailsIfStartIsBeforeEnd(t *testing.T) {
	time1, _ := CreateTime(14, 00)
	time2, _ := CreateTime(13, 59)
	tr, err := CreateTimeRange(time1, time2)
	assert.Nil(t, tr)
	assert.Equal(t, errors.New("ILLEGAL_RANGE"), err)
}
