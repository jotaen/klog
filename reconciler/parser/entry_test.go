package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinimalValidEntry(t *testing.T) {
	yaml := `
date: 2020-01-01
`
	e, errs := Parse(yaml)
	assert.Equal(t, e.Summary(), "")
	assert.Nil(t, errs)
}

func TestFailOnUnknownProperties(t *testing.T) {
	yaml := `
date: 2020-01-01
foo: 1
bar: test
`
	e, errs := Parse(yaml)
	assert.Equal(t, e, nil)
	assert.Contains(t, errs, parserError(MALFORMED_YAML))
}
