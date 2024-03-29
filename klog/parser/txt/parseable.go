/*
Package txt is a generic utility for parsing and processing a text
that is structured in individual lines.
*/
package txt

import "unicode/utf8"

// Parseable is utility data structure for parsing a piece of text.
type Parseable struct {
	Chars           []rune
	PointerPosition int
}

// NewParseable creates a parseable from the given line.
func NewParseable(l Line, startPointerPosition int) *Parseable {
	return &Parseable{
		PointerPosition: startPointerPosition,
		Chars:           []rune(l.Text),
	}
}

// Peek returns the next character, or `utf8.RuneError` if there is none anymore.
func (p *Parseable) Peek() rune {
	char := SubRune(p.Chars, p.PointerPosition, 1)
	if char == nil {
		return utf8.RuneError
	}
	return char[0]
}

// PeekUntil moves the cursor forward until the condition is satisfied, or until the end
// of the line is reached. It returns a Parseable containing the consumed part of the line,
// as well as a bool to indicate whether the condition was met (`true`) or the end of the
// line was encountered (`false`).
func (p *Parseable) PeekUntil(isMatch func(rune) bool) (Parseable, bool) {
	matchLength := 0
	hasMatched := false
	for i := p.PointerPosition; i < len(p.Chars); i++ {
		matchLength++
		if isMatch(SubRune(p.Chars, i, 1)[0]) {
			matchLength -= 1
			hasMatched = true
			break
		}
	}
	return Parseable{
		PointerPosition: p.PointerPosition,
		Chars:           SubRune(p.Chars, p.PointerPosition, matchLength),
	}, hasMatched
}

// Remainder returns the rest of the text.
func (p *Parseable) Remainder() Parseable {
	rest, _ := p.PeekUntil(Is(utf8.RuneError))
	return rest
}

// Advance moves forward the cursor position.
func (p *Parseable) Advance(increment int) {
	p.PointerPosition += increment
}

// SkipWhile consumes all upcoming characters that match the predicate.
func (p *Parseable) SkipWhile(isMatch func(rune) bool) {
	for isMatch(p.Peek()) {
		p.Advance(1)
	}
}

// Length returns the total length of the line.
func (p *Parseable) Length() int {
	return len(p.Chars)
}

// RemainingLength returns the number of chars until the end of the line.
func (p *Parseable) RemainingLength() int {
	return p.Length() - p.PointerPosition
}

// ToString returns the line text as string.
func (p *Parseable) ToString() string {
	return string(p.Chars)
}
