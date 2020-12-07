package parser

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAbsentDatePropertyFails(t *testing.T) {
	yaml := `
summary: Just a normal day
`
	w, errs := Parse(yaml)
	assert.Equal(t, w, nil)
	assert.Contains(t, errs, errors.New("INVALID_DATE"))
}

func TestMalformedDateFails(t *testing.T) {
	yaml := `
date: 01.01.2020
`
	w, errs := Parse(yaml)
	assert.Equal(t, w, nil)
	assert.Contains(t, errs, errors.New("INVALID_DATE"))
}
