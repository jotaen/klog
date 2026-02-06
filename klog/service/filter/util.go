package filter

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
		return token{}
	}
	next := t.tokens[t.pointer]
	t.pointer += 1
	return next
}

func (t *tokenParser) checkNextIsOperand() ParseError {
	if t.pointer >= len(t.tokens) {
		return parseError{
			err:      ErrOperandExpected,
			position: t.tokens[len(t.tokens)-1].position,
			length:   1,
		}
	}
	for _, k := range []tokenKind{
		tokenOpenBracket, tokenTag, tokenDate, tokenDateRange, tokenPeriod, tokenNot, tokenEntryType,
	} {
		if t.tokens[t.pointer].kind == k {
			return nil
		}
	}
	return parseError{
		err:      ErrOperandExpected,
		position: t.tokens[t.pointer].position,
		length:   len(t.tokens[t.pointer].value),
	}
}

func (t *tokenParser) checkNextIsOperatorOrEnd() ParseError {
	if t.pointer >= len(t.tokens) {
		return nil
	}
	for _, k := range []tokenKind{
		tokenCloseBracket, tokenAnd, tokenOr,
	} {
		if t.tokens[t.pointer].kind == k {
			return nil
		}
	}
	return parseError{
		err:      ErrOperatorExpected,
		position: t.tokens[t.pointer].position,
		length:   len(t.tokens[t.pointer].value),
	}
}

type predicateGroup struct {
	ps            []Predicate
	operator      tokenKind // -1 (unset) or tokenAnd or tokenOr
	isNextNegated bool
}

func newPredicateGroup() predicateGroup {
	return predicateGroup{
		ps:            nil,
		operator:      -1,
		isNextNegated: false,
	}
}

func (g *predicateGroup) append(p Predicate) {
	if g.isNextNegated {
		g.isNextNegated = false
		p = Not{p}
	}
	g.ps = append(g.ps, p)
}

func (g *predicateGroup) setOperator(operatorT token, position int) ParseError {
	if g.operator == -1 {
		g.operator = operatorT.kind
	}
	if g.operator != operatorT.kind {
		return parseError{
			err:      ErrCannotMixAndOr,
			position: position,
			length:   2,
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
	} else if g.operator == tokenAnd {
		return And{g.ps}, nil
	} else if g.operator == tokenOr {
		return Or{g.ps}, nil
	} else {
		// This would happen for an empty group.
		return nil, parseError{
			err:      ErrMalformedFilterQuery,
			position: 0,
		}
	}
}
