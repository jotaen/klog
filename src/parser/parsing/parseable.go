package parsing

type Parseable struct {
	Line
	Chars           []rune
	PointerPosition int
}

var END_OF_TEXT int32 = -1

func NewParseable(l Line) Parseable {
	return Parseable{
		PointerPosition: 0,
		Chars:           []rune(l.Text),
		Line:            l,
	}
}

func (p *Parseable) Peek() rune {
	char := SubRune(p.Chars, p.PointerPosition, 1)
	if char == nil {
		return END_OF_TEXT
	}
	return char[0]
}

func (p *Parseable) PeekUntil(isMatch func(rune) bool) (Parseable, bool) {
	result := Parseable{
		PointerPosition: p.PointerPosition,
		Line:            Line{},
	}
	for i := p.PointerPosition; i < len(p.Chars); i++ {
		next := SubRune(p.Chars, i, 1)
		if isMatch(next[0]) {
			return result, true
		}
		result.Chars = append(result.Chars, next[0])
	}
	return result, false
}

func (p *Parseable) Advance(increment int) {
	p.PointerPosition += increment
}

func (p *Parseable) SkipWhitespace() {
	for IsWhitespace(p.Peek()) {
		p.Advance(1)
	}
	return
}

func (p *Parseable) Length() int {
	return len(p.Chars)
}

func (p *Parseable) RemainingLength() int {
	return p.Length() - p.PointerPosition
}

func (p *Parseable) ToString() string {
	return string(p.Chars)
}
