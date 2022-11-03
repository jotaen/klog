package txt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateError(t *testing.T) {
	err := NewError(Line{
		Text:       "Hello World",
		LineNumber: 2,
		LineEnding: "\n",
	}, 0, 5, "CODE", "Title", "Details")
	assert.Equal(t, "CODE", err.Code())
	assert.Equal(t, "Title", err.Title())
	assert.Equal(t, "Details", err.Details())
	assert.Equal(t, 0, err.Position())
	assert.Equal(t, 1, err.Column())
	assert.Equal(t, 5, err.Length())
	assert.Equal(t, "Title: Details", err.Message())
}
