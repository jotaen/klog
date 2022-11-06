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
		text      string
		expect    string
		lineCount int
	}{
		{"a", "a", 1},
		{"\nfoo", "foo", 2},
		{"\n12345\n", "12345", 2},
		{"   \ntest ", "test ", 2},
		{"   \na\ta\n", "a\ta", 2},
		{"\t\na1\n\t \t ", "a1", 3},
		{"\n\na1\n\n", "a1", 4},
		{"喜左衛門", "喜左衛門", 1},
		{"喜左衛門\n", "喜左衛門", 1},
		{"😀·½\n", "😀·½", 1},
	} {
		block, _ := ParseBlock(x.text, 0)

		require.NotNil(t, block)
		require.Len(t, block.Lines(), x.lineCount)
		sgLines, _, _ := block.SignificantLines()
		require.Len(t, sgLines, 1)
		assert.Equal(t, sgLines[0].Text, x.expect)
	}
}

func TestGroupLinesOfSingleBlockWithMultipleLines(t *testing.T) {
	for _, x := range []struct {
		text      string
		expect    []string // significant lines
		lineCount int
	}{
		{"a1\na2", []string{"a1", "a2"}, 2},
		{"\nasdf\nasdf", []string{"asdf", "asdf"}, 3},
		{"\nHey 🥰!\n«How is it?»\n", []string{"Hey 🥰!", "«How is it?»"}, 3},
		{"\n    \t\nA\nB", []string{"A", "B"}, 4},
		{"\n    \t\na b c \n a b c\n", []string{"a b c ", " a b c"}, 4},
		{"\n    \t\n       _       \n     -     \n\n", []string{"       _       ", "     -     "}, 5},
		{" \t \t\nAS:FLKJH\n!(@* #&\n\t", []string{"AS:FLKJH", "!(@* #&"}, 4},
		{" \n\t\n1—2\n·½⅓•ÄﬂÑ\n\n", []string{"1—2", "·½⅓•ÄﬂÑ"}, 5},
	} {
		block, _ := ParseBlock(x.text, 0)

		require.NotNil(t, block)
		require.Len(t, block.Lines(), x.lineCount)
		sgLines, _, _ := block.SignificantLines()
		require.Len(t, sgLines, 2)
		assert.Equal(t, sgLines[0].Text, x.expect[0])
		assert.Equal(t, sgLines[1].Text, x.expect[1])
	}
}
