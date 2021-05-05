package lib

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"klog/app"
	"klog/parser/parsing"
	"regexp"
	"testing"
)

var ansiSequencePattern = regexp.MustCompile(`\x1b\[.+?m`)

func stripAllAnsiSequences(text string) string {
	return ansiSequencePattern.ReplaceAllString(text, "")
}

func TestFormatParserError(t *testing.T) {
	err := parsing.NewErrors([]parsing.Error{
		func() parsing.Error {
			err := parsing.NewError(parsing.NewLineFromString("Foo bar", 2), 4, 3)
			return err.Set("CODE", "Some Title", "A verbose description with details, potentially spanning multiple lines with a comprehensive text and tremendously helpful information.")
		}(),
		func() parsing.Error {
			err := parsing.NewError(parsing.NewLineFromString("Some malformed text", 39), 0, 4)
			return err.Set("CODE", "Error", "Short explanation.")
		}(),
	})
	text := PrettifyError(err, false).Error()
	assert.Equal(t, ` ERROR in line 2: 
    Foo bar
        ^^^
    Some Title: A verbose description with details, potentially
    spanning multiple lines with a comprehensive text
    and tremendously helpful information.

 ERROR in line 39: 
    Some malformed text
    ^^^^
    Error: Short explanation.

`, stripAllAnsiSequences(text))
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
