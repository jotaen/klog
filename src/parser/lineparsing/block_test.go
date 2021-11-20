package lineparsing

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGroupLines(t *testing.T) {
	blocks := GroupIntoBlocks([]Line{
		{Text: ""},
		{Text: "  "},
		{Text: "a1"},
		{Text: "a2"},
		{Text: ""},
		{Text: "b1"},
		{Text: "    "},
		{Text: "\t"},
		{Text: "c1"},
		{Text: ""},
	})
	require.Len(t, blocks, 3)

	require.Len(t, blocks[0], 5)
	assert.Len(t, blocks[0].SignificantLines(), 2)
	assert.Equal(t, blocks[0][2].Text, "a1")
	assert.Equal(t, blocks[0][3].Text, "a2")

	require.Len(t, blocks[1], 3)
	assert.Len(t, blocks[1].SignificantLines(), 1)
	assert.Equal(t, blocks[1][0].Text, "b1")

	require.Len(t, blocks[2], 2)
	assert.Len(t, blocks[2].SignificantLines(), 1)
	assert.Equal(t, blocks[2][0].Text, "c1")
}

func TestDisregardLinesAllEmpty(t *testing.T) {
	blocks := GroupIntoBlocks([]Line{
		{Text: ""},
		{Text: "  "},
		{Text: "\t"},
		{Text: ""},
	})
	require.Nil(t, blocks)
}
