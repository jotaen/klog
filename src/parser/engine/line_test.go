package engine

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSplitTextIntoLines(t *testing.T) {
	text := "foo\nbar\r\n\n \nbaz"
	ls := Split(text)
	require.Len(t, ls, 5)

	assert.Equal(t, ls[0].Original, "foo\n")
	assert.Equal(t, ls[0].LineNumber, 1)

	assert.Equal(t, ls[1].Original, "bar\r\n")
	assert.Equal(t, ls[1].LineNumber, 2)

	assert.Equal(t, ls[2].Original, "\n")
	assert.Equal(t, ls[2].LineNumber, 3)

	assert.Equal(t, ls[3].Original, " \n")
	assert.Equal(t, ls[3].LineNumber, 4)

	assert.Equal(t, ls[4].Original, "baz")
	assert.Equal(t, ls[4].LineNumber, 5)
}

func TestStripsLineEndingsFromValues(t *testing.T) {
	text := "foo\nbar\r\n"
	ls := Split(text)
	require.Len(t, ls, 2)
	assert.Equal(t, ls[0].Value, []rune{'f', 'o', 'o'})
	assert.Equal(t, ls[1].Value, []rune{'b', 'a', 'r'})
}

func TestRecognisesIndentation(t *testing.T) {
	text := "  x\n\ty\n        z"
	ls := Split(text)
	require.Len(t, ls, 3)
	assert.Equal(t, ls[0].Value, []rune{'x'})
	assert.Equal(t, ls[0].IndentationLevel, 1)
	assert.Equal(t, ls[1].Value, []rune{'y'})
	assert.Equal(t, ls[1].IndentationLevel, 1)
	assert.Equal(t, ls[2].Value, []rune{'z'})
	assert.Equal(t, ls[2].IndentationLevel, 2)
}

func TestRejectsInvalidIndentation(t *testing.T) {
	text := " x\n     y"
	ls := Split(text)
	require.Len(t, ls, 2)
	assert.Equal(t, ls[0].Value, []rune("x"))
	assert.Less(t, ls[0].IndentationLevel, 0)
	assert.Equal(t, ls[1].Value, []rune("y"))
	assert.Less(t, ls[1].IndentationLevel, 0)
}

func TestSplitAndJoinResultsInOriginalText(t *testing.T) {
	text := "x\n1293871jh23981y293j\n asdfkj     askdlfjh\n\nalkdjhf\r\n\tasdkljfh\n"
	assert.Equal(t, text, Join(Split(text)))
}
