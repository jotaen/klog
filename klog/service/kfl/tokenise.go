package kfl

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type token interface{}

type tokenOpenBracket struct{}
type tokenCloseBracket struct{}
type tokenAnd struct{}
type tokenOr struct{}
type tokenNot struct{}
type tokenDate struct {
	date string
}
type tokenPeriod struct {
	period string
}
type tokenDateRange struct {
	bounds []string
}
type tokenTag struct {
	tag string
}
type tokenEntryType struct {
	entryType string
}

var (
	tagRegex       = regexp.MustCompile(`^#(([\p{L}\d_-]+)(=(("[^"]*")|('[^']*')|([\p{L}\d_-]*)))?)`)
	dateRangeRegex = regexp.MustCompile(`^((\d{4}-\d{2}-\d{2})?\.\.\.(\d{4}-\d{2}-\d{2})?)`)
	dateRegex      = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})`)
	periodRegex    = regexp.MustCompile(`^((\d{4}-\p{L}?\d+)|(\d{4}))`)
	typeRegex      = regexp.MustCompile(`^type:([\p{L}\-_]+)`)
)

var (
	ErrMissingWhiteSpace = errors.New("Missing whitespace. Please separate operands and operators with whitespace.")
	ErrUnrecognisedToken = errors.New("Unrecognised query token. Please make sure to use valid query syntax.")
)

func tokenise(filterQuery string) ([]token, error) {
	txtParser := newTextParser(filterQuery)
	tokens := []token{}
	for {
		if txtParser.isFinished() {
			break
		}

		if txtParser.peekString(" ") {
			txtParser.advance(1)
		} else if txtParser.peekString("(") {
			tokens = append(tokens, tokenOpenBracket{})
			txtParser.advance(1)
		} else if txtParser.peekString(")") {
			tokens = append(tokens, tokenCloseBracket{})
			txtParser.advance(1)
			if !txtParser.peekString(EOT, " ", ")") {
				return nil, ErrMissingWhiteSpace
			}
		} else if txtParser.peekString("&&") {
			tokens = append(tokens, tokenAnd{})
			txtParser.advance(2)
			if !txtParser.peekString(EOT, " ") {
				return nil, ErrMissingWhiteSpace
			}
		} else if txtParser.peekString("||") {
			tokens = append(tokens, tokenOr{})
			txtParser.advance(2)
			if !txtParser.peekString(EOT, " ") {
				return nil, ErrMissingWhiteSpace
			}
		} else if txtParser.peekString("!") {
			tokens = append(tokens, tokenNot{})
			txtParser.advance(1)
		} else if tm := txtParser.peekRegex(tagRegex); tm != nil {
			value := tm[1]
			tokens = append(tokens, tokenTag{value})
			txtParser.advance(1 + len(value))
			if !txtParser.peekString(EOT, " ", ")") {
				return nil, ErrMissingWhiteSpace
			}
		} else if ym := txtParser.peekRegex(typeRegex); ym != nil {
			tokens = append(tokens, tokenEntryType{ym[1]})
			txtParser.advance(5 + len(ym[1]))
			if !txtParser.peekString(EOT, " ", ")") {
				return nil, ErrMissingWhiteSpace
			}
		} else if rm := txtParser.peekRegex(dateRangeRegex); rm != nil {
			value := rm[1]
			parts := strings.Split(value, "...")
			tokens = append(tokens, tokenDateRange{parts})
			txtParser.advance(len(value))
			if !txtParser.peekString(EOT, " ", ")") {
				return nil, ErrMissingWhiteSpace
			}
		} else if dm := txtParser.peekRegex(dateRegex); dm != nil {
			value := dm[1]
			tokens = append(tokens, tokenDate{value})
			txtParser.advance(len(value))
			if !txtParser.peekString(EOT, " ", ")") {
				return nil, ErrMissingWhiteSpace
			}
		} else if pm := txtParser.peekRegex(periodRegex); pm != nil {
			value := pm[1]
			tokens = append(tokens, tokenPeriod{value})
			txtParser.advance(len(value))
			if !txtParser.peekString(EOT, " ", ")") {
				return nil, ErrMissingWhiteSpace
			}
		} else {
			return nil, fmt.Errorf("%w: %s", ErrUnrecognisedToken, txtParser.remainder())
		}
	}
	return tokens, nil
}
