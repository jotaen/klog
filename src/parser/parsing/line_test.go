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
	assert.Equal(t, "\n", ls[0].LineEnding)
	assert.Equal(t, "bar", ls[1].Text)
	assert.Equal(t, "\r\n", ls[1].LineEnding)
	assert.Equal(t, "baz", ls[2].Text)
	assert.Equal(t, "", ls[2].LineEnding)
}

func TestDeterminesPrecedingWhitespace(t *testing.T) {
	for _, x := range []struct {
		text               string
		expectedWhitespace string
	}{
		{" one space", " "},
		{"  two spaces", "  "},
		{"   three spaces", "   "},
		{"    four spaces", "    "},
		{"\tone tab", "\t"},
		{"\t\ttwo tabs", "\t\t"},
		{"   \t     \t\t \t   wild mix", "   \t     \t\t \t   "},
	} {
		line := NewLineFromString(x.text, 0)
		assert.Equal(t, x.expectedWhitespace, line.PrecedingWhitespace)
	}
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
