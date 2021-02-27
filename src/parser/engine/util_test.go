package engine

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGroupLines(t *testing.T) {
	blocks := GroupIntoBlocks([]Line{
		{Value: []rune("a1")},
		{Value: []rune("a2")},
		{Value: []rune("")},
		{Value: []rune("b1")},
		{Value: []rune("    ")},
		{Value: []rune("\t")},
		{Value: []rune("c1")},
	})
	require.Len(t, blocks, 3)

	require.Len(t, blocks[0], 2)
	assert.Equal(t, blocks[0][0].Value, []rune("a1"))
	assert.Equal(t, blocks[0][1].Value, []rune("a2"))

	require.Len(t, blocks[1], 1)
	assert.Equal(t, blocks[1][0].Value, []rune("b1"))

	require.Len(t, blocks[2], 1)
	assert.Equal(t, blocks[2][0].Value, []rune("c1"))
}
