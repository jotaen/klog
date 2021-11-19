package commands

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPrintOutEmptyInput(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(``)._Run((&Print{}).Run)
	require.Nil(t, err)
	assert.Equal(t, "", state.printBuffer)
}

func TestPrintOutRecord(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
2018-01-31
Hello #world
	1h
`)._Run((&Print{}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-31
Hello #world
    1h

`, state.printBuffer)
}
