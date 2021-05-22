package terminalformat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStripAllAnsiSequences(t *testing.T) {
	assert.Equal(t, "test 123", StripAllAnsiSequences("test 123"))
	assert.Equal(t, "test          123", StripAllAnsiSequences("test          123"))
	assert.Equal(t, "test 123", StripAllAnsiSequences("test \x1b[0m\x1b[4m123\x1b[0m"))
	assert.Equal(t, "test 123", StripAllAnsiSequences("\x1b[0m\x1b[4mtest\x1b[0m \x1b[0m\x1b[4m123\x1b[0m"))
}
