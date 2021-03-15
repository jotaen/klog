package parsing

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSplitTextIntoLines(t *testing.T) {
	text := "foo\nbar\r\n\n  \nbaz"
	ls := Split(text)
	require.Len(t, ls, 5)

	assert.Equal(t, ls[0].Original(), "foo\n")
	assert.Equal(t, ls[0].LineNumber, 1)

	assert.Equal(t, ls[1].Original(), "bar\r\n")
	assert.Equal(t, ls[1].LineNumber, 2)

	assert.Equal(t, ls[2].Original(), "\n")
	assert.Equal(t, ls[2].LineNumber, 3)

	assert.Equal(t, ls[3].Original(), "  \n")
	assert.Equal(t, ls[3].LineNumber, 4)

	assert.Equal(t, ls[4].Original(), "baz")
	assert.Equal(t, ls[4].LineNumber, 5)
}

func TestDeterminesLineEndings(t *testing.T) {
	text := "foo\nbar\r\nbaz"
	ls := Split(text)
	require.Len(t, ls, 3)
	assert.Equal(t, "foo", ls[0].Text)
	assert.Equal(t, "\n", ls[0].originalLineEnding)
	assert.Equal(t, "bar", ls[1].Text)
	assert.Equal(t, "\r\n", ls[1].originalLineEnding)
	assert.Equal(t, "baz", ls[2].Text)
	assert.Equal(t, "", ls[2].originalLineEnding)
}

func TestDeterminesIndentation(t *testing.T) {
	text := "  two spaces\n   three spaces\n    four spaces\n\tone tab"
	ls := Split(text)
	require.Len(t, ls, 4)
	assert.Equal(t, "  ", ls[0].originalIndentation)
	assert.Equal(t, 1, ls[0].IndentationLevel())
	assert.Equal(t, "   ", ls[1].originalIndentation)
	assert.Equal(t, 1, ls[1].IndentationLevel())
	assert.Equal(t, "    ", ls[2].originalIndentation)
	assert.Equal(t, 1, ls[2].IndentationLevel())
	assert.Equal(t, "\t", ls[3].originalIndentation)
	assert.Equal(t, 1, ls[3].IndentationLevel())
}

func TestInvalidIndentation(t *testing.T) {
	text := "     NO: five spaces\n\t\tNO: two tabs\n NO: one space"
	ls := Split(text)
	require.Len(t, ls, 3)
	assert.Equal(t, "     ", ls[0].originalIndentation)
	assert.Equal(t, "     NO: five spaces\n", ls[0].Original())
	assert.Equal(t, -1, ls[0].IndentationLevel())
	assert.Equal(t, "\t\t", ls[1].originalIndentation)
	assert.Equal(t, "\t\tNO: two tabs\n", ls[1].Original())
	assert.Equal(t, -1, ls[1].IndentationLevel())
	assert.Equal(t, " ", ls[2].originalIndentation)
	assert.Equal(t, -1, ls[2].IndentationLevel())
	assert.Equal(t, " NO: one space", ls[2].Original())
}

func TestToStringRestoresOriginal(t *testing.T) {
	text := "  Hello World\n\tTest 123\r\n      Foo Bar BAZ"
	ls := Split(text)
	require.Len(t, ls, 3)
	assert.Equal(t, "  Hello World\n", ls[0].Original())
	assert.Equal(t, "\tTest 123\r\n", ls[1].Original())
	assert.Equal(t, "      Foo Bar BAZ", ls[2].Original())
}

func TestSplitAndJoinResultsInOriginalText(t *testing.T) {
	text := "x\n1293871jh23981y293j\n  asdfkj     askdlfjh\n\nalkdjhf\r\n\tasdkljfh\n"
	assert.Equal(t, text, Join(Split(text)))
}
