package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"klog/parser/parsing"
	"testing"
	"time"
)

var now = time.Date(1995, 3, 31, 13, 15, 29, 0, time.UTC)

func TestRenderTemplate(t *testing.T) {
	result, err := RenderTemplate(`
{{ TODAY }}
Foo #xyz

{{YESTERDAY}} (8h30m!)
	1h
	{{ NOW }} - ?
`, now)
	require.Nil(t, err)
	assert.Equal(t, []parsing.Text{
		{"", 0},
		{"1995-03-31", 0},
		{"Foo #xyz", 0},
		{"", 0},
		{"1995-03-30 (8h30m!)", 0},
		{"1h", 1},
		{"13:15 - ?", 1},
	}, result)
}

func TestTemplateFailsIfNoValidRecord(t *testing.T) {
	result, err := RenderTemplate(`
{{ TODAY }} foo
	This is all malformed
`, now)
	require.Error(t, err)
	require.Nil(t, result)
}
