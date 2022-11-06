package lib

import (
	"errors"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib/terminalformat"
	"github.com/jotaen/klog/klog/parser/txt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatParserError(t *testing.T) {
	block1, _ := txt.ParseBlock("Good text\nSome malformed text", 37)
	block2, _ := txt.ParseBlock("Another issue!", 133)
	err := app.NewParserErrors([]txt.Error{
		txt.NewError(block1, 1, 0, 4, "CODE", "Error", "Short explanation."),
		txt.NewError(block2, 0, 8, 5, "CODE", "Problem", "More info."),
	})
	text := PrettifyError(err, false).Error()
	assert.Equal(t, ` ERROR in line 39: 
    Some malformed text
    ^^^^
    Error: Short explanation.

 ERROR in line 134: 
    Another issue!
            ^^^^^
    Problem: More info.

`, terminalformat.StripAllAnsiSequences(text))
}

func TestReflowsLongMessages(t *testing.T) {
	block, _ := txt.ParseBlock("Foo bar", 1)
	err := app.NewParserErrors([]txt.Error{
		txt.NewError(block, 0, 4, 3, "CODE", "Some Title", "A verbose description with details, potentially spanning multiple lines with a comprehensive text and tremendously helpful information.\nBut it respects newlines."),
	})
	text := PrettifyError(err, false).Error()
	assert.Equal(t, ` ERROR in line 2: 
    Foo bar
        ^^^
    Some Title: A verbose description with details, potentially
    spanning multiple lines with a comprehensive text
    and tremendously helpful information.
    But it respects newlines.

`, terminalformat.StripAllAnsiSequences(text))
}

func TestConvertsTabToSpaces(t *testing.T) {
	block, _ := txt.ParseBlock("\tFoo\tbar", 13)
	err := app.NewParserErrors([]txt.Error{
		txt.NewError(block, 0, 0, 8, "CODE", "Error title", "Error details"),
	})
	text := PrettifyError(err, false).Error()
	assert.Equal(t, ` ERROR in line 14: 
     Foo bar
    ^^^^^^^^
    Error title: Error details

`, terminalformat.StripAllAnsiSequences(text))
}

func TestFormatAppError(t *testing.T) {
	err := app.NewError("Some message", "A more detailed explanation", nil)
	text := PrettifyError(err, false).Error()
	assert.Equal(t, `Error: Some message
A more detailed explanation`, text)
}

func TestFormatAppErrorWithDebugFlag(t *testing.T) {
	textWithNilErr := PrettifyError(
		app.NewError("Some message", "A more detailed explanation", nil),
		true).Error()
	assert.Equal(t, `Error: Some message
A more detailed explanation`, textWithNilErr)

	textWithErr := PrettifyError(
		app.NewError("Some message", "A more detailed explanation", errors.New("ORIG_ERR")),
		true).Error()
	assert.Equal(t, `Error: Some message
A more detailed explanation

Original Error:
ORIG_ERR`, textWithErr)
}

func TestFormatRegularError(t *testing.T) {
	textWithNilErr := PrettifyError(errors.New("Some plain error"), true).Error()
	assert.Equal(t, `Error: Some plain error`, textWithNilErr)
}
