package reconciler

import (
	"github.com/jotaen/klog/src/parser/lineparsing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInsertInBetween(t *testing.T) {
	before := lineparsing.Split("first\nthird\nfourth")
	after := insert(before, 1, []InsertableText{
		{"second", 0},
	}, stylePreferencesOrDefault(nil))
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

func TestInsertAtBeginningAndEnd(t *testing.T) {
	before := lineparsing.Split("beginning\nend")
	after := insert(before, 0, []InsertableText{
		{"first", 0},
	}, stylePreferencesOrDefault(nil))
	after = insert(after, 3, []InsertableText{
		{"last", 0},
	}, stylePreferencesOrDefault(nil))
	require.Len(t, after, 4)
	assert.Equal(t, "first\n", after[0].Original())
	assert.Equal(t, "beginning\n", after[1].Original())
	assert.Equal(t, "end\n", after[2].Original())
	assert.Equal(t, "last\n", after[3].Original())
}

func TestInsertMultipleTexts(t *testing.T) {
	before := lineparsing.Split("first\nfourth\nfifth\n")
	after := insert(before, 1, []InsertableText{
		{"second", 0},
		{"third", 1},
	}, stylePreferencesOrDefault(nil))
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
	before := lineparsing.Split("bar")
	after := insert(before, 0, []InsertableText{{"foo", 0}}, stylePreferencesOrDefault(nil))
	after = insert(after, 2, []InsertableText{{"baz", 1}}, stylePreferences{"\t", "\r\n"})
	after = insert(after, 0, []InsertableText{{"hello", 1}}, stylePreferencesOrDefault(nil))
	require.Len(t, after, 4)
	assert.Equal(t, "    hello\n", after[0].Original())
	assert.Equal(t, "foo\n", after[1].Original())
	assert.Equal(t, "bar\r\n", after[2].Original())
	assert.Equal(t, "\tbaz\r\n", after[3].Original())
}

func TestInsertIntoEmptySlice(t *testing.T) {
	var before []lineparsing.Line
	after := insert(before, 0, []InsertableText{{"Hello World", 0}}, stylePreferencesOrDefault(nil))
	require.Len(t, after, 1)
	assert.Equal(t, "Hello World\n", after[0].Original())
}
