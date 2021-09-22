package service

import (
	. "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashYieldsDistinctValues(t *testing.T) {
	dayHashes := make(map[DayHash]bool)
	weekHashes := make(map[WeekHash]bool)
	monthHashes := make(map[MonthHash]bool)
	quarterHashes := make(map[QuarterHash]bool)
	yearHashes := make(map[YearHash]bool)

	// 1.1.1000 is a Wednesday. 1000 days later it’s Sunday, 27.9.1002
	initialDate := Ɀ_Date_(1000, 1, 1)
	for i := 0; i < 1000; i++ {
		d := initialDate.PlusDays(i)
		dayHashes[NewDayHash(d)] = true
		weekHashes[NewWeekHash(d)] = true
		monthHashes[NewMonthHash(d)] = true
		quarterHashes[NewQuarterHash(d)] = true
		yearHashes[NewYearHash(d)] = true
	}

	assert.Len(t, dayHashes, 1000)
	assert.Len(t, weekHashes, 145)
	assert.Len(t, monthHashes, 33)
	assert.Len(t, quarterHashes, 11)
	assert.Len(t, yearHashes, 3)
}
