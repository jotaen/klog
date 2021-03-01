package parsing

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGroupLines(t *testing.T) {
	blocks := GroupIntoBlocks([]Line{
		{Text: "a1"},
		{Text: "a2"},
		{Text: ""},
		{Text: "b1"},
		{Text: "    "},
		{Text: "\t"},
		{Text: "c1"},
	})
	require.Len(t, blocks, 3)

	require.Len(t, blocks[0], 2)
	assert.Equal(t, blocks[0][0].Text, "a1")
	assert.Equal(t, blocks[0][1].Text, "a2")

	require.Len(t, blocks[1], 1)
	assert.Equal(t, blocks[1][0].Text, "b1")

	require.Len(t, blocks[2], 1)
	assert.Equal(t, blocks[2][0].Text, "c1")
}
