package engine

// Text is a single line of characters.
type Text struct {
	Value            []rune
	PointerPosition  int
	LineNumber       int
	IndentationLevel int
}

// Chunk is a paragraph of Text (i.e. a block of subsequent lines).
type Chunk []Text

var END_OF_TEXT int32 = -1

func (c Chunk) Pop() Chunk {
	return c[1:]
}

func (t *Text) Peek() rune {
	char := SubRune(t.Value, t.PointerPosition, 1)
	if char == nil {
		return END_OF_TEXT
	}
	return char[0]
}

func (t *Text) PeekUntil(isMatch func(rune) bool) (Text, bool) {
	result := Text{
		PointerPosition: t.PointerPosition,
		Value:           nil,
		LineNumber:      t.LineNumber,
	}
	for i := t.PointerPosition; i < len(t.Value); i++ {
		next := SubRune(t.Value, i, 1)
		if isMatch(next[0]) {
			return result, true
		}
		result.Value = append(result.Value, next[0])
	}
	return result, false
}

func (t *Text) Advance(increment int) {
	t.PointerPosition += increment
}

func (t *Text) SkipWhitespace() {
	for IsWhitespace(t.Peek()) {
		t.Advance(1)
	}
	return
}

func (t *Text) Length() int {
	return len(t.Value)
}

func (t *Text) RemainingLength() int {
	return t.Length() - t.PointerPosition
}

func (t *Text) ToString() string {
	return string(t.Value)
}
