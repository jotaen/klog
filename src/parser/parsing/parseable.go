package parsing

type Parseable struct {
	Line
	PointerPosition int
}

var END_OF_TEXT int32 = -1

func NewParseable(l Line) Parseable {
	return Parseable{
		PointerPosition: 0,
		Line:            l,
	}
}

func (p *Parseable) Peek() rune {
	char := SubRune(p.Value, p.PointerPosition, 1)
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
	for i := p.PointerPosition; i < len(p.Value); i++ {
		next := SubRune(p.Value, i, 1)
		if isMatch(next[0]) {
			return result, true
		}
		result.Value = append(result.Value, next[0])
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
	return len(p.Value)
}

func (p *Parseable) RemainingLength() int {
	return p.Length() - p.PointerPosition
}

func (p *Parseable) ToString() string {
	return string(p.Value)
}
