package terminalformat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var reflower = NewReflower(60, "\n")

func TestLineBreakerReflowsText(t *testing.T) {
	original := `This is a very long line and it should be reflowed so that it doesn’t run so wide, because that’s easier to read and it looks better.
However, existing line breaks are respected.

The same is true for blank-line-separated paragraphs`
	assert.Equal(t, `This is a very long line and it should be reflowed so that
it doesn’t run so wide, because that’s easier to read and
it looks better.
However, existing line breaks are respected.

The same is true for blank-line-separated paragraphs`, reflower.Reflow(original, ""))
}

func TestLineBreakerDoesNotDoAnythingIfEmptyInput(t *testing.T) {
	assert.Equal(t, "", reflower.Reflow("", ""))
	assert.Equal(t, "", reflower.Reflow("   ", ""))
	assert.Equal(t, "\n", reflower.Reflow("\n", ""))
}

func TestLineBreakerPrependsPrefix(t *testing.T) {
	original := "This is a very long line and it should be reflowed so that it doesn’t run so wide, because that’s easier to read."
	assert.Equal(t, `  This is a very long line and it should be reflowed so that
  it doesn’t run so wide, because that’s easier to read.`, reflower.Reflow(original, "  "))
}
