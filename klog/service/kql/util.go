package kql

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
	pointer int
}

func newTokenParser(ts []token) tokenParser {
	return tokenParser{
		tokens:  ts,
		pointer: 0,
	}
}

func (t *tokenParser) next() token {
	if t.pointer >= len(t.tokens) {
		return nil
	}
	next := t.tokens[t.pointer]
	t.pointer += 1
	return next
}

func (t *tokenParser) checkNextIsOperand() error {
	if t.pointer >= len(t.tokens) {
		return ErrOperandExpected
	}
	switch t.tokens[t.pointer].(type) {
	case tokenOpenBracket, tokenTag, tokenDate, tokenDateRange, tokenPeriod, tokenNot:
		return nil
	}
	return ErrOperandExpected
}

func (t *tokenParser) checkNextIsOperatorOrEnd() error {
	if t.pointer >= len(t.tokens) {
		return nil
	}
	switch t.tokens[t.pointer].(type) {
	case tokenCloseBracket, tokenAnd, tokenOr:
		return nil
	}
	return ErrOperatorExpected
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

func (g *predicateGroup) setOperator(operatorT token) error {
	if g.operator == nil {
		g.operator = operatorT
	}
	if g.operator != operatorT {
		return ErrCannotMixAndOr
	}
	return nil
}

func (g *predicateGroup) negateNextOperand() {
	g.isNextNegated = true
}

func (g *predicateGroup) make() (Predicate, error) {
	if len(g.ps) == 1 {
		return g.ps[0], nil
	} else if g.operator == (tokenAnd{}) {
		return And{g.ps}, nil
	} else if g.operator == (tokenOr{}) {
		return Or{g.ps}, nil
	} else {
		return nil, ErrMalformedQuery
	}
}
