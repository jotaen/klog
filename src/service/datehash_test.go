package service

import (
	"github.com/stretchr/testify/assert"
	. "klog"
	"testing"
)

func TestHashYieldsDistinctValues(t *testing.T) {
	hashes := make(map[DayHash]bool)
	for i, d := 0, Ɀ_Date_(1000, 1, 1); i < 1000; i++ {
		d = d.PlusDays(i)
		hashes[NewDayHash(d)] = true
	}
	assert.Len(t, hashes, 1000)
}
