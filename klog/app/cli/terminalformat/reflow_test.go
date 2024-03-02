package terminalformat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLineBreakerReflowsText(t *testing.T) {
	reflower := NewReflower(60, "\n")
	original := `This is a very long line and it should be reflowed so that it doesn’t run so wide, because that’s easier to read and it looks better.
However, existing line breaks are respected.

The same is true for blank-line-separated paragraphs`
	assert.Equal(t, `This is a very long line and it should be reflowed so that
it doesn’t run so wide, because that’s easier to read and
it looks better.
However, existing line breaks are respected.

The same is true for blank-line-separated paragraphs`, reflower.Reflow(original, nil))
}

func TestLineBreakerDoesNotDoAnythingIfEmptyInput(t *testing.T) {
	reflower := NewReflower(60, "\n")
	assert.Equal(t, "", reflower.Reflow("", nil))
	assert.Equal(t, "", reflower.Reflow("   ", nil))
	assert.Equal(t, "\n", reflower.Reflow("\n", nil))
}

func TestLineBreakerPrependsPrefix(t *testing.T) {
	reflower := NewReflower(60, "\n")
	original := "This is a very long line and it should be reflowed so that it doesn’t run so wide, because that’s easier to read."
	assert.Equal(t, `  This is a very long line and it should be reflowed so that
  it doesn’t run so wide, because that’s easier to read.`, reflower.Reflow(original, []string{"  "}))
}

func TestLineBreakerPrependsMultiplePrefixes(t *testing.T) {
	reflower := NewReflower(30, "\n")
	original := "This is a very long line and it should be reflowed so that it doesn’t run so wide, because that’s easier to read."
	assert.Equal(t, `This is a very long line and
| it should be reflowed so that
| it doesn’t run so wide,
| because that’s easier to read.`, reflower.Reflow(original, []string{"", "| "}))
}
