package datetime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerialiseDuration(t *testing.T) {
	assert.Equal(t, "0:01", Duration(1).ToString())
	assert.Equal(t, "2:20", Duration(140).ToString())
	assert.Equal(t, "15:00", Duration(900).ToString())
	assert.Equal(t, "68:59", Duration(4139).ToString())
}
