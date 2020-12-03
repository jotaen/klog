package reconciler

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSumUpHours(t *testing.T) {
	day := Entry{
		Times: []int64 { int64(60), int64(120) },
	}
	assert.Equal(t, day.TotalTime(), int64(180))
}
