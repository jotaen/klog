package prettify

import (
	"errors"
	"testing"

	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/parser/txt"
	tf "github.com/jotaen/klog/lib/terminalformat"
	"github.com/stretchr/testify/assert"
)

var styler = tf.NewStyler(tf.COLOUR_THEME_NO_COLOUR)

func TestFormatParserError(t *testing.T) {
	block1, _ := txt.ParseBlock("Good text\nSome malformed text", 37)
	block2, _ := txt.ParseBlock("Another issue!", 133)
	err := app.NewParserErrors([]txt.Error{
		txt.NewError(block1, 1, 0, 4, "CODE", "Error", "Short explanation."),
		txt.NewError(block2, 0, 8, 5, "CODE", "Problem", "More info.").SetOrigin("some-file.klg"),
	})
	text := PrettifyParsingError(err, styler).Error()
	assert.Equal(t, `
[SYNTAX ERROR] in line 39
    Some malformed text
    ^^^^
    Error: Short explanation.

[SYNTAX ERROR] in line 134 of file some-file.klg
    Another issue!
            ^^^^^
    Problem: More info.
`, tf.StripAllAnsiSequences(text))
}

func TestReflowsLongMessages(t *testing.T) {
	block, _ := txt.ParseBlock("Foo bar", 1)
	err := app.NewParserErrors([]txt.Error{
		txt.NewError(block, 0, 4, 3, "CODE", "Some Title", "A verbose description with details, potentially spanning multiple lines with a comprehensive text and tremendously helpful information.\nBut\nit\nrespects\nnewlines."),
	})
	text := PrettifyParsingError(err, styler).Error()
	assert.Equal(t, `
[SYNTAX ERROR] in line 2
    Foo bar
        ^^^
    Some Title: A verbose description with details, potentially spanning
    multiple lines with a comprehensive text and tremendously helpful
    information.
    But
    it
    respects
    newlines.
`, tf.StripAllAnsiSequences(text))
}

func TestConvertsTabToSpaces(t *testing.T) {
	block, _ := txt.ParseBlock("\tFoo\tbar", 13)
	err := app.NewParserErrors([]txt.Error{
		txt.NewError(block, 0, 0, 8, "CODE", "Error title", "Error details"),
	})
	text := PrettifyParsingError(err, styler).Error()
	assert.Equal(t, `
[SYNTAX ERROR] in line 14
     Foo bar
    ^^^^^^^^
    Error title: Error details
`, tf.StripAllAnsiSequences(text))
}

func TestFormatAppError(t *testing.T) {
	err := app.NewError("Some message", "A more detailed explanation", nil)
	text := PrettifyAppError(err, false).Error()
	assert.Equal(t, `Error: Some message
A more detailed explanation`, text)
}

func TestFormatAppErrorWithDebugFlag(t *testing.T) {
	textWithNilErr := PrettifyAppError(
		app.NewError("Some message", "A more detailed explanation", nil),
		true).Error()
	assert.Equal(t, `Error: Some message
A more detailed explanation`, textWithNilErr)

	textWithErr := PrettifyAppError(
		app.NewError("Some message", "A more detailed explanation", errors.New("ORIG_ERR")),
		true).Error()
	assert.Equal(t, `Error: Some message
A more detailed explanation

Original Error:
ORIG_ERR`, textWithErr)
}
