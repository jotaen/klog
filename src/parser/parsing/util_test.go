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

func TestInsertInBetween(t *testing.T) {
	before := Split("first\nthird\nfourth")
	after := Insert(before, 1, []Text{
		{"second", 0},
	}, DefaultPreferences())
	require.Len(t, after, 4)
	assert.Equal(t, before[0].Original(), after[0].Original())
	assert.Equal(t, 1, after[0].LineNumber)

	assert.Equal(t, "second\n", after[1].Original())
	assert.Equal(t, 2, after[1].LineNumber)

	assert.Equal(t, before[1].Original(), after[2].Original())
	assert.Equal(t, 3, after[2].LineNumber)

	assert.Equal(t, before[2].Original(), after[3].Original())
	assert.Equal(t, 4, after[3].LineNumber)
}

func TestInsertMultipleTexts(t *testing.T) {
	before := Split("first\nfourth\nfifth\n")
	after := Insert(before, 1, []Text{
		{"second", 0},
		{"third", 1},
	}, DefaultPreferences())
	require.Len(t, after, 5)
	assert.Equal(t, "first\n", after[0].Original())
	assert.Equal(t, 1, after[0].LineNumber)
	assert.Equal(t, "second\n", after[1].Original())
	assert.Equal(t, 2, after[1].LineNumber)
	assert.Equal(t, "    third\n", after[2].Original())
	assert.Equal(t, 3, after[2].LineNumber)
	assert.Equal(t, "fourth\n", after[3].Original())
	assert.Equal(t, 4, after[3].LineNumber)
	assert.Equal(t, "fifth\n", after[4].Original())
	assert.Equal(t, 5, after[4].LineNumber)
}

func TestInsertWithLineEndingsAndIndentation(t *testing.T) {
	before := Split("bar")
	after := Insert(before, 0, []Text{{"foo", 0}}, DefaultPreferences())
	after = Insert(after, 2, []Text{{"baz", 1}}, Preferences{"\r\n", "\t"})
	after = Insert(after, 0, []Text{{"hello", 1}}, DefaultPreferences())
	require.Len(t, after, 4)
	assert.Equal(t, "    hello\n", after[0].Original())
	assert.Equal(t, "foo\n", after[1].Original())
	assert.Equal(t, "bar\r\n", after[2].Original())
	assert.Equal(t, "\tbaz\r\n", after[3].Original())
}

func TestInsertIntoEmptySlice(t *testing.T) {
	var before []Line
	after := Insert(before, 0, []Text{{"Hello World", 0}}, DefaultPreferences())
	require.Len(t, after, 1)
	assert.Equal(t, "Hello World\n", after[0].Original())
}
