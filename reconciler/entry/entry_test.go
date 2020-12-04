package entry

import (
	"cloud.google.com/go/civil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSumUpTimes(t *testing.T) {
	day := Entry{
		Times: []Minutes{Minutes(60), Minutes(120)},
	}
	assert.Equal(t, day.TotalTime(), Minutes(180))
}

func TestSumUpRanges(t *testing.T) {
	day := Entry{
		Ranges: []Range{
			Range{Start: civil.Time{Hour: 9, Minute: 07}, End: civil.Time{Hour: 12, Minute: 59}},
			Range{Start: civil.Time{Hour: 13, Minute: 49}, End: civil.Time{Hour: 17, Minute: 12}},
		},
	}
	assert.Equal(t, day.TotalTime(), Minutes(435))
}

func TestSumUpTimesAndRanges(t *testing.T) {
	day := Entry{
		Times: []Minutes{Minutes(90)},
		Ranges: []Range{
			Range{Start: civil.Time{Hour: 8, Minute: 00}, End: civil.Time{Hour: 12, Minute: 00}},
		},
	}
	assert.Equal(t, day.TotalTime(), Minutes(330))
}
