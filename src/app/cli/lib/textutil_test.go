package lib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLineBreakerReflowsText(t *testing.T) {
	original := `This is a very long line and it should be reflowed so that it doesn’t run so wide, because that’s easier to read and it looks better.
However, existing line breaks are respected.

The same is true for blank-line-separated paragraphs`
	assert.Equal(t, `This is a very long line and it should be reflowed so that
it doesn’t run so wide, because that’s easier to read and
it looks better.
However, existing line breaks are respected.

The same is true for blank-line-separated paragraphs`, lineBreaker.reflow(original, ""))
}

func TestLineBreakerDoesNotDoAnythingIfEmptyInput(t *testing.T) {
	assert.Equal(t, "", lineBreaker.reflow("", ""))
	assert.Equal(t, "", lineBreaker.reflow("   ", ""))
	assert.Equal(t, "\n", lineBreaker.reflow("\n", ""))
}

func TestLineBreakerPrependsPrefix(t *testing.T) {
	original := "This is a very long line and it should be reflowed so that it doesn’t run so wide, because that’s easier to read."
	assert.Equal(t, `  This is a very long line and it should be reflowed so that
  it doesn’t run so wide, because that’s easier to read.`, lineBreaker.reflow(original, "  "))
}
