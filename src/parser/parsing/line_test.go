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

	assert.Equal(t, ls[0].ToString(), "foo\n")
	assert.Equal(t, ls[0].LineNumber, 1)

	assert.Equal(t, ls[1].ToString(), "bar\r\n")
	assert.Equal(t, ls[1].LineNumber, 2)

	assert.Equal(t, ls[2].ToString(), "\n")
	assert.Equal(t, ls[2].LineNumber, 3)

	assert.Equal(t, ls[3].ToString(), "  \n")
	assert.Equal(t, ls[3].LineNumber, 4)

	assert.Equal(t, ls[4].ToString(), "baz")
	assert.Equal(t, ls[4].LineNumber, 5)
}

func TestDeterminesLineEndings(t *testing.T) {
	text := "foo\nbar\r\nbaz"
	ls := Split(text)
	require.Len(t, ls, 3)
	assert.Equal(t, "foo", ls[0].Text)
	assert.Equal(t, "\n", ls[0].lineEnding)
	assert.Equal(t, "bar", ls[1].Text)
	assert.Equal(t, "\r\n", ls[1].lineEnding)
	assert.Equal(t, "baz", ls[2].Text)
	assert.Equal(t, "", ls[2].lineEnding)
}

func TestDeterminesIndentation(t *testing.T) {
	text := "  two spaces\n   three spaces\n    four spaces\n     five spaces\n\tone tab\n invalid"
	ls := Split(text)
	require.Len(t, ls, 6)
	assert.Equal(t, "  ", ls[0].indentation)
	assert.Equal(t, 1, ls[0].IndentationLevel())
	assert.Equal(t, "   ", ls[1].indentation)
	assert.Equal(t, 1, ls[1].IndentationLevel())
	assert.Equal(t, "    ", ls[2].indentation)
	assert.Equal(t, 1, ls[2].IndentationLevel())
	assert.Equal(t, "     ", ls[3].indentation)
	assert.Equal(t, 2, ls[3].IndentationLevel())
	assert.Equal(t, "\t", ls[4].indentation)
	assert.Equal(t, 1, ls[4].IndentationLevel())
	assert.Equal(t, " ", ls[5].indentation)
	assert.Equal(t, -1, ls[5].IndentationLevel())
}

func TestToStringRestoresOriginal(t *testing.T) {
	text := "  Hello World\n\tTest 123\r\n      Foo Bar BAZ"
	ls := Split(text)
	require.Len(t, ls, 3)
	assert.Equal(t, "  Hello World\n", ls[0].ToString())
	assert.Equal(t, "\tTest 123\r\n", ls[1].ToString())
	assert.Equal(t, "      Foo Bar BAZ", ls[2].ToString())
}

func TestSplitAndJoinResultsInOriginalText(t *testing.T) {
	text := "x\n1293871jh23981y293j\n  asdfkj     askdlfjh\n\nalkdjhf\r\n\tasdkljfh\n"
	assert.Equal(t, text, Join(Split(text)))
}

func TestJoinAddsMissingLineEndings(t *testing.T) {
	ls := []Line{
		NewLineFromString("First Line", 1),
		NewLineFromString("Second Line\n", 1),
		NewLineFromString("Third Line", 2),
	}
	text := Join(ls)
	assert.Equal(t, "First Line\nSecond Line\nThird Line\n", text)
}

func TestJoinAddsMissingLineEndingsAndGuessesFromPreviousValue(t *testing.T) {
	ls := []Line{
		NewLineFromString("First Line\r\n", 1),
		NewLineFromString("Second Line", 2),
	}
	text := Join(ls)
	assert.Equal(t, "First Line\r\nSecond Line\r\n", text)
}

func TestInsertInBetween(t *testing.T) {
	before := Split("first\nthird\nfourth")
	after := Insert(before, 1, "second\n")
	require.Len(t, after, 4)
	assert.Equal(t, before[0].ToString(), after[0].ToString())
	assert.Equal(t, after[0].LineNumber, 1)

	assert.Equal(t, "second\n", after[1].ToString())
	assert.Equal(t, after[1].LineNumber, 2)

	assert.Equal(t, before[1].ToString(), after[2].ToString())
	assert.Equal(t, after[2].LineNumber, 3)

	assert.Equal(t, before[2].ToString(), after[3].ToString())
	assert.Equal(t, after[3].LineNumber, 4)
}

func TestInsert(t *testing.T) {
	before := Split("bar")
	after := Insert(before, 0, "foo\n")
	after = Insert(after, 2, "baz")
	require.Len(t, after, 3)
	assert.Equal(t, after[0].ToString(), "foo\n")
	assert.Equal(t, after[1].ToString(), "bar")
	assert.Equal(t, after[2].ToString(), "baz")
}
