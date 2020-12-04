package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinimalValidWorkDay(t *testing.T) {
	yaml := `
date: 2020-01-01
`
	w, errs := Parse(yaml)
	assert.Equal(t, w.Summary(), "")
	assert.Nil(t, errs)
}

func TestFailOnUnknownProperties(t *testing.T) {
	yaml := `
date: 2020-01-01
foo: 1
bar: test
`
	w, errs := Parse(yaml)
	assert.Equal(t, w, nil)
	assert.Contains(t, errs, parserError(MALFORMED_YAML))
}
