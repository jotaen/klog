package txt

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateIndentatorFromLine(t *testing.T) {
	for _, indentator := range []*Indentator{
		NewIndentator([]string{"  ", "    "}, NewLineFromString("  Hello")),
		NewIndentator([]string{"  ", "    "}, NewLineFromString("    Hello")),
		NewIndentator([]string{"  ", "    "}, NewLineFromString("          Hello")),
		NewIndentator([]string{"\t"}, NewLineFromString("\tHello")),
	} {
		require.NotNil(t, indentator)
	}
}

func TestCreatesNoIndentatorIfLineIsNotIndentatedAccordingly(t *testing.T) {
	for _, indentator := range []*Indentator{
		NewIndentator([]string{"  ", "    "}, NewLineFromString("Hello")),
		NewIndentator([]string{"  ", "    "}, NewLineFromString(" Hello")),
		NewIndentator([]string{"\t"}, NewLineFromString("  Hello")),
	} {
		require.Nil(t, indentator)
	}
}

func TestCreatesIndentedParseable(t *testing.T) {
	indentator := Indentator{"\t"}

	p1 := indentator.NewIndentedParseable(NewLineFromString("Hello"), 0)
	require.NotNil(t, p1)
	assert.Equal(t, p1.PointerPosition, 0)
	assert.Equal(t, []rune{'H', 'e', 'l', 'l', 'o'}, p1.Chars)

	p2 := indentator.NewIndentedParseable(NewLineFromString("\tHello"), 1)
	require.NotNil(t, p2)
	assert.Equal(t, 1, p2.PointerPosition)

	p3 := indentator.NewIndentedParseable(NewLineFromString("\t\tHello"), 1)
	require.NotNil(t, p3)
	assert.Equal(t, 1, p3.PointerPosition)
}

func TestCreatesNoParseableForIndentationMismatch(t *testing.T) {
	indentator := Indentator{"\t"}
	for _, p := range []*Parseable{
		indentator.NewIndentedParseable(NewLineFromString("Hello"), 1),
		indentator.NewIndentedParseable(NewLineFromString("\tHello"), 2),
		indentator.NewIndentedParseable(NewLineFromString("\t\tHello"), 5),
	} {
		require.Nil(t, p)
	}
}
