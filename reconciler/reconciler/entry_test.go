package reconciler

import (
	"testing"
	"cloud.google.com/go/civil"
	"github.com/stretchr/testify/assert"
)

func TestSumUpTimes(t *testing.T) {
	day := Entry{
		Times: []Minutes { Minutes(60), Minutes(120) },
	}
	assert.Equal(t, day.TotalTime(), Minutes(180))
}

func TestSumUpRanges(t *testing.T) {
	x1_start, _ := civil.ParseTime("9:00:00")
	x1_end, _ := civil.ParseTime("13:00:00")
	x2_start, _ := civil.ParseTime("14:00:00")
	x2_end, _ := civil.ParseTime("17:00:00")
	day := Entry{
		Ranges: []Range {
			Range{ Start: x1_start, End: x1_end },
			Range{ Start: x2_start, End: x2_end },
		},
	}
	assert.Equal(t, day.TotalTime(), Minutes(240 + 180))
}
