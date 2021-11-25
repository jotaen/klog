package lib

import (
	"errors"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/lib/terminalformat"
	"github.com/jotaen/klog/src/parser/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatParserError(t *testing.T) {
	err := engine.NewErrors([]engine.Error{
		engine.NewError(engine.NewLineFromString("Foo bar", 2), 4, 3, "CODE", "Some Title", "A verbose description with details, potentially spanning multiple lines with a comprehensive text and tremendously helpful information.\nBut it respects newlines."),
		engine.NewError(engine.NewLineFromString("Some malformed text", 39), 0, 4, "CODE", "Error", "Short explanation."),
	})
	text := PrettifyError(err, false).Error()
	assert.Equal(t, ` ERROR in line 2: 
    Foo bar
        ^^^
    Some Title: A verbose description with details, potentially
    spanning multiple lines with a comprehensive text
    and tremendously helpful information.
    But it respects newlines.

 ERROR in line 39: 
    Some malformed text
    ^^^^
    Error: Short explanation.

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
