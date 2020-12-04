package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseSummary(t *testing.T) {
	yaml := `
date: 2020-01-01
summary: Just a normal day
`
	e, _ := Parse(yaml)
	assert.Equal(t, e.Summary, "Just a normal day")
}
