package engine

import "strings"

// NewIndentator creates an indentator object, if the given line is indented
// according to the allowed styles.
func NewIndentator(allowedIndentationStyles []string, l Line) *Indentator {
	for _, s := range allowedIndentationStyles {
		if strings.HasPrefix(l.Text, s) {
			return &Indentator{s}
		}
	}
	return nil
}

// Indentator is a utility to check to consistently process indentated text.
type Indentator struct {
	indentationStyle string
}

// NewIndentedParseable returns a Parseable with already skipped indentation.
// It returns `nil` if the encountered indentation level is smaller than `atLevel`.
// It only consumes the desired indentation and disregards any additional indentation.
func (i *Indentator) NewIndentedParseable(l Line, atLevel int) *Parseable {
	expectedIndentation := strings.Repeat(i.indentationStyle, atLevel)
	if !strings.HasPrefix(l.Text, expectedIndentation) {
		return nil
	}
	return NewParseable(l, len(expectedIndentation))
}

func (i *Indentator) Style() string {
	return i.indentationStyle
}
