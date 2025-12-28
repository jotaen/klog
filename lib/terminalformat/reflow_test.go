package terminalformat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReflowsText(t *testing.T) {
	reflower := NewReflower(60)
	original := `This is a very long line and it should be reflowed so that it doesn’t run so wide, because that’s easier to read and it looks better.`
	assert.Equal(t, `This is a very long line and it should be reflowed so that
it doesn’t run so wide, because that’s easier to read and it
looks better.`, reflower.Reflow(original, ""))
}

func TestReflowsTextThatMeetsLineLength(t *testing.T) {
	reflower := NewReflower(38)
	original := `This is a very long line and it should
not be reflowed so that it doesn’t run
so wide, because that’s easier to read
and it looks better.`
	assert.Equal(t, `This is a very long line and it should
not be reflowed so that it doesn’t run
so wide, because that’s easier to read
and it looks better.`, reflower.Reflow(original, ""))
}

func TestReflowsTextThatAlmostMeetsLineLength(t *testing.T) {
	reflower := NewReflower(39)
	original := `This is a very long line and it should not be re-flowed so that it does not run so wide, because that’s easier to read and and it looks better.`
	assert.Equal(t, `This is a very long line and it should
not be re-flowed so that it does not
run so wide, because that’s easier to
read and and it looks better.`, reflower.Reflow(original, ""))
}

func TestReflowsTextRespectLineBreaks(t *testing.T) {
	reflower := NewReflower(60)
	original := `Existing line breaks are respected.
Even
for
super
short
lines.`
	assert.Equal(t, `Existing line breaks are respected.
Even
for
super
short
lines.`, reflower.Reflow(original, ""))
}

func TestReflowsTextPreservesParagraphs(t *testing.T) {
	reflower := NewReflower(40)
	original := `For a text that consists of multiple paragraphs, the reflower preserves the existing line breaks.

See? It only reflows within the respective paragraphs!



And it even works for multiple blank lines in between.`
	assert.Equal(t, `For a text that consists of multiple
paragraphs, the reflower preserves the
existing line breaks.

See? It only reflows within the
respective paragraphs!



And it even works for multiple blank
lines in between.`, reflower.Reflow(original, ""))
}

func TestReflowsTextWithVeryLongWords(t *testing.T) {
	reflower := NewReflower(18)
	original := `But what happens if thereIsASingleWordThatIsSuperDuperLong? How does it behave then?`
	assert.Equal(t, `But what happens
if
thereIsASingleWord
ThatIsSuperDuperLo
ng? How does it
behave then?`, reflower.Reflow(original, ""))
}

func TestReflowsTextWithUtf8(t *testing.T) {
	reflower := NewReflower(10)
	original := `😊😊😊😊😊😊😊😊 😊😊😊 😊😊😊 😊😊😊😊 😊 😊😊😊😊😊 😊😊😊 😊😊.`
	assert.Equal(t, `😊😊😊😊😊😊😊😊
😊😊😊 😊😊😊
😊😊😊😊 😊
😊😊😊😊😊 😊😊😊
😊😊.`, reflower.Reflow(original, ""))
}

func TestReflowDoesNotDoAnythingIfEmptyInput(t *testing.T) {
	reflower := NewReflower(60)
	assert.Equal(t, "", reflower.Reflow("", ""))
	assert.Equal(t, "   ", reflower.Reflow("   ", ""))
	assert.Equal(t, "\n", reflower.Reflow("\n", ""))
	assert.Equal(t, "\n\n", reflower.Reflow("\n\n", ""))
}

func TestReflowPrependsPrefix(t *testing.T) {
	reflower := NewReflower(45)
	original := "This is a very long line and it should be reflowed so that it doesn’t run so wide, because that’s easier to read."
	assert.Equal(t, ` | This is a very long line and it should be
 | reflowed so that it doesn’t run so wide,
 | because that’s easier to read.`, reflower.Reflow(original, " | "))
}

func TestReflowHandlesLeadingAndTrailingBlankLines(t *testing.T) {
	reflower := NewReflower(100)
	original := `

This is a very long line and it should be reflowed so that it doesn’t run so wide, because that’s easier to read.


`
	assert.Equal(t, `

This is a very long line and it should be reflowed so that it doesn’t run so wide, because that’s
easier to read.


`, reflower.Reflow(original, ""))
}

func TestReflowPreservesLeadingWhiteSpaceInLines(t *testing.T) {
	reflower := NewReflower(50)
	original := `This is a very long line and it should be reflowed.
    This line has leading whitespace.
      This one too.
         Leading white space should be preserved. Also, if an indented line overflows, the indentation should be preserved until the next line break.
See?`
	assert.Equal(t, `This is a very long line and it should be
reflowed.
    This line has leading whitespace.
      This one too.
         Leading white space should be preserved.
         Also, if an indented line overflows, the
         indentation should be preserved until the
         next line break.
See?`, reflower.Reflow(original, ""))
}

func TestReflowDisregardsRedundantWhiteSpace(t *testing.T) {
	reflower := NewReflower(50)
	original := `This    is    a        very       long   line    and    it   should     be       reflowed.      `
	assert.Equal(t, `This is a very long line and it should be
reflowed.`, reflower.Reflow(original, ""))
}
