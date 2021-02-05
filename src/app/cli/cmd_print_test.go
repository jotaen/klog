package cli

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPrintOutEmptyInput(t *testing.T) {
	out, err := RunWithContext(``, (&Print{}).Run)
	require.Nil(t, err)
	assert.Equal(t, "\n", out)
}

func TestPrintOutRecord(t *testing.T) {
	out, err := RunWithContext(`
2018-01-31
Hello #world
	1h
`, (&Print{}).Run)
	require.Nil(t, err)
	assert.Equal(t, "\n2018-01-31\nHello #world\n    1h\n", out)
}
