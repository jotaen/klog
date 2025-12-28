package kfl

import (
	"regexp"
	"strings"
)

const EOT = "" // End of text

type textParser struct {
	text    string
	pointer int
}

func newTextParser(text string) textParser {
	return textParser{
		text:    text,
		pointer: 0,
	}
}

func (t *textParser) isFinished() bool {
	return t.pointer == len(t.text)
}

func (t *textParser) peekString(lookup ...string) bool {
	r := t.remainder()
	for _, l := range lookup {
		if l == EOT {
			if r == EOT {
				return true
			}
		} else if strings.HasPrefix(r, l) {
			return true
		}
	}
	return false
}

func (t *textParser) peekRegex(lookup *regexp.Regexp) []string {
	return lookup.FindStringSubmatch(t.remainder())
}

func (t *textParser) advance(i int) {
	t.pointer += i
}

func (t *textParser) remainder() string {
	if t.isFinished() {
		return ""
	}
	return t.text[t.pointer:]
}

type tokenParser struct {
	tokens  []token
	pos     []int
	pointer int
}

func newTokenParser(ts []token, pos []int) tokenParser {
	return tokenParser{
		tokens:  ts,
		pos:     pos,
		pointer: 0,
	}
}

func (t *tokenParser) next() (token, int) {
	if t.pointer >= len(t.tokens) {
		return nil, -1
	}
	next := t.tokens[t.pointer]
	pos := t.pos[t.pointer]
	t.pointer += 1
	return next, pos
}

func (t *tokenParser) checkNextIsOperand() ParseError {
	if t.pointer >= len(t.tokens) {
		return parseError{
			err:      ErrOperandExpected,
			position: t.pos[t.pointer],
		}
	}
	switch t.tokens[t.pointer].(type) {
	case tokenOpenBracket, tokenTag, tokenDate, tokenDateRange, tokenPeriod, tokenNot, tokenEntryType:
		return nil
	}
	return parseError{
		err:      ErrOperandExpected,
		position: t.pos[t.pointer],
	}
}

func (t *tokenParser) checkNextIsOperatorOrEnd() ParseError {
	if t.pointer >= len(t.tokens) {
		return nil
	}
	switch t.tokens[t.pointer].(type) {
	case tokenCloseBracket, tokenAnd, tokenOr:
		return nil
	}
	return parseError{
		err:      ErrOperatorExpected,
		position: t.pos[t.pointer],
	}
}

type predicateGroup struct {
	ps            []Predicate
	operator      token // nil or tokenAnd or tokenOr
	isNextNegated bool
}

func (g *predicateGroup) append(p Predicate) {
	if g.isNextNegated {
		g.isNextNegated = false
		p = Not{p}
	}
	g.ps = append(g.ps, p)
}

func (g *predicateGroup) setOperator(operatorT token, position int) ParseError {
	if g.operator == nil {
		g.operator = operatorT
	}
	if g.operator != operatorT {
		return parseError{
			err:      ErrCannotMixAndOr,
			position: position,
		}
	}
	return nil
}

func (g *predicateGroup) negateNextOperand() {
	g.isNextNegated = true
}

func (g *predicateGroup) make() (Predicate, ParseError) {
	if len(g.ps) == 1 {
		return g.ps[0], nil
	} else if g.operator == (tokenAnd{}) {
		return And{g.ps}, nil
	} else if g.operator == (tokenOr{}) {
		return Or{g.ps}, nil
	} else {
		return nil, parseError{err: ErrMalformedFilterQuery, position: 0}
	}
}
