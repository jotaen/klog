package reconciler

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSumUpHours(t *testing.T) {
	day := Entry{
		Times: []Minutes { Minutes(60), Minutes(120) },
	}
	assert.Equal(t, day.TotalTime(), Minutes(180))
}
