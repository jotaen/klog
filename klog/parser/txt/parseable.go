/*
Package engine is a generic utility for parsing and processing a text
that is structured in individual lines.
*/
package txt

// Parseable is utility data structure for processing a Line in a parser.
type Parseable struct {
	Line
	Chars           []rune
	PointerPosition int
}

var END_OF_TEXT int32 = -1

func NewParseable(l Line, startPointerPosition int) *Parseable {
	return &Parseable{
		PointerPosition: startPointerPosition,
		Chars:           []rune(l.Text),
		Line:            l,
	}
}

// Peek returns the next character, or END_OF_TEXT if there is none anymore.
func (p *Parseable) Peek() rune {
	char := SubRune(p.Chars, p.PointerPosition, 1)
	if char == nil {
		return END_OF_TEXT
	}
	return char[0]
}

// PeekUntil moves the cursor forward until the condition is satisfied, or until the end
// of the line is reached. It returns a Parseable containing the consumed part of the line,
// as well as a bool to indicate whether the condition was met (`true`) or the end of the
// line was encountered (`false`).
func (p *Parseable) PeekUntil(isMatch func(rune) bool) (Parseable, bool) {
	result := Parseable{
		PointerPosition: p.PointerPosition,
		Line:            Line{},
	}
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
	result.Chars = SubRune(p.Chars, p.PointerPosition, matchLength)
	return result, hasMatched
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
