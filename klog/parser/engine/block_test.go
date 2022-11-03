package engine

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGroupEmptyInput(t *testing.T) {
	for _, ls := range []string{
		``,
		"\n \n\t\t\n  ",
	} {
		blocks := GroupIntoBlocks(ls)
		require.Nil(t, blocks)
	}
}

func TestGroupLinesOfSingleBlock(t *testing.T) {
	for _, x := range []struct {
		text string
		ls   int
	}{
		{"a1", 1},
		{"\na1", 2},
		{"\na1\n", 2},
		{"   \na1", 2},
		{"   \na1\n", 2},
		{"\t\na1\n\t \t ", 3},
		{"\n\na1\n\n", 4},
	} {
		blocks := GroupIntoBlocks(x.text)
		require.Len(t, blocks, 1)

		require.Len(t, blocks[0], x.ls)
		require.Len(t, blocks[0].SignificantLines(), 1)
		assert.Equal(t, blocks[0].SignificantLines()[0].Text, "a1")
	}
}

func TestGroupLinesOfSingleBlockWithMultipleLines(t *testing.T) {
	for _, x := range []struct {
		text string
		ls   int
	}{
		{"a1\na2", 2},
		{"\na1\na2", 3},
		{"\na1\na2\n", 3},
		{"\n    \t\na1\na2", 4},
		{"\n    \t\na1\na2\n", 4},
		{"\n    \t\na1\na2\n\n", 5},
		{" \t \t\na1\na2\n\t", 4},
		{" \n\t\na1\na2\n\n", 5},
	} {
		blocks := GroupIntoBlocks(x.text)
		require.Len(t, blocks, 1)

		require.Len(t, blocks[0], x.ls)
		require.Len(t, blocks[0].SignificantLines(), 2)
		assert.Equal(t, blocks[0].SignificantLines()[0].Text, "a1")
		assert.Equal(t, blocks[0].SignificantLines()[1].Text, "a2")
	}
}

func TestGroupLinesOfMultipleBlocks(t *testing.T) {
	blocks := GroupIntoBlocks(`
  
a1
a2

b1
        
	
c1

`)
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
	blocks := GroupIntoBlocks(`

 
	
 

`)
	require.Nil(t, blocks)
}
