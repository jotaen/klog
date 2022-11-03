package txt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeterminesLineEndings(t *testing.T) {
	ls := []Line{NewLineFromString(
		"foo\n", 1),
		NewLineFromString("bar\r\n", 2),
		NewLineFromString("baz", 3),
	}
	assert.Equal(t, "foo", ls[0].Text)
	assert.Equal(t, "\n", ls[0].LineEnding)
	assert.Equal(t, "bar", ls[1].Text)
	assert.Equal(t, "\r\n", ls[1].LineEnding)
	assert.Equal(t, "baz", ls[2].Text)
	assert.Equal(t, "", ls[2].LineEnding)
}
