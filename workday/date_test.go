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
