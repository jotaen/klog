package kfl

import (
	"errors"
	"regexp"
)

type tokenKind int

const (
	tokenOpenBracket tokenKind = iota
	tokenCloseBracket
	tokenAnd
	tokenOr
	tokenNot
	tokenDate
	tokenPeriod
	tokenDateRange
	tokenTag
	tokenEntryType
)

type token struct {
	kind     tokenKind
	value    string
	position int
}

var (
	tagRegex       = regexp.MustCompile(`^(#([\p{L}\d_-]+)(=(("[^"]*")|('[^']*')|([\p{L}\d_-]*)))?)`)
	dateRangeRegex = regexp.MustCompile(`^(((\d{4}-\d{2}-\d{2})|(\d{4}-\p{L}?\d+)|(\d{4}))?\.{3}((\d{4}-\d{2}-\d{2})|(\d{4}-\p{L}?\d+)|(\d{4}))?)`)
	dateRegex      = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})`)
	periodRegex    = regexp.MustCompile(`^((\d{4}-\p{L}?\d+)|(\d{4}))`)
	typeRegex      = regexp.MustCompile(`^(type:[\p{L}\-_]+)`)
)

var (
	ErrMissingWhiteSpace = errors.New("Missing whitespace. Please separate operands and operators with whitespace.")
	ErrUnrecognisedToken = errors.New("Unrecognised query token. Please make sure to use valid query syntax.")
)

func tokenise(filterQuery string) ([]token, ParseError) {
	txtParser := newTextParser(filterQuery)
	tokens := []token{}
	for {
		if txtParser.isFinished() {
			break
		}

		if txtParser.peekString(" ") {
			txtParser.advance(1)
		} else if txtParser.peekString("(") {
			tokens = append(tokens, token{tokenOpenBracket, "(", txtParser.pointer})
			txtParser.advance(1)
		} else if txtParser.peekString(")") {
			tokens = append(tokens, token{tokenCloseBracket, ")", txtParser.pointer})
			txtParser.advance(1)
			if !txtParser.peekString(EOT, " ", ")") {
				return nil, parseError{
					err:      ErrMissingWhiteSpace,
					position: txtParser.pointer - 1,
					length:   1,
				}
			}
		} else if txtParser.peekString("&&") {
			tokens = append(tokens, token{tokenAnd, "&&", txtParser.pointer})
			txtParser.advance(2)
			if !txtParser.peekString(EOT, " ") {
				return nil, parseError{
					err:      ErrMissingWhiteSpace,
					position: txtParser.pointer,
					length:   1,
				}
			}
		} else if txtParser.peekString("||") {
			tokens = append(tokens, token{tokenOr, "||", txtParser.pointer})
			txtParser.advance(2)
			if !txtParser.peekString(EOT, " ") {
				return nil, parseError{
					err:      ErrMissingWhiteSpace,
					position: txtParser.pointer,
					length:   1,
				}
			}
		} else if txtParser.peekString("!") {
			tokens = append(tokens, token{tokenNot, "!", txtParser.pointer})
			txtParser.advance(1)
		} else if tm := txtParser.peekRegex(tagRegex); tm != nil {
			value := tm[1]
			tokens = append(tokens, token{tokenTag, value, txtParser.pointer})
			txtParser.advance(len(value))
			if !txtParser.peekString(EOT, " ", ")") {
				return nil, parseError{
					err:      ErrMissingWhiteSpace,
					position: txtParser.pointer,
					length:   1,
				}
			}
		} else if ym := txtParser.peekRegex(typeRegex); ym != nil {
			tokens = append(tokens, token{tokenEntryType, ym[1], txtParser.pointer})
			txtParser.advance(len(ym[1]))
			if !txtParser.peekString(EOT, " ", ")") {
				return nil, parseError{
					err:      ErrMissingWhiteSpace,
					position: txtParser.pointer - 1,
					length:   1,
				}
			}
		} else if rm := txtParser.peekRegex(dateRangeRegex); rm != nil {
			value := rm[1]
			tokens = append(tokens, token{tokenDateRange, value, txtParser.pointer})
			txtParser.advance(len(value))
			if !txtParser.peekString(EOT, " ", ")") {
				return nil, parseError{
					err:      ErrMissingWhiteSpace,
					position: txtParser.pointer - 1,
					length:   1,
				}
			}
		} else if dm := txtParser.peekRegex(dateRegex); dm != nil {
			value := dm[1]
			tokens = append(tokens, token{tokenDate, value, txtParser.pointer})
			txtParser.advance(len(value))
			if !txtParser.peekString(EOT, " ", ")") {
				return nil, parseError{
					err:      ErrMissingWhiteSpace,
					position: txtParser.pointer - 1,
					length:   1,
				}
			}
		} else if pm := txtParser.peekRegex(periodRegex); pm != nil {
			value := pm[1]
			tokens = append(tokens, token{tokenPeriod, value, txtParser.pointer})
			txtParser.advance(len(value))
			if !txtParser.peekString(EOT, " ", ")") {
				return nil, parseError{
					err:      ErrMissingWhiteSpace,
					position: txtParser.pointer - 1,
					length:   1,
				}
			}
		} else {
			return nil, parseError{
				err:      ErrUnrecognisedToken,
				position: txtParser.pointer,
				length:   1,
			}
		}
	}
	return tokens, nil
}
