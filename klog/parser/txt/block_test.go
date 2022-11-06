package txt

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
		block, _ := ParseBlock(ls, 0)
		assert.Nil(t, block)
	}
}

func TestGroupLinesOfSingleBlock(t *testing.T) {
	for _, x := range []struct {
		text  string
		count int
	}{
		{"a1", 1},
		{"\na1", 2},
		{"\na1\n", 2},
		{"   \na1", 2},
		{"   \na1\n", 2},
		{"\t\na1\n\t \t ", 3},
		{"\n\na1\n\n", 4},
	} {
		block, _ := ParseBlock(x.text, 0)

		require.NotNil(t, block)
		require.Len(t, block.Lines(), x.count)
		sgLines, _, _ := block.SignificantLines()
		require.Len(t, sgLines, 1)
		assert.Equal(t, sgLines[0].Text, "a1")
	}
}

func TestGroupLinesOfSingleBlockWithMultipleLines(t *testing.T) {
	for _, x := range []struct {
		text  string
		count int
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
		block, _ := ParseBlock(x.text, 0)

		require.NotNil(t, block)
		require.Len(t, block.Lines(), x.count)
		sgLines, _, _ := block.SignificantLines()
		require.Len(t, sgLines, 2)
		assert.Equal(t, sgLines[0].Text, "a1")
		assert.Equal(t, sgLines[1].Text, "a2")
	}
}
